package internal

import (
	"log"
	"os"
	"sclls/lsp"
	"testing"
)



// Helper function to create a test state
func createTestState() State {
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)
	config := ServerConfig{
		NeedsJsonPath:    "/test/needs.json",
		DocumentRootPath: "/test/docs",
		TemplateStrings:  []string{"# req-Id: ", "# req-traceability: "},
	}

	// Create mock needs list
	needsList := NeedsInfo{
		"REQ_001":  Need{ID: "REQ_001", Docname: "requirements", Lineno: 10},
		"REQ_002":  Need{ID: "REQ_002", Docname: "design", Lineno: 20},
		"TOOL_001": Need{ID: "TOOL_001", Docname: "tools", Lineno: 5},
	}

	return State{
		Documents:    make(map[string]*DocumentInfo),
		NeedsList:    needsList,
		ServerConfig: config,
		Logger:       logger,
	}
}

// Tests for NewState
func TestNewState(t *testing.T) {
	tests := []struct {
		name        string
		config      ServerConfig
		expectPanic bool
	}{
		{
			name: "valid config",
			config: ServerConfig{
				NeedsJsonPath:    "/test/needs.json",
				DocumentRootPath: "/test/docs",
				TemplateStrings:  []string{"# req-Id: "},
			},
			expectPanic: false,
		},
		{
			name: "empty config",
			config: ServerConfig{
				NeedsJsonPath:    "",
				DocumentRootPath: "",
				TemplateStrings:  []string{},
			},
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic but didn't get one")
					}
				}()
			}

			state := NewState(tt.config, logger)

			if state.Documents == nil {
				t.Error("Expected Documents map to be initialized")
			}
			if state.Logger == nil {
				t.Error("Expected Logger to be set")
			}
		})
	}
}

// Tests for OpenDocument
func TestOpenDocument(t *testing.T) {
	tests := []struct {
		name          string
		uri           string
		content       string
		expectedDiags int
		shouldHaveDoc bool
	}{
		{
			name:          "new document with valid content",
			uri:           "file:///test.rst",
			content:       "# req-Id: REQ_001\nSome content",
			expectedDiags: 0,
			shouldHaveDoc: true,
		},
		{
			name:          "new document with invalid need",
			uri:           "file:///test2.rst",
			content:       "# req-Id: INVALID_NEED\nSome content",
			expectedDiags: 1,
			shouldHaveDoc: true,
		},
		{
			name:          "document with empty template",
			uri:           "file:///test3.rst",
			content:       "# req-Id: \nSome content",
			expectedDiags: 1,
			shouldHaveDoc: true,
		},
		{
			name:          "document without templates",
			uri:           "file:///test4.rst",
			content:       "Just some regular content\nNo templates here",
			expectedDiags: 0,
			shouldHaveDoc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := createTestState()

			diagnostics := state.OpenDocument(tt.uri, tt.content)

			if len(diagnostics) != tt.expectedDiags {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiags, len(diagnostics))
			}

			if tt.shouldHaveDoc {
				if _, exists := state.Documents[tt.uri]; !exists {
					t.Error("Expected document to be stored in state")
				}
				if state.Documents[tt.uri].Content != tt.content {
					t.Error("Document content doesn't match expected")
				}
			}
		})
	}
}

// Tests for UpdateDocument
func TestUpdateDocument(t *testing.T) {
	tests := []struct {
		name           string
		uri            string
		initialContent string
		updatedContent string
		expectedDiags  int
	}{
		{
			name:           "update existing document",
			uri:            "file:///existing.rst",
			initialContent: "# req-Id: REQ_001",
			updatedContent: "# req-Id: REQ_002",
			expectedDiags:  0,
		},
		{
			name:           "update non-existing document",
			uri:            "file:///new.rst",
			initialContent: "",
			updatedContent: "# req-Id: INVALID",
			expectedDiags:  1,
		},
		{
			name:           "update with multiple needs",
			uri:            "file:///multi.rst",
			initialContent: "# req-Id: REQ_001",
			updatedContent: "# req-Id: REQ_001, REQ_002, INVALID",
			expectedDiags:  1,
		},
		{
			name:           "update with empty content",
			uri:            "file:///empty.rst",
			initialContent: "# req-Id: REQ_001",
			updatedContent: "",
			expectedDiags:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := createTestState()

			// Setup initial document if needed
			if tt.initialContent != "" {
				state.OpenDocument(tt.uri, tt.initialContent)
			}

			diagnostics := state.UpdateDocument(tt.uri, tt.updatedContent)

			if len(diagnostics) != tt.expectedDiags {
				t.Errorf("Expected %d diagnostics, got %d", tt.expectedDiags, len(diagnostics))
			}

			if doc, exists := state.Documents[tt.uri]; exists {
				if doc.Content != tt.updatedContent {
					t.Error("Document content wasn't updated correctly")
				}
			} else {
				t.Error("Document should exist after update")
			}
		})
	}
}

