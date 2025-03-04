package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/peder1981/advpl-tlpp-compiler/pkg/ast"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/compiler"
	"github.com/peder1981/advpl-tlpp-compiler/pkg/parser"
)

// Server implementa o servidor LSP para AdvPL/TLPP
type Server struct {
	documents     sync.Map
	capabilities  ServerCapabilities
	compiler      *compiler.Compiler
	ide           *compiler.IDEIntegration
	configuration Configuration
	client        *Client
}

// Configuration representa as configurações do servidor LSP
type Configuration struct {
	MaxNumberOfProblems int    `json:"maxNumberOfProblems"`
	Dialect            string `json:"dialect"`
	LogLevel           string `json:"logLevel"`
}

// ServerCapabilities define as capacidades do servidor LSP
type ServerCapabilities struct {
	TextDocumentSync           int  `json:"textDocumentSync"`
	DocumentFormattingProvider bool `json:"documentFormattingProvider"`
	DocumentSymbolProvider     bool `json:"documentSymbolProvider"`
	CompletionProvider        struct {
		TriggerCharacters []string `json:"triggerCharacters"`
	} `json:"completionProvider"`
	DefinitionProvider bool `json:"definitionProvider"`
	HoverProvider      bool `json:"hoverProvider"`
	DiagnosticProvider bool `json:"diagnosticProvider"`
}

// Client representa o cliente LSP
type Client struct {
	notifyDiagnostics func(uri string, diagnostics []Diagnostic) error
}

// NewServer cria uma nova instância do servidor LSP
func NewServer() *Server {
	return &Server{
		capabilities: ServerCapabilities{
			TextDocumentSync:           1, // Incremental
			DocumentFormattingProvider: true,
			DocumentSymbolProvider:     true,
			CompletionProvider: struct {
				TriggerCharacters []string `json:"triggerCharacters"`
			}{
				TriggerCharacters: []string{".", ":"},
			},
			DefinitionProvider: true,
			HoverProvider:      true,
			DiagnosticProvider: true,
		},
		configuration: Configuration{
			MaxNumberOfProblems: 100,
			Dialect:            "advpl",
			LogLevel:           "info",
		},
		client: &Client{},
	}
}

// SetClient configura o cliente LSP
func (s *Server) SetClient(notifyDiagnostics func(uri string, diagnostics []Diagnostic) error) {
	s.client.notifyDiagnostics = notifyDiagnostics
}

// Initialize inicializa o servidor LSP
func (s *Server) Initialize(ctx context.Context, params InitializeParams) (InitializeResult, error) {
	log.Printf("Inicializando servidor LSP com configurações: %+v", s.configuration)
	return InitializeResult{
		Capabilities: s.capabilities,
	}, nil
}

// DidOpen é chamado quando um documento é aberto
func (s *Server) DidOpen(ctx context.Context, params DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI
	text := params.TextDocument.Text

	log.Printf("Documento aberto: %s", uri)
	
	// Armazena o documento
	s.documents.Store(uri, text)

	// Analisa o documento
	return s.analyzeDocument(uri, text)
}

// DidChange é chamado quando um documento é modificado
func (s *Server) DidChange(ctx context.Context, params DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI
	
	// Aplica as mudanças ao documento
	if text, ok := s.documents.Load(uri); ok {
		newText := applyChanges(text.(string), params.ContentChanges)
		s.documents.Store(uri, newText)
		return s.analyzeDocument(uri, newText)
	}
	
	return fmt.Errorf("documento não encontrado: %s", uri)
}

