package lsp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/ast"
)

func TestInitialize(t *testing.T) {
	server := NewServer()
	
	// Cria uma requisição de inicialização
	params := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	result, err := server.Initialize(context.Background(), &params)
	
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	if result == nil {
		t.Fatal("Resultado da inicialização é nulo")
	}
	
	// Verifica se o servidor retornou as capacidades esperadas
	if !result.Capabilities.TextDocumentSync.OpenClose {
		t.Error("Capacidade OpenClose não está habilitada")
	}
	
	if !result.Capabilities.CompletionProvider.ResolveProvider {
		t.Error("CompletionProvider.ResolveProvider não está habilitado")
	}
	
	if !result.Capabilities.HoverProvider {
		t.Error("HoverProvider não está habilitado")
	}
	
	if !result.Capabilities.DefinitionProvider {
		t.Error("DefinitionProvider não está habilitado")
	}
	
	if !result.Capabilities.DocumentSymbolProvider {
		t.Error("DocumentSymbolProvider não está habilitado")
	}
	
	if !result.Capabilities.DocumentFormattingProvider {
		t.Error("DocumentFormattingProvider não está habilitado")
	}
}

func TestDidOpenTextDocument(t *testing.T) {
	server := NewServer()
	
	// Inicializa o servidor
	initParams := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	_, err := server.Initialize(context.Background(), &initParams)
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	// Simula a abertura de um documento
	params := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test/workspace/test.prw",
			LanguageID: "advpl",
			Version:    1,
			Text:       "Function Test()\nReturn\nEndFunction",
		},
	}
	
	err = server.DidOpenTextDocument(context.Background(), &params)
	
	if err != nil {
		t.Fatalf("Erro ao abrir documento: %v", err)
	}
	
	// Verifica se o documento foi adicionado ao gerenciador de documentos
	doc, exists := server.documents[params.TextDocument.URI]
	
	if !exists {
		t.Fatal("Documento não foi adicionado ao gerenciador de documentos")
	}
	
	if doc.Text != params.TextDocument.Text {
		t.Errorf("Texto do documento não corresponde. Esperado: %s, Obtido: %s", 
			params.TextDocument.Text, doc.Text)
	}
}

func TestCompletion(t *testing.T) {
	server := NewServer()
	
	// Inicializa o servidor
	initParams := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	_, err := server.Initialize(context.Background(), &initParams)
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	// Adiciona um documento
	didOpenParams := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test/workspace/test.prw",
			LanguageID: "advpl",
			Version:    1,
			Text:       "Function Test()\nLocal a := 5\nReturn\nEndFunction",
		},
	}
	
	err = server.DidOpenTextDocument(context.Background(), &didOpenParams)
	if err != nil {
		t.Fatalf("Erro ao abrir documento: %v", err)
	}
	
	// Solicita completions
	completionParams := CompletionParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///test/workspace/test.prw",
			},
			Position: Position{
				Line:      1,
				Character: 7, // Após "Local "
			},
		},
	}
	
	completions, err := server.Completion(context.Background(), &completionParams)
	
	if err != nil {
		t.Fatalf("Erro ao obter completions: %v", err)
	}
	
	// Verifica se retornou alguma sugestão
	if len(completions.Items) == 0 {
		t.Fatal("Nenhuma sugestão de completação retornada")
	}
}

func TestHover(t *testing.T) {
	server := NewServer()
	
	// Inicializa o servidor
	initParams := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	_, err := server.Initialize(context.Background(), &initParams)
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	// Adiciona um documento
	didOpenParams := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test/workspace/test.prw",
			LanguageID: "advpl",
			Version:    1,
			Text:       "Function Test()\nLocal a := 5\nReturn a\nEndFunction",
		},
	}
	
	err = server.DidOpenTextDocument(context.Background(), &didOpenParams)
	if err != nil {
		t.Fatalf("Erro ao abrir documento: %v", err)
	}
	
	// Solicita hover
	hoverParams := HoverParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///test/workspace/test.prw",
			},
			Position: Position{
				Line:      2,
				Character: 8, // Sobre a variável "a" no Return
			},
		},
	}
	
	hover, err := server.Hover(context.Background(), &hoverParams)
	
	if err != nil {
		t.Fatalf("Erro ao obter hover: %v", err)
	}
	
	// Verifica se retornou alguma informação
	if hover == nil || hover.Contents.Value == "" {
		t.Fatal("Nenhuma informação de hover retornada")
	}
}

