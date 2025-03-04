package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"advpl-tlpp-compiler/pkg/lsp"
)

// Constantes para os métodos LSP
const (
	MethodInitialize           = "initialize"
	MethodInitialized          = "initialized"
	MethodShutdown             = "shutdown"
	MethodExit                 = "exit"
	MethodTextDocumentDidOpen  = "textDocument/didOpen"
	MethodTextDocumentDidChange = "textDocument/didChange"
	MethodTextDocumentDidClose = "textDocument/didClose"
	MethodTextDocumentDocumentSymbol = "textDocument/documentSymbol"
	MethodTextDocumentCompletion = "textDocument/completion"
	MethodTextDocumentHover    = "textDocument/hover"
	MethodTextDocumentDefinition = "textDocument/definition"
	MethodTextDocumentFormatting = "textDocument/formatting"
	MethodWorkspaceDidChangeConfiguration = "workspace/didChangeConfiguration"
)

// Constantes para os códigos de erro JSON-RPC
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// Estrutura para mensagens JSON-RPC
type jsonrpcMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	} `json:"error,omitempty"`
}

func main() {
	// Configura o log para um arquivo
	setupLogging()
	
	log.Println("Iniciando servidor LSP para AdvPL/TLPP")
	log.Printf("Versão Go: %s", runtime.Version())
	log.Printf("GOOS: %s, GOARCH: %s", runtime.GOOS, runtime.GOARCH)

	// Cria o servidor LSP
	server := lsp.NewServer()
	
	// Configura o cliente para enviar notificações
	server.SetClient(func(uri string, diagnostics []lsp.Diagnostic) error {
		// Envia notificação de diagnósticos para o cliente
		notification := jsonrpcMessage{
			JSONRPC: "2.0",
			Method:  "textDocument/publishDiagnostics",
			Params:  marshal(lsp.PublishDiagnosticsParams{
				URI:         uri,
				Diagnostics: diagnostics,
			}),
		}
		
		// Serializa e envia a notificação
		notificationJSON, err := json.Marshal(notification)
		if err != nil {
			log.Printf("Erro ao serializar notificação de diagnósticos: %v\n", err)
			return err
		}
		
		fmt.Printf("Content-Length: %d\r\n\r\n%s", len(notificationJSON), notificationJSON)
		return nil
	})

	// Lê a entrada padrão
	go processMessages(server)

	// Aguarda indefinidamente
	select {}
}

// setupLogging configura o log para um arquivo
func setupLogging() {
	// Determina o diretório de logs
	logDir := os.TempDir()
	if homeDir, err := os.UserHomeDir(); err == nil {
		logDir = filepath.Join(homeDir, ".advpl-lsp")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			logDir = os.TempDir()
		}
	}
	
	logFile := filepath.Join(logDir, "advpl-lsp.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Não foi possível abrir o arquivo de log %s: %v\n", logFile, err)
		return
	}
	
	log.SetOutput(file)
	log.Printf("Log configurado em: %s", logFile)
}

// processMessages processa as mensagens JSON-RPC
func processMessages(server *lsp.Server) {
	for {
		message, err := readMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Fim da entrada padrão, encerrando servidor")
				os.Exit(0)
			}
			log.Printf("Erro ao ler mensagem: %v\n", err)
			continue
		}
		
		// Processa a mensagem
		go handleMessage(server, message)
	}
}

// readMessage lê uma mensagem JSON-RPC da entrada padrão
func readMessage() (jsonrpcMessage, error) {
	var message jsonrpcMessage
	
	// Lê o cabeçalho Content-Length
	contentLength, err := readContentLength()
	if err != nil {
		return message, err
	}
	
	// Lê o corpo da mensagem
	body := make([]byte, contentLength)
	_, err = io.ReadFull(os.Stdin, body)
	if err != nil {
		return message, err
	}
	
	// Parse a mensagem JSON-RPC
	if err := json.Unmarshal(body, &message); err != nil {
		return message, err
	}
	
	log.Printf("Recebida mensagem: método=%s, id=%v\n", message.Method, message.ID)
	return message, nil
}

// readContentLength lê o cabeçalho Content-Length da entrada padrão
func readContentLength() (int, error) {
	var contentLength int
	
	// Lê linhas até encontrar uma linha em branco
	for {
		header, err := readLine()
		if err != nil {
			return 0, err
		}
		
		// Linha em branco indica o fim do cabeçalho
		if header == "" {
			break
		}
		
		// Parse o Content-Length
		if strings.HasPrefix(header, "Content-Length:") {
			_, err := fmt.Sscanf(header, "Content-Length: %d", &contentLength)
			if err != nil {
				return 0, fmt.Errorf("erro ao parsear Content-Length: %v", err)
			}
		}
	}
	
	if contentLength == 0 {
		return 0, fmt.Errorf("Content-Length não encontrado ou inválido")
	}
	
	return contentLength, nil
}