// analyzeDocument analisa um documento e publica diagnósticos
func (s *Server) analyzeDocument(uri string, text string) error {
	// Cria uma nova instância do IDE Integration
	s.ide = compiler.NewIDEIntegration(uri)

	// Parse o documento
	program, err := parser.ParseSource(text)
	if err != nil {
		// Mesmo com erro, publicamos os diagnósticos disponíveis
		diagnostics := []Diagnostic{
			{
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 0, Character: 0},
				},
				Severity: 1, // Error
				Source:   "advpl-compiler",
				Message:  fmt.Sprintf("Erro de parse: %v", err),
			},
		}
		return s.publishDiagnostics(uri, diagnostics)
	}

	// Processa o programa com o IDE Integration
	s.ide.ProcessProgram(program)

	// Obtém os diagnósticos
	diagnostics, err := s.ide.GetDiagnostics()
	if err != nil {
		return err
	}

	// Publica os diagnósticos
	return s.publishDiagnostics(uri, parseDiagnostics(diagnostics))
}

// DocumentSymbol retorna os símbolos do documento
func (s *Server) DocumentSymbol(ctx context.Context, params DocumentSymbolParams) ([]DocumentSymbol, error) {
	if text, ok := s.documents.Load(params.TextDocument.URI); ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text.(string))
		if err != nil {
			return nil, err
		}

		s.ide.ProcessProgram(program)
		symbols, err := s.ide.GetSymbols()
		if err != nil {
			return nil, err
		}

		return parseSymbols(symbols), nil
	}

	return nil, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
}

// Completion fornece sugestões de completação
func (s *Server) Completion(ctx context.Context, params CompletionParams) ([]CompletionItem, error) {
	if text, ok := s.documents.Load(params.TextDocument.URI); ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text.(string))
		if err != nil {
			return nil, err
		}

		s.ide.ProcessProgram(program)
		
		// Obter o contexto de completação baseado na posição do cursor
		context := ""
		if params.Context.TriggerCharacter != "" {
			context = params.Context.TriggerCharacter
		}
		
		completions, err := s.ide.GetCompletions(context)
		if err != nil {
			return nil, err
		}

		return parseCompletions(completions), nil
	}

	return nil, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
}

// Hover fornece informações ao passar o mouse sobre um símbolo
func (s *Server) Hover(ctx context.Context, params HoverParams) (Hover, error) {
	if text, ok := s.documents.Load(params.TextDocument.URI); ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text.(string))
		if err != nil {
			return Hover{}, err
		}

		s.ide.ProcessProgram(program)
		
		// Obter informações do símbolo na posição do cursor
		symbolInfo, err := s.ide.GetSymbolAtPosition(params.Position.Line, params.Position.Character)
		if err != nil || symbolInfo == "" {
			return Hover{
				Contents: MarkupContent{
					Kind:  "markdown",
					Value: "Nenhuma informação disponível",
				},
			}, nil
		}
		
		// Parsear as informações do símbolo
		var info struct {
			Name        string `json:"name"`
			Kind        int    `json:"kind"`
			Description string `json:"description"`
			Type        string `json:"type"`
			Location    string `json:"location"`
		}
		
		if err := json.Unmarshal([]byte(symbolInfo), &info); err != nil {
			return Hover{}, err
		}
		
		// Formatar as informações para exibição
		content := fmt.Sprintf("## %s\n\n", info.Name)
		content += fmt.Sprintf("**Tipo:** %s\n\n", getSymbolKindName(info.Kind))
		
		if info.Type != "" {
			content += fmt.Sprintf("**Tipo de dado:** %s\n\n", info.Type)
		}
		
		if info.Description != "" {
			content += fmt.Sprintf("**Descrição:** %s\n\n", info.Description)
		}
		
		if info.Location != "" {
			content += fmt.Sprintf("**Definido em:** %s\n", info.Location)
		}
		
		return Hover{
			Contents: MarkupContent{
				Kind:  "markdown",
				Value: content,
			},
			Range: &Range{
				Start: params.Position,
				End:   Position{Line: params.Position.Line, Character: params.Position.Character + len(info.Name)},
			},
		}, nil
	}

	return Hover{}, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
}

