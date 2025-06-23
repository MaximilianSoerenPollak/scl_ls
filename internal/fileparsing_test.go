package internal

import (
	"log"
	"os"
	"sclls/lsp"
	"testing"
)

func TestGetDocumentNameFromURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    string
		wantErr bool
	}{
		{
			name:    "valid file URI",
			uri:     "file:///home/user/project/main.go",
			want:    "main.go",
			wantErr: false,
		},
		{
			name:    "valid file URI with spaces",
			uri:     "file:///home/user/my%20project/test.go",
			want:    "test.go",
			wantErr: false,
		},
		{
			name:    "invalid URI",
			uri:     "not a valid uri ://",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty URI",
			uri:     "",
			want:    ".",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDocumentNameFromURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDocumentNameFromURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDocumentNameFromURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindNeedPosition(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		needID  string
		want    []NeedPositionInfo
	}{
		{
			name:    "single match",
			content: []byte("hello world\nfind me here"),
			needID:  "find",
			want: []NeedPositionInfo{
				{Line: 1, StartCol: 0, EndCol: 4},
			},
		},
		{
			name:    "multiple matches",
			content: []byte("test test\nmore test"),
			needID:  "test",
			want: []NeedPositionInfo{
				{Line: 0, StartCol: 0, EndCol: 4},
				{Line: 0, StartCol: 5, EndCol: 9},
				{Line: 1, StartCol: 5, EndCol: 9},
			},
		},
		{
			name:    "no matches",
			content: []byte("hello world"),
			needID:  "missing",
			want:    []NeedPositionInfo{},
		},
		{
			name:    "empty content",
			content: []byte(""),
			needID:  "anything",
			want:    []NeedPositionInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindNeedPoisiton(tt.content, tt.needID)
			if len(got) != len(tt.want) {
				t.Errorf("FindNeedPoisiton() returned %d results, want %d", len(got), len(tt.want))
				return
			}
			for i, pos := range got {
				if pos.Line != tt.want[i].Line || pos.StartCol != tt.want[i].StartCol || pos.EndCol != tt.want[i].EndCol {
					t.Errorf("FindNeedPoisiton()[%d] = %+v, want %+v", i, pos, tt.want[i])
				}
			}
		})
	}
}

func TestGetLineAndColumn(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		position int
		wantLine int
		wantCol  int
	}{
		{
			name:     "first character",
			content:  []byte("hello world"),
			position: 0,
			wantLine: 0,
			wantCol:  0,
		},
		{
			name:     "after newline",
			content:  []byte("hello\nworld"),
			position: 6,
			wantLine: 1,
			wantCol:  0,
		},
		{
			name:     "position beyond content",
			content:  []byte("short"),
			position: 100,
			wantLine: 0,
			wantCol:  5,
		},
		{
			name:     "negative position",
			content:  []byte("hello"),
			position: -1,
			wantLine: 0,
			wantCol:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLine, gotCol := getLineAndColumn(tt.content, tt.position)
			if gotLine != tt.wantLine {
				t.Errorf("getLineAndColumn() line = %v, want %v", gotLine, tt.wantLine)
			}
			if gotCol != tt.wantCol {
				t.Errorf("getLineAndColumn() col = %v, want %v", gotCol, tt.wantCol)
			}
		})
	}
}

func TestDocumentInfo_FindNeedsInPosition(t *testing.T) {
	// Setup test data
	testNeed := Need{ID: "test-need", Title: "Test Need"}

	di := DocumentInfo{
		DocumentNeeds: DocumentNeeds{
			Needs: []NeedDocInfo{
				{
					Positions: []NeedPositionInfo{
						{Line: 1, StartCol: 5, EndCol: 10},
						{Line: 3, StartCol: 0, EndCol: 4},
					},
					Need: testNeed,
				},
			},
		},
	}

	tests := []struct {
		name    string
		pos     lsp.Position
		want    Need
		wantErr bool
	}{
		{
			name:    "position matches first range",
			pos:     lsp.Position{Line: 1, Character: 7},
			want:    testNeed,
			wantErr: false,
		},
		{
			name:    "position matches second range",
			pos:     lsp.Position{Line: 3, Character: 2},
			want:    testNeed,
			wantErr: false,
		},
		{
			name:    "position outside any range",
			pos:     lsp.Position{Line: 0, Character: 0},
			want:    Need{},
			wantErr: true,
		},
		{
			name:    "position on wrong line",
			pos:     lsp.Position{Line: 2, Character: 7},
			want:    Need{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := di.FindNeedsInPosition(tt.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindNeedsInPosition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.want.ID || got.Title != tt.want.Title {
				t.Errorf("FindNeedsInPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDocumentNeeds(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)

	tests := []struct {
		name string
		uri  string
		want string // expected DocName
	}{
		{
			name: "valid URI",
			uri:  "file:///home/user/test.go",
			want: "test.go",
		},
		{
			name: "invalid URI still creates DocumentNeeds",
			uri:  "invalid://uri",
			want: ".", // empty because parsing fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDocumentNeeds(tt.uri, logger)
			if got.URI != tt.uri {
				t.Errorf("NewDocumentNeeds() URI = %v, want %v", got.URI, tt.uri)
			}
			if got.DocName != tt.want {
				t.Errorf("NewDocumentNeeds() DocName = %v, want %v", got.DocName, tt.want)
			}
		})
	}
}
