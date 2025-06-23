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
func (s *State) OpenDocument(uri string, content string) []lsp.Diagnostic {
	di, ok := s.Documents[uri] //
	if !ok {
		// Document not yet in our map
		s.Logger.Printf("Document %s not found in state.documents. Initializing new DocumentInfo.", uri)
		newDocInfo := &DocumentInfo{}
		s.Documents[uri] = newDocInfo
		di = newDocInfo
	}
	documentNeeds := NewDocumentNeeds(uri, s.Logger)
	di.Content = content
	byteContent := []byte(content)
	ndi := FindAllNeedsPositions(byteContent, s.NeedsList)
	diagnostics := s.FindDiagnosticsInDocument(byteContent)
	documentNeeds.Needs = ndi
	di.DocumentNeeds = documentNeeds
	di.Needs = ndi
	di.Diagnostics = diagnostics
	if diagnostics == nil {
		s.Logger.Printf("I think diagnostics is empty: %v", diagnostics)
		return []lsp.Diagnostic{}
	}
	return diagnostics
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
	di.Content = content
	di.Diagnostics = diagnostics
	if diagnostics == nil {
		s.Logger.Printf("I think diagnostics is empty: %v", diagnostics)
		return []lsp.Diagnostic{}
	}
	return diagnostics
}

