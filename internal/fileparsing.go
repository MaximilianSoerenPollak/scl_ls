package internal

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"path/filepath"

	"sclls/lsp"
)

type DocumentInfo struct {
	Content string
	DocumentNeeds
	Diagnostics []lsp.Diagnostic
}

type DocumentNeeds struct {
	DocName string        `json:"docName"`
	URI     string        `json:"uri"`
	Needs   []NeedDocInfo `json:"needs"`
}

type NeedDocInfo struct {
	Positions []NeedPositionInfo
	Need      `json:"need"` // unsure if this needs a json encoding
}

type NeedPositionInfo struct {
	Line     int `json:"line"`
	StartCol int `json:"startCol"`
	EndCol   int `json:"endCol"`
}

func GetDocumentNameFromURI(uri string) (string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	// Get the path and extract just the filename
	return filepath.Base(parsed.Path), nil
}

func GetURIFromDocumentName(filename string, docPath string) string {
	// Make sure this isn't double-encoding the path
	absPath, err := filepath.Abs(docPath + "/" + filename)
	if err != nil {
		absPath = filename
	}
	// Convert to forward slashes for URI
	uriPath := filepath.ToSlash(absPath)

	// Don't use url.PathEscape here if you're already creating a proper URL
	fileURI := &url.URL{
		Scheme: "file",
		Path:   uriPath, // This should handle encoding properly
	}

	// Might need to make this better?
	return fileURI.String()
}

// Helper I guess?
func FindAllNeedsPositions(content []byte, toBeSearchedNeeds NeedsInfo) []NeedDocInfo {
	var result []NeedDocInfo
	// Search for each string one by one
	for id, need := range toBeSearchedNeeds {
		var ndi NeedDocInfo
		positions := FindNeedPoisiton(content, id)
		if len(positions) == 0 {
			// No positons were found for this need
			continue
		}
		ndi.Positions = positions
		ndi.Need = need
		result = append(result, ndi)
	}
	return result
}

// FindString searches for one string and returns all its positions
func FindNeedPoisiton(content []byte, needID string) []NeedPositionInfo {
	var positions []NeedPositionInfo
	searchBytes := []byte(needID)

	// Keep looking until we can't find any more matches
	start := 0
	for {

		index := bytes.Index(content[start:], searchBytes)
		if index == -1 {
			break // No more matches found
		}

		// Calculate the actual position in the full content
		actualNeedPositionInfo := start + index

		// TODO: Easier way to do this?
		line, col := getLineAndColumn(content, actualNeedPositionInfo)

		// Save this match
		positions = append(positions, NeedPositionInfo{
			Line:     line,
			StartCol: col,
			EndCol:   col + len(needID),
		})

		// Move past this match to look for the next one
		start = actualNeedPositionInfo + 1
	}

	return positions
}

// getLineAndColumn figures out what line and column a position is at
func getLineAndColumn(content []byte, position int) (int, int) {
	line := 0
	col := 0

	// Go through the content character by character until we reach our position
	for i := 0; i < position && i < len(content); i++ {
		if content[i] == '\n' {
			line++  // New line found
			col = 0 // Reset column to start of new line
		} else {
			col++ // Move to next column
		}
	}

	return line, col
}

// TODO: Return error?
func NewDocumentNeeds(uri string, logger *log.Logger) DocumentNeeds {
	docName, err := GetDocumentNameFromURI(uri)
	if err != nil {
		logger.Printf("could not convert URI to document name. URI: %s Error: %s", uri, err.Error())
	}
	return DocumentNeeds{
		DocName: docName,
		URI:     uri,
	}
}

func (di DocumentInfo) FindNeedsInPosition(pos lsp.Position) (Need, error) {
	fmt.Println("INSIDE Find NEEDS POSITION")
	for _, need := range di.Needs {
		for _, p := range need.Positions {
			if pos.Line == p.Line && pos.Character >= p.StartCol && pos.Character <= p.EndCol {
				return need.Need, nil
			}
		}
	}
	return Need{}, errors.New("could not find a known need at requested position in document")
}
