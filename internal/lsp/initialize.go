package lsp

import "go.lsp.dev/protocol"

type InitializeRequest struct {
	Request
	Params protocol.InitializeParams `json:"params"`
}

type InitializeResponse struct {
	Response
	Result protocol.InitializeResult `json:"result"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: RPCVersion,
			ID:  &id,
		},
		Result: protocol.InitializeResult{
			Capabilities: protocol.ServerCapabilities{
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CompletionProvider: &protocol.CompletionOptions{},
			},
			ServerInfo: &protocol.ServerInfo{
				Name:    name,
				Version: version,
			},
		},
	}
}
