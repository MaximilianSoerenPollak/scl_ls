package internal

import (
	"log"

	"errors"

	"sclls/lsp"
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
	return State{Documents: m, NeedsList: needsList, ServerConfig: srvConfig}
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
	for _, need := range docInfo.Needs {
		for _, p := range need.Positions {
			if pos.Line == p.Line && pos.Character >= p.StartCol && pos.Character <= p.EndCol {
				return need.Need, nil
			}
		}
	}
	return Need{}, errors.New("could not find information for requested position")
}

func (s *State) GoToDefinition(id int, docUri string, pos lsp.Position) lsp.DefinitionResponse {		
	l := lsp.Location{URI: docUri}

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.DefinitionResult{
			Location: l,

		},
	}
}