// readLine lê uma linha da entrada padrão
func readLine() (string, error) {
	var buf strings.Builder
	
	for {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			return "", err
		}
		
		// CR
		if b[0] == 13 {
			continue
		}
		
		// LF
		if b[0] == 10 {
			return buf.String(), nil
		}
		
		buf.WriteByte(b[0])
	}
}

// handleMessage processa uma mensagem JSON-RPC
func handleMessage(server *lsp.Server, message jsonrpcMessage) {
	var result interface{}
	var err error
	
	ctx := context.Background()
	
	switch message.Method {
	case MethodInitialize:
		var params lsp.InitializeParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de initialize: %v", err), nil)
			return
		}
		result, err = server.Initialize(ctx, params)
		
	case MethodInitialized:
		// Apenas registra o evento, não requer resposta
		log.Println("Servidor inicializado")
		return
		
	case MethodShutdown:
		// Responde com null
		log.Println("Recebido comando de shutdown")
		sendResponse(message.ID, nil)
		return
		
	case MethodExit:
		// Encerra o servidor
		log.Println("Recebido comando de exit, encerrando servidor")
		os.Exit(0)
		
	case MethodTextDocumentDidOpen:
		var params lsp.DidOpenTextDocumentParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de didOpen: %v", err), nil)
			return
		}
		err = server.DidOpen(ctx, params)
		// Não requer resposta
		if err != nil {
			log.Printf("Erro em didOpen: %v\n", err)
		}
		return
		
	case MethodTextDocumentDidChange:
		var params lsp.DidChangeTextDocumentParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de didChange: %v", err), nil)
			return
		}
		err = server.DidChange(ctx, params)
		// Não requer resposta
		if err != nil {
			log.Printf("Erro em didChange: %v\n", err)
		}
		return
		
	case MethodTextDocumentDocumentSymbol:
		var params lsp.DocumentSymbolParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de documentSymbol: %v", err), nil)
			return
		}
		result, err = server.DocumentSymbol(ctx, params)
		
	case MethodTextDocumentCompletion:
		var params lsp.CompletionParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de completion: %v", err), nil)
			return
		}
		result, err = server.Completion(ctx, params)
		
	case MethodTextDocumentHover:
		var params lsp.HoverParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de hover: %v", err), nil)
			return
		}
		result, err = server.Hover(ctx, params)
		
	case MethodTextDocumentDefinition:
		var params lsp.DefinitionParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de definition: %v", err), nil)
			return
		}
		result, err = server.Definition(ctx, params)
		
	case MethodTextDocumentFormatting:
		var params lsp.DocumentFormattingParams
		if err := json.Unmarshal(message.Params, &params); err != nil {
			sendErrorResponse(message.ID, InvalidParams, fmt.Sprintf("Erro ao parsear parâmetros de formatting: %v", err), nil)
			return
		}
		result, err = server.Formatting(ctx, params)
		
	case MethodWorkspaceDidChangeConfiguration:
		// Atualiza a configuração do servidor
		log.Println("Recebida atualização de configuração")
		return
		
	default:
		log.Printf("Método não suportado: %s\n", message.Method)
		sendErrorResponse(message.ID, MethodNotFound, fmt.Sprintf("Método não suportado: %s", message.Method), nil)
		return
	}
	
	// Verifica se ocorreu algum erro
	if err != nil {
		log.Printf("Erro ao processar método %s: %v\n", message.Method, err)
		sendErrorResponse(message.ID, InternalError, err.Error(), nil)
		return
	}
	
	// Envia a resposta
	sendResponse(message.ID, result)
}

// sendResponse envia uma resposta JSON-RPC
func sendResponse(id interface{}, result interface{}) {
	if id == nil {
		return // Não envia resposta para notificações
	}
	
	response := jsonrpcMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	
	// Serializa e envia a resposta
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Erro ao serializar resposta: %v\n", err)
		return
	}
	
	fmt.Printf("Content-Length: %d\r\n\r\n%s", len(responseJSON), responseJSON)
}

// sendErrorResponse envia uma resposta de erro JSON-RPC
func sendErrorResponse(id interface{}, code int, message string, data interface{}) {
	if id == nil {
		return // Não envia resposta para notificações
	}
	
	response := jsonrpcMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: &struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		}{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	
	// Serializa e envia a resposta
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Erro ao serializar resposta de erro: %v\n", err)
		return
	}
	
	fmt.Printf("Content-Length: %d\r\n\r\n%s", len(responseJSON), responseJSON)
}

// marshal serializa um objeto para JSON
func marshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Erro ao serializar objeto: %v\n", err)
		return nil
	}
	return data
}
