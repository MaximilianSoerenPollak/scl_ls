package internal

import (
	"log"
	"sclls/lsp"
	"strings"
)

type State struct {
	// Document URI => Information
	Documents map[string]DocumentInfo
	NeedsList NeedsInfo
	ServerConfig
	Logger *log.Logger
}

func NewState(srvConfig ServerConfig, logger *log.Logger) State {
	needsJson := ParseNeedsJson(srvConfig.NeedsJsonPath, logger)
	needsList := GetNeedsList(needsJson)
	m := make(map[string]DocumentInfo)
	return State{Documents: m, NeedsList: needsList, ServerConfig: srvConfig, Logger: logger}
}

// Need to have a check here if the document is already in the thing
func (s *State) OpenDocument(uri string, content string) {
	// Document already exists and was opened again
	_, ok := s.Documents[uri]
	if ok {
		return
	}
	var di DocumentInfo
	documentNeeds := NewDocumentNeeds(uri, s.Logger)
	di.Content = content
	byteContent := []byte(content)
	ndi := FindAllNeedsPositions(byteContent, s.NeedsList)
	documentNeeds.Needs = ndi
	di.DocumentNeeds = documentNeeds
	s.Documents[uri] = di
}

func (s *State) UpdateDocument(uri string, content string) {
	di := s.Documents[uri]
	di.Content = content
	byteContent := []byte(content)
	ndi := FindAllNeedsPositions(byteContent, s.NeedsList)
	di.Needs = ndi
	s.Documents[uri] = di
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
	if toBeCompletedItem == "" {
		for _, need := range s.NeedsList {
			items = append(items, need.GenerateCompletionInfo())
		}
	}

	for _, need := range s.NeedsList {
		if strings.HasPrefix(need.ID, toBeCompletedItem) {
			items = append(items, need.GenerateCompletionInfo())
		}
	}
	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
}