func TestDefinition(t *testing.T) {
	server := NewServer()
	
	// Inicializa o servidor
	initParams := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	_, err := server.Initialize(context.Background(), &initParams)
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	// Adiciona um documento
	didOpenParams := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test/workspace/test.prw",
			LanguageID: "advpl",
			Version:    1,
			Text:       "Function Test()\nLocal a := 5\nReturn a\nEndFunction",
		},
	}
	
	err = server.DidOpenTextDocument(context.Background(), &didOpenParams)
	if err != nil {
		t.Fatalf("Erro ao abrir documento: %v", err)
	}
	
	// Solicita definição
	definitionParams := DefinitionParams{
		TextDocumentPositionParams: TextDocumentPositionParams{
			TextDocument: TextDocumentIdentifier{
				URI: "file:///test/workspace/test.prw",
			},
			Position: Position{
				Line:      2,
				Character: 8, // Sobre a variável "a" no Return
			},
		},
	}
	
	locations, err := server.Definition(context.Background(), &definitionParams)
	
	if err != nil {
		t.Fatalf("Erro ao obter definição: %v", err)
	}
	
	// Verifica se retornou alguma localização
	if len(locations) == 0 {
		t.Fatal("Nenhuma localização de definição retornada")
	}
	
	// Verifica se a localização aponta para a declaração da variável
	if locations[0].Range.Start.Line != 1 {
		t.Errorf("Definição incorreta. Esperado linha 1, obtido linha %d", 
			locations[0].Range.Start.Line)
	}
}

func TestDocumentSymbols(t *testing.T) {
	server := NewServer()
	
	// Inicializa o servidor
	initParams := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	_, err := server.Initialize(context.Background(), &initParams)
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	// Adiciona um documento com função e classe
	didOpenParams := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test/workspace/test.prw",
			LanguageID: "advpl",
			Version:    1,
			Text: `Function Test()
				Return
			EndFunction
			
			Class MyClass
				Data Name
				
				Method New() Constructor
					Return Self
			EndClass`,
		},
	}
	
	err = server.DidOpenTextDocument(context.Background(), &didOpenParams)
	if err != nil {
		t.Fatalf("Erro ao abrir documento: %v", err)
	}
	
	// Solicita símbolos do documento
	symbolParams := DocumentSymbolParams{
		TextDocument: TextDocumentIdentifier{
			URI: "file:///test/workspace/test.prw",
		},
	}
	
	symbols, err := server.DocumentSymbol(context.Background(), &symbolParams)
	
	if err != nil {
		t.Fatalf("Erro ao obter símbolos: %v", err)
	}
	
	// Verifica se retornou os símbolos esperados
	if len(symbols) < 2 {
		t.Fatalf("Número incorreto de símbolos. Esperado pelo menos 2, obtido %d", len(symbols))
	}
	
	// Verifica se encontrou a função e a classe
	foundFunction := false
	foundClass := false
	
	for _, symbol := range symbols {
		if symbol.Name == "Test" && symbol.Kind == SymbolKind_Function {
			foundFunction = true
		}
		if symbol.Name == "MyClass" && symbol.Kind == SymbolKind_Class {
			foundClass = true
		}
	}
	
	if !foundFunction {
		t.Error("Função 'Test' não encontrada nos símbolos")
	}
	
	if !foundClass {
		t.Error("Classe 'MyClass' não encontrada nos símbolos")
	}
}

func TestFormatting(t *testing.T) {
	server := NewServer()
	
	// Inicializa o servidor
	initParams := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	_, err := server.Initialize(context.Background(), &initParams)
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	// Adiciona um documento mal formatado
	didOpenParams := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test/workspace/test.prw",
			LanguageID: "advpl",
			Version:    1,
			Text: `Function Test()
			Local a:=5
			Local b   :=   10
			Return a+b
			EndFunction`,
		},
	}
	
	err = server.DidOpenTextDocument(context.Background(), &didOpenParams)
	if err != nil {
		t.Fatalf("Erro ao abrir documento: %v", err)
	}
	
	// Solicita formatação
	formatParams := DocumentFormattingParams{
		TextDocument: TextDocumentIdentifier{
			URI: "file:///test/workspace/test.prw",
		},
		Options: FormattingOptions{
			TabSize:                2,
			InsertSpaces:           true,
			TrimTrailingWhitespace: true,
		},
	}
	
	edits, err := server.DocumentFormatting(context.Background(), &formatParams)
	
	if err != nil {
		t.Fatalf("Erro ao formatar documento: %v", err)
	}
	
	// Verifica se retornou edições
	if len(edits) == 0 {
		t.Fatal("Nenhuma edição de formatação retornada")
	}
}

