package lsp

// InitializeParams representa os parâmetros de inicialização do LSP
type InitializeParams struct {
	ProcessID             int                `json:"processId"`
	RootURI              string             `json:"rootUri"`
	InitializationOptions map[string]interface{} `json:"initializationOptions,omitempty"`
}

// InitializeResult representa o resultado da inicialização do LSP
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

// TextDocumentItem representa um documento de texto
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// DidOpenTextDocumentParams representa os parâmetros do evento didOpen
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// TextDocumentIdentifier identifica um documento de texto
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// VersionedTextDocumentIdentifier identifica um documento de texto com versão
type VersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

// TextDocumentContentChangeEvent representa uma mudança no conteúdo do documento
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidChangeTextDocumentParams representa os parâmetros do evento didChange
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// Position representa uma posição no documento
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range representa um intervalo no documento
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location representa uma localização no documento
type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// DocumentSymbolParams representa os parâmetros para símbolos do documento
type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DocumentSymbol representa um símbolo no documento
type DocumentSymbol struct {
	Name           string           `json:"name"`
	Detail         string           `json:"detail,omitempty"`
	Kind           int             `json:"kind"`
	Range          Range           `json:"range"`
	SelectionRange Range           `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

// CompletionParams representa os parâmetros para completação
type CompletionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position              `json:"position"`
	Context      CompletionContext     `json:"context,omitempty"`
}

// CompletionContext representa o contexto da completação
type CompletionContext struct {
	TriggerKind      int    `json:"triggerKind"`
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
}

// CompletionItem representa um item de completação
type CompletionItem struct {
	Label         string `json:"label"`
	Kind         int    `json:"kind"`
	Detail       string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
	InsertText   string `json:"insertText,omitempty"`
	SortText     string `json:"sortText,omitempty"`
}

// HoverParams representa os parâmetros para hover
type HoverParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position              `json:"position"`
}

// MarkupContent representa conteúdo formatado
type MarkupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

// Hover representa informações de hover
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range       `json:"range,omitempty"`
}

// DefinitionParams representa os parâmetros para definição
type DefinitionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position              `json:"position"`
}

// DocumentFormattingParams representa os parâmetros para formatação
type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions     `json:"options"`
}

// FormattingOptions representa opções de formatação
type FormattingOptions struct {
	TabSize      int  `json:"tabSize"`
	InsertSpaces bool `json:"insertSpaces"`
}

// TextEdit representa uma edição de texto
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// Diagnostic representa um diagnóstico
type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity,omitempty"`
	Code     string `json:"code,omitempty"`
	Source   string `json:"source,omitempty"`
	Message  string `json:"message"`
}
