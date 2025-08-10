package lsp

type LSP struct {
	workspaceRoot string
}

func New(workspaceRoot string) *LSP {
	return &LSP{
		workspaceRoot: workspaceRoot,
	}
}