// Tests for FindDiagnosticsInDocument
func TestFindDiagnosticsInDocument(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectedNum int
		description string
	}{
		{
			name:        "valid needs",
			content:     "# req-Id: REQ_001\n# req-traceability: REQ_002",
			expectedNum: 0,
			description: "Should have no diagnostics for valid needs",
		},
		{
			name:        "invalid needs",
			content:     "# req-Id: INVALID_NEED\n# req-traceability: ANOTHER_INVALID",
			expectedNum: 2,
			description: "Should have diagnostics for invalid needs",
		},
		{
			name:        "empty template",
			content:     "# req-Id: \n# req-traceability: ",
			expectedNum: 2,
			description: "Should have diagnostics for empty templates",
		},
		{
			name:        "mixed valid and invalid",
			content:     "# req-Id: REQ_001, INVALID, REQ_002",
			expectedNum: 1,
			description: "Should have diagnostic only for invalid need",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := createTestState()
			diagnostics := state.FindDiagnosticsInDocument([]byte(tt.content))

			if len(diagnostics) != tt.expectedNum {
				t.Errorf("%s: Expected %d diagnostics, got %d", tt.description, tt.expectedNum, len(diagnostics))
			}

			// Verify diagnostic properties
			for _, diag := range diagnostics {
				if diag.Source != "scl_lsp" {
					t.Error("Expected diagnostic source to be 'scl_lsp'")
				}
				if diag.Severity < 1 || diag.Severity > 2 {
					t.Error("Expected diagnostic severity to be 1 or 2")
				}
			}
		})
	}
}

// Tests for TextDocumentCompletion
func TestTextDocumentCompletion(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		position      lsp.Position
		expectedItems int
		description   string
	}{
		{
			name:          "completion after req- prefix",
			content:       "# req-",
			position:      lsp.Position{Line: 0, Character: 6},
			expectedItems: 2, // req-Id: and req-traceability:
			description:   "Should suggest template completions",
		},
		{
			name:          "completion after req-Id:",
			content:       "# req-Id: ",
			position:      lsp.Position{Line: 0, Character: 10},
			expectedItems: 3, // All needs in mock data
			description:   "Should suggest all available needs",
		},
		{
			name:          "completion with partial need",
			content:       "# req-Id: REQ",
			position:      lsp.Position{Line: 0, Character: 13},
			expectedItems: 2, // REQ_001 and REQ_002
			description:   "Should filter needs by prefix",
		},
		{
			name:          "no completion for regular text",
			content:       "Just regular text",
			position:      lsp.Position{Line: 0, Character: 10},
			expectedItems: 0,
			description:   "Should not suggest completions for regular text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := createTestState()
			uri := "file:///test.rst"

			// Setup document
			state.OpenDocument(uri, tt.content)

			response := state.TextDocumentCompletion(1, uri, tt.position)

			if len(response.Result) != tt.expectedItems {
				t.Errorf("%s: Expected %d completion items, got %d", tt.description, tt.expectedItems, len(response.Result))
			}

			// Verify response structure
			if response.Response.RPC != "2.0" {
				t.Error("Expected RPC version 2.0")
			}
			if response.Response.ID == nil || *response.Response.ID != 1 {
				t.Error("Expected response ID to be 1")
			}
		})
	}
}

// Tests for GoToDefinition
func TestGoToDefinition(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		position        lsp.Position
		expectedResults int
		description     string
	}{
		{
			name:            "valid need position",
			content:         "# req-Id: REQ_001",
			position:        lsp.Position{Line: 0, Character: 12}, // Position on REQ_001
			expectedResults: 1,
			description:     "Should find definition for valid need",
		},
		{
			name:            "invalid position",
			content:         "# req-Id: REQ_001",
			position:        lsp.Position{Line: 0, Character: 5}, // Position not on need
			expectedResults: 0,
			description:     "Should not find definition for invalid position",
		},
		{
			name:            "empty document",
			content:         "",
			position:        lsp.Position{Line: 0, Character: 0},
			expectedResults: 0,
			description:     "Should handle empty document gracefully",
		},
		{
			name:            "position beyond document",
			content:         "# req-Id: REQ_001",
			position:        lsp.Position{Line: 10, Character: 0},
			expectedResults: 0,
			description:     "Should handle out-of-bounds position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := createTestState()
			uri := "file:///test.rst"

			// Setup document
			state.OpenDocument(uri, tt.content)

			response := state.GoToDefinition(1, uri, tt.position)

			if len(response.Result) != tt.expectedResults {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedResults, len(response.Result))
			}

			// Verify response structure
			if response.Response.RPC != "2.0" {
				t.Error("Expected RPC version 2.0")
			}
			if response.Response.ID == nil || *response.Response.ID != 1 {
				t.Error("Expected response ID to be 1")
			}

			// If we have results, verify they have proper structure
			for _, result := range response.Result {
				if result.URI == "" {
					t.Error("Expected result to have URI")
				}
				if result.Range.Start.Line < 0 {
					t.Error("Expected valid line number in result")
				}
			}
		})
	}
}