func (s *State) FindDiagnosticsInDocument(content []byte) []lsp.Diagnostic {
	var diagnostics = []lsp.Diagnostic{}

	reader := bytes.NewReader(content)
	scanner := bufio.NewScanner(reader)
	lineNr := -1

	for scanner.Scan() {
		lineNr++
		lineTxt := scanner.Text()

		s.Logger.Printf("Diagnostics: Processing line %d: '%s'", lineNr, lineTxt)

		var templateFound bool
		var matchedTemplatePrefix string

		for _, tmplStr := range s.TemplateStrings {
			if strings.HasPrefix(lineTxt, tmplStr) {
				matchedTemplatePrefix = tmplStr
				templateFound = true
				break
			}
		}

		if templateFound { // Only proceed if a template prefix was found
			s.Logger.Println("We found a template string, now going further")
			contentAfterPrefix := strings.TrimPrefix(lineTxt, matchedTemplatePrefix)
			prefixLength := len(matchedTemplatePrefix)
			if strings.TrimSpace(contentAfterPrefix) == "" {
				diagnostics = append(diagnostics, lsp.Diagnostic{
					Range: lsp.Range{
						Start: lsp.Position{
							Line:      lineNr,
							Character: len(matchedTemplatePrefix),
						},
						End: lsp.Position{
							Line:      lineNr,
							Character: len(matchedTemplatePrefix),
						},
					},
					Severity: 2,
					Source:   "scl_lsp",
					Message:  fmt.Sprintf("Found template string but no need after."),
				})
				continue
			}

			potentialNeedsFound := strings.Split(contentAfterPrefix, ",")
			s.Logger.Printf("Potential Needs is: %d for line: %d", len(potentialNeedsFound), lineNr)

			currentOffsetInSuffix := 0

			for _, drtyNeed := range potentialNeedsFound {
				trimmedNeed := strings.TrimSpace(drtyNeed)

				offsetWithinDirtyPart := strings.Index(drtyNeed, trimmedNeed)
				if offsetWithinDirtyPart == -1 {
					offsetWithinDirtyPart = 0
				}

				charStart := prefixLength + currentOffsetInSuffix + offsetWithinDirtyPart
				charEnd := charStart + len(trimmedNeed)

				s.Logger.Printf("Diagnostics: On line %d, found dirty part: '%s', trimmed: '%s'. Range: [%d,%d]",
					lineNr, drtyNeed, trimmedNeed, charStart, charEnd)

				if trimmedNeed == "" {
					// This case handles empty parts from trailing commas or double commas (e.g., "ID1,,ID2")
					// We need to advance the offset for the next iteration.
					// currentOffsetInSuffix needs to move past the original drtyNeed (even if empty) and the comma.
					currentOffsetInSuffix += len(drtyNeed) + len(",") // len(drtyNeed) for empty is 0
					continue                                          // Skip processing this empty need
				}

				// Check if the trimmedNeed exists in your NeedsList
				_, ok := s.NeedsList[trimmedNeed] // Assuming s.NeedsList is map[string]Need
				if !ok {
					s.Logger.Printf("Diagnostics: Unknown need '%s' on line %d.", trimmedNeed, lineNr)
					diagnostics = append(diagnostics, lsp.Diagnostic{
						Range: lsp.Range{
							Start: lsp.Position{
								Line:      lineNr,
								Character: charStart,
							},
							End: lsp.Position{
								Line:      lineNr,
								Character: charEnd,
							},
						},
						Severity: 1,
						Source:   "scl_lsp",
						Message:  fmt.Sprintf("Need '%s' not found. Typo or missing definition?", trimmedNeed),
					})
				}

				// This must account for the full length of the *original* dirty part, plus the comma that separated it.
				// For example, if "A, B", first drtyNeed is "A", next part starts after "A,".
				currentOffsetInSuffix += len(drtyNeed) + len(",")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		s.Logger.Printf("Error scanning document content for diagnostics: %v", err)
	}

	s.Logger.Printf("FindDiagnosticsInDocument returning %d diagnostics.", len(diagnostics))
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
		s.Logger.Printf("Definition: Did not find need definition requested. Error: %s", err.Error())
		// Need to send error repsonse instead then in the future
		return lsp.DefinitionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: []lsp.Location{},
		}
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

// Make this only activate when you write one of the template strings
func (s *State) TextDocumentCompletion(id int, docURI string, pos lsp.Position) lsp.CompletionResponse {
	s.Logger.Printf("=== COMPLETION DEBUG ===")
	docInfo := s.Documents[docURI]
	if docInfo == nil {
		s.Logger.Printf("ERROR: Document not found for URI: %s", docURI)
		return lsp.CompletionResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: []lsp.CompletionItem{},
		}
	}
	s.Logger.Printf("Position: Line=%d, Character=%d\n", pos.Line, pos.Character)
	s.Logger.Printf("Document found, content length: %d\n", len(docInfo.Content))
	s.Logger.Printf("Raw document content: %q\n", docInfo.Content)
	lines := strings.Split(docInfo.Content, "\n")
	completionLine := ""
	if int(pos.Line) < len(lines) {
		completionLine = lines[pos.Line]
	} else if int(pos.Line) == len(lines) {
		s.Logger.Printf("Cursor is on a logically new line (%d), treating as empty.", pos.Line)
		completionLine = "" // It's an empty line
	} else {
		s.Logger.Printf("Completion ERROR: Line %d is significantly out of bounds (doc has %d lines). Returning empty.", pos.Line, len(lines))
		s.Logger.Printf("Current document content state:\n---\n%q\n---", docInfo.Content)
		return lsp.CompletionResponse{
			Response: lsp.Response{RPC: "2.0", ID: &id},
			Result:   []lsp.CompletionItem{},
		}
	}
	linePrefix := ""
	if int(pos.Character) <= len(completionLine) { // Use <= here, if cursor is *at* the end of line
		linePrefix = completionLine[:pos.Character]
	} else {
		// This case means pos.Character is beyond the actual length of the text on the line.
		// It implies the user typed past the end or the line is still empty but they moved cursor.
		// We'll treat linePrefix as the entire content of the line, even if character is invalid.
		linePrefix = completionLine
		s.Logger.Printf("Cursor character %d is beyond line content length %d. Using full line as prefix.", pos.Character, len(completionLine))
	}
	s.Logger.Printf("linePrefix: '%s'", linePrefix)

	// toBeCompletedItem should be the fragment *after* the last space/word boundary
	// This is useful for filtering specific keywords (like "req-" or "tool_")
	lastSpace := strings.LastIndexAny(linePrefix, " \t\n\r")
	toBeCompletedItem := linePrefix
	if lastSpace != -1 {
		toBeCompletedItem = linePrefix[lastSpace+1:]
	}
	s.Logger.Printf("toBeCompletedItem: '%s'", toBeCompletedItem)

	// Label = What we want to complete
	var items []lsp.CompletionItem
	hasReqPrefix := strings.HasPrefix(toBeCompletedItem, "# req-")
	hasTraceability := strings.Contains(completionLine, "# req-traceability:")
	hasReqId := strings.Contains(completionLine, "# req-Id:")

	s.Logger.Printf("Condition 1 - HasPrefix('# req-'): %v", hasReqPrefix)
	s.Logger.Printf("Condition 2a - Contains('# req-traceability:'): %v", hasTraceability)
	s.Logger.Printf("Condition 2b - Contains('# req-Id:'): %v", hasReqId)
	s.Logger.Printf("NeedsList length: %d", len(s.NeedsList))
	// NeedToCheck if this is okay.
	// TODO: Pre-compute this once? Might be nicer to do this.
	//s.Logger.Printf("This is toBeCompletedItem: %v", toBeCompletedItem)
	s.Logger.Printf("To be completed item: %+v", toBeCompletedItem)
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
	}
	s.Logger.Printf("NeedsList after case 1. length: %d", len(s.NeedsList))
	if strings.Contains(linePrefix, "req-Id: ") || strings.Contains(linePrefix, "req-traceability: ") {
		// Find what comes after the colon and space
		var afterColon string
		if idx := strings.Index(linePrefix, "req-Id: "); idx != -1 {
			afterColon = linePrefix[idx+len("req-Id: "):]
		} else if idx := strings.Index(linePrefix, "req-traceability: "); idx != -1 {
			afterColon = linePrefix[idx+len("req-traceability: "):]
		}

		// Remove any spaces at the beginning
		afterColon = strings.TrimLeft(afterColon, " \t")

		// If nothing typed yet, show all needs
		if afterColon == "" {
			for _, need := range s.NeedsList {
				items = append(items, need.GenerateCompletionInfo())
			}
		} else {
			// Filter needs based on what's already typed
			for _, need := range s.NeedsList {
				if strings.HasPrefix(strings.ToLower(need.ID), strings.ToLower(afterColon)) {
					items = append(items, need.GenerateCompletionInfo())
				}
			}
		}
	}
	s.Logger.Printf("Final items count: %d", len(items))
	s.Logger.Printf("=== END COMPLETION DEBUG ===")
	return lsp.CompletionResponse{
		Response: lsp.Response{RPC: "2.0", ID: &id},
		Result:   items,
	}
}