// Definition retorna a definição de um símbolo
func (s *Server) Definition(ctx context.Context, params DefinitionParams) ([]Location, error) {
	if text, ok := s.documents.Load(params.TextDocument.URI); ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text.(string))
		if err != nil {
			return nil, err
		}

		s.ide.ProcessProgram(program)
		
		// Obter a definição do símbolo na posição do cursor
		definitionInfo, err := s.ide.GetDefinitionAtPosition(params.Position.Line, params.Position.Character)
		if err != nil || definitionInfo == "" {
			return []Location{}, nil
		}
		
		// Parsear as informações da definição
		var definitions []struct {
			URI   string `json:"uri"`
			Line  int    `json:"line"`
			Col   int    `json:"column"`
			EndLine int  `json:"endLine"`
			EndCol  int  `json:"endColumn"`
		}
		
		if err := json.Unmarshal([]byte(definitionInfo), &definitions); err != nil {
			return nil, err
		}
		
		// Converter para o formato LSP
		locations := make([]Location, len(definitions))
		for i, def := range definitions {
			locations[i] = Location{
				URI: def.URI,
				Range: Range{
					Start: Position{Line: def.Line, Character: def.Col},
					End:   Position{Line: def.EndLine, Character: def.EndCol},
				},
			}
		}
		
		return locations, nil
	}

	return nil, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
}

// Formatting formata um documento
func (s *Server) Formatting(ctx context.Context, params DocumentFormattingParams) ([]TextEdit, error) {
	if text, ok := s.documents.Load(params.TextDocument.URI); ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		// Parse o documento
		program, err := parser.ParseSource(text.(string))
		if err != nil {
			return nil, err
		}
		
		// Formatar o código
		formattedCode, err := s.ide.FormatCode(text.(string), params.Options.TabSize, params.Options.InsertSpaces)
		if err != nil {
			return nil, err
		}
		
		// Criar um TextEdit para substituir todo o documento
		return []TextEdit{
			{
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 999999, Character: 999999}, // Fim do documento
				},
				NewText: formattedCode,
			},
		}, nil
	}
	
	return nil, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
}

// publishDiagnostics publica diagnósticos para um documento
func (s *Server) publishDiagnostics(uri string, diagnostics []Diagnostic) error {
	if s.client.notifyDiagnostics != nil {
		return s.client.notifyDiagnostics(uri, diagnostics)
	}
	
	log.Printf("Diagnósticos para %s: %d problemas encontrados", uri, len(diagnostics))
	return nil
}

// getSymbolKindName retorna o nome do tipo de símbolo
func getSymbolKindName(kind int) string {
	switch kind {
	case 0:
		return "Função"
	case 1:
		return "Classe"
	case 2:
		return "Método"
	case 3:
		return "Variável"
	case 4:
		return "Parâmetro"
	case 5:
		return "Atributo"
	default:
		return "Desconhecido"
	}
}

// parseDiagnostics converte os diagnósticos do IDE Integration para o formato LSP
func parseDiagnostics(diagnosticsJSON string) []Diagnostic {
	var diagnostics []Diagnostic
	var ideDiagnostics []struct {
		Line     int    `json:"line"`
		Column   int    `json:"column"`
		EndLine  int    `json:"endLine"`
		EndColumn int   `json:"endColumn"`
		Severity int    `json:"severity"`
		Message  string `json:"message"`
		Code     string `json:"code"`
		Source   string `json:"source"`
	}

	if err := json.Unmarshal([]byte(diagnosticsJSON), &ideDiagnostics); err != nil {
		return diagnostics
	}

	for _, d := range ideDiagnostics {
		diagnostic := Diagnostic{
			Range: Range{
				Start: Position{Line: d.Line, Character: d.Column},
				End:   Position{Line: d.EndLine, Character: d.EndColumn},
			},
			Severity: d.Severity,
			Code:     d.Code,
			Source:   d.Source,
			Message:  d.Message,
		}
		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}
