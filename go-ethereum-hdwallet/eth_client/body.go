package eth

import (
	"encoding/json"
)

type (
	// Body is a struct used as the body of a POST HTTP JSON-RPC request.
	Body struct {
		ID      int64         `json:"id"`
		Method  string        `json:"method"`
		Jsonrpc string        `json:"jsonrpc"`
		Params  []interface{} `json:"params"`
	}
)

const (
	JsonrpcVersion = "2.0"
)

// NewBodyWithParameters creates a new Body struct, using the provided parameters slice.
func NewBody(method string, bodyid int64, parameters []interface{}) ([]byte, error) {
	if len(parameters) == 0 {
		parameters = []interface{}{}
	}

	body := Body{
		ID:      bodyid,
		Method:  method,
		Params:  parameters,
		Jsonrpc: JsonrpcVersion,
	}

	return json.Marshal(body)
}
