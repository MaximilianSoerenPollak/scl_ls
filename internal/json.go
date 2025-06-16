package internal

import (
	"encoding/json"
	"strings"
)

// Handling custom JSON parsing here

// StringSlice is a custom type that can unmarshal from a comma-separated string or JSON array
type StringSlice []string

func (ss *StringSlice) UnmarshalJSON(data []byte) error {
	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*ss = StringSlice(slice)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" {
		*ss = StringSlice{}
		return nil
	}

	parts := strings.Split(str, ",")
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimSpace(part)
	}

	*ss = StringSlice(result)
	return nil
}
