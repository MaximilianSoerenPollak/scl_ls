package internal

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringSlice_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    StringSlice
		wantErr bool
	}{
		{
			name:    "JSON array",
			input:   []byte(`["apple", "banana", "cherry"]`),
			want:    StringSlice{"apple", "banana", "cherry"},
			wantErr: false,
		},
		{
			name:    "comma-separated string",
			input:   []byte(`"apple,banana,cherry"`),
			want:    StringSlice{"apple", "banana", "cherry"},
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   []byte(`""`),
			want:    StringSlice{},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid json`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ss StringSlice
			err := ss.UnmarshalJSON(tt.input)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("StringSlice.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !reflect.DeepEqual(ss, tt.want) {
				t.Errorf("StringSlice.UnmarshalJSON() = %v, want %v", ss, tt.want)
			}
		})
	}
}

func TestStringSlice_UnmarshalJSON_WithSpaces(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  StringSlice
	}{
		{
			name:  "string with spaces",
			input: []byte(`"apple, banana , cherry"`),
			want:  StringSlice{"apple", "banana", "cherry"},
		},
		{
			name:  "single item with spaces",
			input: []byte(`"  single item  "`),
			want:  StringSlice{"single item"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ss StringSlice
			err := ss.UnmarshalJSON(tt.input)
			
			if err != nil {
				t.Errorf("StringSlice.UnmarshalJSON() unexpected error = %v", err)
				return
			}
			
			if !reflect.DeepEqual(ss, tt.want) {
				t.Errorf("StringSlice.UnmarshalJSON() = %v, want %v", ss, tt.want)
			}
		})
	}
}

// Test integration with actual JSON unmarshaling
func TestStringSlice_WithJSONUnmarshal(t *testing.T) {
	type TestStruct struct {
		Items StringSlice `json:"items"`
	}

	tests := []struct {
		name    string
		input   string
		want    StringSlice
		wantErr bool
	}{
		{
			name:    "JSON with array",
			input:   `{"items": ["a", "b", "c"]}`,
			want:    StringSlice{"a", "b", "c"},
			wantErr: false,
		},
		{
			name:    "JSON with string",
			input:   `{"items": "a,b,c"}`,
			want:    StringSlice{"a", "b", "c"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts TestStruct
			err := json.Unmarshal([]byte(tt.input), &ts)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !reflect.DeepEqual(ts.Items, tt.want) {
				t.Errorf("json.Unmarshal() items = %v, want %v", ts.Items, tt.want)
			}
		})
	}
}
