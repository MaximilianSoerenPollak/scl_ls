package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	// we specify the params later
	// params..
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"`

	//
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}