func TestDiagnostics(t *testing.T) {
	server := NewServer()
	
	// Inicializa o servidor
	initParams := InitializeParams{
		RootURI:      "file:///test/workspace",
		Capabilities: ClientCapabilities{},
	}
	
	_, err := server.Initialize(context.Background(), &initParams)
	if err != nil {
		t.Fatalf("Erro ao inicializar servidor LSP: %v", err)
	}
	
	// Adiciona um documento com erro
	didOpenParams := DidOpenTextDocumentParams{
		TextDocument: TextDocumentItem{
			URI:        "file:///test/workspace/test.prw",
			LanguageID: "advpl",
			Version:    1,
			Text: `Function Test()
				Local a := 5
				Return a +
			EndFunction`,
		},
	}
	
	// Captura as publicações de diagnósticos
	var diagnostics *PublishDiagnosticsParams
	server.client = &mockClient{
		publishDiagnostics: func(ctx context.Context, params *PublishDiagnosticsParams) error {
			diagnostics = params
			return nil
		},
	}
	
	err = server.DidOpenTextDocument(context.Background(), &didOpenParams)
	if err != nil {
		t.Fatalf("Erro ao abrir documento: %v", err)
	}
	
	// Verifica se foram publicados diagnósticos
	if diagnostics == nil {
		t.Fatal("Nenhum diagnóstico publicado")
	}
	
	// Verifica se encontrou o erro de expressão incompleta
	if len(diagnostics.Diagnostics) == 0 {
		t.Fatal("Nenhum diagnóstico encontrado")
	}
	
	foundError := false
	for _, diag := range diagnostics.Diagnostics {
		if diag.Severity == DiagnosticSeverity_Error && diag.Range.Start.Line == 2 {
			foundError = true
			break
		}
	}
	
	if !foundError {
		t.Error("Erro de expressão incompleta não detectado")
	}
}

// Mock do cliente LSP para testes
type mockClient struct {
	publishDiagnostics func(ctx context.Context, params *PublishDiagnosticsParams) error
}

func (c *mockClient) PublishDiagnostics(ctx context.Context, params *PublishDiagnosticsParams) error {
	if c.publishDiagnostics != nil {
		return c.publishDiagnostics(ctx, params)
	}
	return nil
}

// Estruturas necessárias para os testes

type ClientCapabilities struct {
	TextDocument struct {
		Completion struct {
			CompletionItem struct {
				SnippetSupport bool `json:"snippetSupport"`
			} `json:"completionItem"`
		} `json:"completion"`
	} `json:"textDocument"`
}

type InitializeParams struct {
	RootURI      string             `json:"rootUri"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	TextDocumentSync struct {
		OpenClose bool `json:"openClose"`
		Change    int  `json:"change"`
	} `json:"textDocumentSync"`
	CompletionProvider struct {
		ResolveProvider bool     `json:"resolveProvider"`
		TriggerCharacters []string `json:"triggerCharacters"`
	} `json:"completionProvider"`
	HoverProvider            bool `json:"hoverProvider"`
	DefinitionProvider       bool `json:"definitionProvider"`
	DocumentSymbolProvider   bool `json:"documentSymbolProvider"`
	DocumentFormattingProvider bool `json:"documentFormattingProvider"`
}

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type CompletionParams struct {
	TextDocumentPositionParams
}

type CompletionItem struct {
	Label  string `json:"label"`
	Kind   int    `json:"kind"`
	Detail string `json:"detail"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type HoverParams struct {
	TextDocumentPositionParams
}

type MarkupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

type DefinitionParams struct {
	TextDocumentPositionParams
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DocumentSymbol struct {
	Name           string          `json:"name"`
	Detail         string          `json:"detail,omitempty"`
	Kind           int             `json:"kind"`
	Range          Range           `json:"range"`
	SelectionRange Range           `json:"selectionRange"`
	Children       []DocumentSymbol `json:"children,omitempty"`
}

type FormattingOptions struct {
	TabSize                int  `json:"tabSize"`
	InsertSpaces           bool `json:"insertSpaces"`
	TrimTrailingWhitespace bool `json:"trimTrailingWhitespace,omitempty"`
}

type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
}

type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// Constantes para tipos de símbolos
const (
	SymbolKind_Function = 12
	SymbolKind_Class    = 5
)

// Constantes para severidade de diagnósticos
const (
	DiagnosticSeverity_Error       = 1
	DiagnosticSeverity_Warning     = 2
	DiagnosticSeverity_Information = 3
	DiagnosticSeverity_Hint        = 4
)
