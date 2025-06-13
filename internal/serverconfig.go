package internal

type ServerConfig struct {
	NeedsJsonPath    string `json:"needsJsonPath"`
	DocumentRootPath string `json:"documentRootPath"`
	Enabled          bool   `json:"enabled"`
}
