package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"advpl-tlpp-compiler/pkg/compiler"
	"advpl-tlpp-compiler/pkg/parser"
)

// Server implementa o servidor LSP para AdvPL/TLPP
type Server struct {
	documents     map[string]string
	capabilities  ServerCapabilities
	compiler      *compiler.Compiler
	ide           *compiler.IDEIntegration
	configuration Configuration
	client        *Client
	mu            sync.Mutex
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
		documents: make(map[string]string),
		capabilities: ServerCapabilities{
			TextDocumentSync:           1, // Full sync
			DocumentFormattingProvider: true,
			DocumentSymbolProvider:     true,
			CompletionProvider: struct {
				TriggerCharacters []string `json:"triggerCharacters"`
			}{
				TriggerCharacters: []string{".", ":", "$"},
			},
			DefinitionProvider: true,
			HoverProvider:      true,
			DiagnosticProvider: true,
		},
		compiler: compiler.New(nil, compiler.Options{
			Dialect: "advpl",
			Verbose: false,
		}),
		configuration: Configuration{
			MaxNumberOfProblems: 100,
			Dialect:            "advpl",
			LogLevel:           "info",
		},
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
	s.documents[uri] = text

	// Analisa o documento
	return s.analyzeDocument(uri, text)
}

// DidChange é chamado quando um documento é modificado
func (s *Server) DidChange(ctx context.Context, params DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI
	
	// Aplica as mudanças ao documento
	if text, ok := s.documents[uri]; ok {
		newText := applyChanges(text, params.ContentChanges)
		s.documents[uri] = newText
		return s.analyzeDocument(uri, newText)
	}
	
	return fmt.Errorf("documento não encontrado: %s", uri)
}

// SymbolInfo representa informações sobre um símbolo
type SymbolInfo struct {
	Name string
	Kind int
	Range Range
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
	if text, ok := s.documents[params.TextDocument.URI]; ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text)
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
	if text, ok := s.documents[params.TextDocument.URI]; ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text)
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
	if text, ok := s.documents[params.TextDocument.URI]; ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text)
		if err != nil {
			return Hover{}, err
		}

		s.ide.ProcessProgram(program)
		
		// Obter informações do símbolo na posição do cursor
		symbolInfo, err := s.getSymbolAtPosition(params.TextDocument.URI, params.Position.Line, params.Position.Character)
		if err != nil || symbolInfo == nil {
			return Hover{
				Contents: MarkupContent{
					Kind:  "markdown",
					Value: "Nenhuma informação disponível",
				},
			}, nil
		}
		
		// Formatar as informações para exibição
		content := fmt.Sprintf("## %s\n\n", symbolInfo.Name)
		content += fmt.Sprintf("**Tipo:** %s\n\n", getSymbolKindName(symbolInfo.Kind))
		
		if symbolInfo.Range.Start.Line != symbolInfo.Range.End.Line || symbolInfo.Range.Start.Character != symbolInfo.Range.End.Character {
			content += fmt.Sprintf("**Definido em:** %d:%d - %d:%d\n", symbolInfo.Range.Start.Line, symbolInfo.Range.Start.Character, symbolInfo.Range.End.Line, symbolInfo.Range.End.Character)
		}
		
		return Hover{
			Contents: MarkupContent{
				Kind:  "markdown",
				Value: content,
			},
			Range: &symbolInfo.Range,
		}, nil
	}

	return Hover{}, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
}

// Definition retorna a definição de um símbolo
func (s *Server) Definition(ctx context.Context, params DefinitionParams) ([]Location, error) {
	if text, ok := s.documents[params.TextDocument.URI]; ok {
		s.ide = compiler.NewIDEIntegration(params.TextDocument.URI)
		
		program, err := parser.ParseSource(text)
		if err != nil {
			return nil, err
		}

		s.ide.ProcessProgram(program)
		
		// Obter a definição do símbolo na posição do cursor
		definitionInfo, err := s.getDefinitionAtPosition(params.TextDocument.URI, params.Position.Line, params.Position.Character)
		if err != nil || definitionInfo == nil {
			return []Location{}, nil
		}
		
		// Converter para o formato LSP
		locations := []Location{*definitionInfo}
		
		return locations, nil
	}

	return nil, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
}

// Formatting formata um documento
func (s *Server) Formatting(ctx context.Context, params DocumentFormattingParams) ([]TextEdit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Obter o texto do documento
	doc, ok := s.documents[params.TextDocument.URI]
	if !ok {
		return nil, fmt.Errorf("documento não encontrado: %s", params.TextDocument.URI)
	}
	
	text := doc
	
	// Formatar o código
	formattedCode, err := s.formatCode(text)
	if err != nil {
		return nil, err
	}
	
	// Criar um único TextEdit que substitui todo o documento
	return []TextEdit{
		{
			Range: Range{
				Start: Position{Line: 0, Character: 0},
				End:   Position{Line: 9999, Character: 0}, // Um valor grande para cobrir todo o documento
			},
			NewText: formattedCode,
		},
	}, nil
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

// getSymbolAtPosition retorna o símbolo na posição especificada
func (s *Server) getSymbolAtPosition(uri string, line, character int) (*SymbolInfo, error) {
	// Implementação simplificada
	return &SymbolInfo{
		Name: "Symbol",
		Kind: 1,
		Range: Range{
			Start: Position{Line: line, Character: character},
			End:   Position{Line: line, Character: character + 5},
		},
	}, nil
}

// getDefinitionAtPosition retorna a definição do símbolo na posição especificada
func (s *Server) getDefinitionAtPosition(uri string, line, character int) (*Location, error) {
	// Implementação simplificada
	return &Location{
		URI: uri,
		Range: Range{
			Start: Position{Line: 0, Character: 0},
			End:   Position{Line: 5, Character: 10},
		},
	}, nil
}

// formatCode formata o código fonte
func (s *Server) formatCode(text string) (string, error) {
	// Implementação simplificada - apenas retorna o mesmo texto
	return text, nil
}
