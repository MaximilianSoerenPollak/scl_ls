package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"sclls/lsp"
	"strings"
)

type State struct {
	// Document URI => Information
	Documents map[string]*DocumentInfo
	NeedsList NeedsInfo
	ServerConfig
	Logger *log.Logger
}

func NewState(srvConfig ServerConfig, logger *log.Logger) State {
	needsJson := ParseNeedsJson(srvConfig.NeedsJsonPath, logger)
	needsList := GetNeedsList(needsJson)
	m := make(map[string]*DocumentInfo)
	return State{Documents: m, NeedsList: needsList, ServerConfig: srvConfig, Logger: logger}
}

// Need to have a check here if the document is already in the thing
func (s *State) OpenDocument(uri string, content string) {
	// Document already exists and was opened again
	di, ok := s.Documents[uri] //
	if !ok {
		// The document doesn't exist in our map yet.
		// This means didOpen wasn't called, or an error occurred.
		// We need to initialize it.
		s.Logger.Printf("Document %s not found in state.documents. Initializing new DocumentInfo.", uri)
		newDocInfo := &DocumentInfo{} // Create a new DocumentInfo instance
		s.Documents[uri] = newDocInfo // Store the pointer to the new instance
		di = newDocInfo               // Use this new instance for current operations
	}
	documentNeeds := NewDocumentNeeds(uri, s.Logger)
	di.Content = content
	byteContent := []byte(content)
	ndi := FindAllNeedsPositions(byteContent, s.NeedsList)
	documentNeeds.Needs = ndi
	di.DocumentNeeds = documentNeeds
}

func (s *State) UpdateDocument(uri string, content string) []lsp.Diagnostic {
	di, ok := s.Documents[uri] //
	if !ok {
		// The document doesn't exist in our map yet.
		// This means didOpen wasn't called, or an error occurred.
		// We need to initialize it.
		s.Logger.Printf("Document %s not found in state.documents. Initializing new DocumentInfo.", uri)
		newDocInfo := &DocumentInfo{} // Create a new DocumentInfo instance
		s.Documents[uri] = newDocInfo // Store the pointer to the new instance
		di = newDocInfo               // Use this new instance for current operations
	}
	byteContent := []byte(content)
	ndi := FindAllNeedsPositions(byteContent, s.NeedsList)
	diagnostics := s.FindDiagnosticsInDocument(byteContent)
	di.Needs = ndi
	di.Diagnostics = diagnostics
	if diagnostics == nil {
		return []lsp.Diagnostic{}
	}
	return diagnostics
}

func (s *State) FindDiagnosticsInDocument(content []byte) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic
	// Maybe we can make it here so that if you have a template string and we don't find a need corresonding
	// We throw an error that this need doesn't exists?
	//
	// 1. Search contents for template string start rows
	reader := bytes.NewReader(content)
	scanner := bufio.NewScanner(reader)
	lineNr := 0
	for scanner.Scan() {
		lineNr++
		lineTxt := scanner.Text()
		var cleanStr string
		for _, tmplStr := range s.TemplateStrings {
			if strings.HasPrefix(lineTxt, tmplStr) {
				// Found a prefix we deem a template string
				cleanStr = strings.ReplaceAll(lineTxt, tmplStr, "")
				break
			}
		}
		potentialNeedsFound := strings.SplitSeq(cleanStr, ",")
		for drtyNeed := range potentialNeedsFound {
			need := strings.TrimSpace(drtyNeed)
			s.Logger.Printf("This is the need we found: %s", need)
			if need == "" {
				break
			}
			_, ok := s.NeedsList[need]
			if !ok {
				idxStrt := strings.Index(need, lineTxt)
				idxEnd := len(need)
				diagnostics = append(diagnostics, lsp.Diagnostic{
					Range: lsp.Range{
						Start: lsp.Position{
							Line:      lineNr,
							Character: idxStrt,
						},
						End: lsp.Position{
							Line:      lineNr,
							Character: idxEnd,
						},
					},
					Severity: 1,
					Source:   "scl_lsp",
					Message:  fmt.Sprintf("Need: %s not found in needs.json. Perhaps you made a typo?", need),
				})
			}
		}
	}
	return diagnostics
}

func (s *State) UpdateNeedsJson(path string) {
	needsJson := ParseNeedsJson(path, s.Logger)
	s.NeedsList = GetNeedsList(needsJson)
}

func (s *State) FindNeedsInRequestedPosition(docURI string, pos lsp.Position) (Need, error) {
	docInfo := s.Documents[docURI]
	// Probably can speed this up somehow
	return docInfo.FindNeedsInPosition(pos)
}

func (s *State) GoToDefinition(id int, docURI string, pos lsp.Position) lsp.DefinitionResponse {
	docInfo := s.Documents[docURI]
	foundNeed, err := docInfo.FindNeedsInPosition(pos)
	if err != nil {
		s.Logger.Println("Did not find need definition requested")
		// Need to send error repsonse instead then in teh future
		s.Logger.Panic(err)
	}
	docName := foundNeed.Docname + ".rst"
	fnDocURI := GetURIFromDocumentName(docName, s.DocumentRootPath)
	s.Logger.Println("Searched for need in document name")
	s.Logger.Printf("DocURI: %s  Lineo: %d", fnDocURI, foundNeed.Lineno)

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: []lsp.Location{{
			URI: fnDocURI,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      foundNeed.Lineno - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      foundNeed.Lineno - 1,
					Character: 0,
				},
			},
		},
		}}
}

// TODO:
// Make this only activate when you write one of the template strings
func (s *State) TextDocumentCompletion(id int, docURI string, pos lsp.Position) lsp.CompletionResponse {
	docInfo := s.Documents[docURI]
	lineNr := 0
	var completionLine string
	for line := range strings.Lines(docInfo.Content) {
		if lineNr == pos.Line {
			completionLine = line
		}
		lineNr++
	}
	lineContentSplit := strings.Split(completionLine, " ")
	toBeCompletedItem := strings.TrimSpace(lineContentSplit[len(lineContentSplit)-1])

	// Label = What we want to complete
	var items []lsp.CompletionItem
	// NeedToCheck if this is okay.
	// TODO: Pre-compute this once? Might be nicer to do this.
	//s.Logger.Printf("This is toBeCompletedItem: %v", toBeCompletedItem)
	if strings.HasPrefix(toBeCompletedItem, "req-") {
		items = append(items, lsp.CompletionItem{
			Label:            "req-Id:",
			Detail:           "Insert a requirement ID placeholder",
			Documentation:    "Use this to link to a specific requirement.",
			InsertText:       "req-Id: ${1:NEED_ID}",
			InsertTextFormat: 2,
		})
		items = append(items, lsp.CompletionItem{
			Label:            "req-traceability:",
			Detail:           "Insert traceability information placeholder",
			Documentation:    "Use this to track the origin or relationships of a component.",
			InsertText:       "req-traceability: ${1:NEED_ID}",
			InsertTextFormat: 2,
		})
	} else if strings.Contains(completionLine, "# req-traceability:") || strings.Contains(completionLine, "# req-Id:") {
		if toBeCompletedItem == "" {
			// Empty start, just give back all suggestions we have
			for _, need := range s.NeedsList {
				items = append(items, need.GenerateCompletionInfo())
			}
		} else {
			// We have a non 'empty' start, therefore can do comparisions
			for _, need := range s.NeedsList {
				if strings.HasPrefix(need.ID, toBeCompletedItem) {
					items = append(items, need.GenerateCompletionInfo())
				}
			}
		}
	}
	//s.Logger.Printf("Sending following completion items: %v", items)
	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
}
