package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	TIPO_CHAVE_TELEFONE     = "TELEFONE"
	TIPO_CHAVE_EMAIL        = "EMAIL"
	TIPO_CHAVE_CPFCNPJ      = "CPFCNPJ"
	TIPO_CHAVE_EVP          = "EVP"
	TIPO_CHAVE_AGENCIACONTA = "AGENCIACONTA"

	TIPO_CONTA_CORRENTE  = "CONTA_CORRENTE"
	TIPO_CONTA_POUPANCA  = "CONTA_POUPANCA"
	TIPO_CONTA_PAGAMENTO = "CONTA_PAGAMENTO"
	TIPO_CONTA_SALARIO   = "CONTA_SALARIO"
)

// Códigos de rejeição
const (
	REJEITADO_GENERICO                     = "2001650"
	REJEITADO_CONTA_DEB_IGUAL_CRED         = "2001651"
	REJEITADO_AGENTE_SAQUE                 = "2001652"
	REJEITADO_RECEBEDOR_CONTATAR_BANCO     = "2001653"
	REJEITADO_CONTA_ORIGEM_RESTRICAO       = "2001654"
	REJEITADO_SEGURANCA                    = "2001655"
	REJEITADO_CONTA_CREDITO_NAO_ENCONTRADA = "2001656"
	REJEITADO_CONTA_CRED_SEM_PERMISSAO     = "2001657"
	REJEITADO_SALDO_INSUFICIENTE           = "2001658"
	REJEITADO_CONTA_DEB_SEM_PERMISSAO      = "2001659"
	REJEITADO_SEM_LIMITE                   = "200165" // 10-15 compartilham prefixo
	REJEITADO_RECEBEDOR_CONTATAR_BANCO2    = "20016516"
	REJEITADO_DEVOLUCAO_NAO_PERMITIDA      = "20016519" // 19,20,24-26,29,31-36
	REJEITADO_CONTA_DEB_IGUAL_CRED2        = "20016521"
	REJEITADO_ISPB_INVALIDO                = "20016522"
	REJEITADO_OPERACAO_NAO_PERMITIDA       = "20016523"
	REJEITADO_SEGURANCA2                   = "20016527" // 27,28,30,37
	REJEITADO_SISTEMA_INDISPONIVEL         = "20016538"
	REJEITADO_OUTRA_INSTITUICAO            = "20016539"
	REJEITADO_CONTA_DEBITO_INVALIDA        = "20016540" // 40-42,44
	REJEITADO_CONTA_CREDITO_INVALIDA       = "20016543" // 43,45
	REJEITADO_GENERICO2                    = "2003333"
	REJEITADO_PSP                          = "2007777"
	REJEITADO_CLIENTE                      = "2009999"
)

// Mapa de códigos de rejeição para mensagens
var motivosRejeicao = map[string]string{
	REJEITADO_GENERICO:                     "Não foi possível efetuar a transação",
	REJEITADO_CONTA_DEB_IGUAL_CRED:         "CONTA DE DEBITO IGUAL CONTA DE CREDITO",
	REJEITADO_AGENTE_SAQUE:                 "AGENTE DE SAQUE NAO PERMITE OPERACAO",
	REJEITADO_RECEBEDOR_CONTATAR_BANCO:     "Não foi possível efetuar a transação. Oriente o recebedor a contatar seu banco",
	REJEITADO_CONTA_ORIGEM_RESTRICAO:       "CLIENTE CONTA ORIGEM POSSUI RESTRICAO",
	REJEITADO_SEGURANCA:                    "MOTOR DE SEGURANCA RECUSOU A TRANSACAO",
	REJEITADO_CONTA_CREDITO_NAO_ENCONTRADA: "CONTA DE CREDITO NAO ENCONTRADA",
	REJEITADO_CONTA_CRED_SEM_PERMISSAO:     "CONTA CREDITO SEM PERMISSAO DE OPERACAO",
	REJEITADO_SALDO_INSUFICIENTE:           "CONTA NAO POSSUI SALDO SUFICIENTE",
	REJEITADO_CONTA_DEB_SEM_PERMISSAO:      "CONTA DEBITO SEM PERMISSAO DE OPERACAO",
	REJEITADO_SEM_LIMITE + "10":            "CONTA SEM LIMITE SUFICIENTE P/ TRANSACAO",
	REJEITADO_SEM_LIMITE + "11":            "CONTA SEM LIMITE SUFICIENTE P/ TRANSACAO",
	REJEITADO_SEM_LIMITE + "12":            "CONTA SEM LIMITE SUFICIENTE P/ TRANSACAO",
	REJEITADO_SEM_LIMITE + "13":            "CONTA SEM LIMITE SUFICIENTE P/ TRANSACAO",
	REJEITADO_SEM_LIMITE + "14":            "CONTA SEM LIMITE SUFICIENTE P/ TRANSACAO",
	REJEITADO_SEM_LIMITE + "15":            "CONTA SEM LIMITE SUFICIENTE P/ TRANSACAO",
	REJEITADO_RECEBEDOR_CONTATAR_BANCO2:    "Não foi possível efetuar a transação. Oriente o recebedor a contatar seu banco",
	REJEITADO_DEVOLUCAO_NAO_PERMITIDA:      "DEVOLUCAO NAO PERMITIDA",
	REJEITADO_CONTA_DEB_IGUAL_CRED2:        "CONTA DE DEBITO IGUAL CONTA DE CREDITO",
	REJEITADO_ISPB_INVALIDO:                "ISPB DO PSP DO RECEBEDOR INVALIDO",
	REJEITADO_OPERACAO_NAO_PERMITIDA:       "OPERACAO NAO PODE SER REALIZADA",
	REJEITADO_SEGURANCA2:                   "MOTOR DE SEGURANCA RECUSOU A TRANSACAO",
	REJEITADO_SISTEMA_INDISPONIVEL:         "SISTEMA TEMPORARIAMENTE INDISPONIVEL",
	REJEITADO_OUTRA_INSTITUICAO:            "PAGAMENTO REJEITADO P/ OUTRA INSTITUICAO",
	REJEITADO_CONTA_DEBITO_INVALIDA:        "DADOS DA CONTA DE DEBITO INVALIDA",
	REJEITADO_CONTA_CREDITO_INVALIDA:       "DADOS DA CONTA DE CREDITO INVALIDA",
	REJEITADO_GENERICO2:                    "Não foi possível efetuar a transação",
	REJEITADO_PSP:                          "Cancelado pelo PSP",
	REJEITADO_CLIENTE:                      "Cancelado pelo cliente",
}

type ErrorResponse struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	Violacoes []struct {
		Razao       string `json:"razao"`
		Propriedade string `json:"propriedade"`
		Valor       string `json:"valor,omitempty"`
	} `json:"violacoes,omitempty"`
}

type PixTransferRequest struct {
	Pagador struct {
		TipoChave string `json:"tipoChave"`
		ChavePix  string `json:"chavePix,omitempty"`
		Agencia   string `json:"agencia,omitempty"`
		Conta     string `json:"conta,omitempty"`
	} `json:"pagador"`
	Recebedor struct {
		CpfCnpj        string `json:"cpfCnpj,omitempty"`
		TipoChave      string `json:"tipoChave"`
		TipoConta      string `json:"tipoConta,omitempty"`
		ChavePix       string `json:"chavePix,omitempty"`
		Ispb           string `json:"ispb,omitempty"`
		Agencia        string `json:"agencia,omitempty"`
		Conta          string `json:"conta,omitempty"`
		DigitoConta    string `json:"digitoConta,omitempty"`
		Banco          string `json:"banco,omitempty"`
		NomeFavorecido string `json:"nomeFavorecido,omitempty"`
	} `json:"recebedor"`
	IdTransacao string `json:"idTransacao"`
	Valor       string `json:"valor"`
	Descricao   string `json:"descricao,omitempty"`
}

type PixTransferResponse struct {
	Pagador struct {
		CpfCnpj   string `json:"cpfCnpj"`
		Agencia   string `json:"agencia"`
		Conta     string `json:"conta"`
		TipoConta string `json:"tipoConta"`
	} `json:"pagador"`
	Recebedor struct {
		CpfCnpj        string `json:"cpfCnpj"`
		TipoChave      string `json:"tipoChave"`
		TipoConta      string `json:"tipoConta,omitempty"`
		ChavePix       string `json:"chavePix,omitempty"`
		Ispb           string `json:"ispb,omitempty"`
		Agencia        string `json:"agencia,omitempty"`
		Conta          string `json:"conta,omitempty"`
		Banco          string `json:"banco,omitempty"`
		NomeFavorecido string `json:"nomeFavorecido"`
	} `json:"recebedor"`
	Valor        string `json:"valor"`
	E2e          string `json:"e2e"`
	IdTransacao  string `json:"idTransacao"`
	Descricao    string `json:"descricao,omitempty"`
	DataCriacao  string `json:"dataCriacao"`
	Status       string `json:"status"`
	ValorTarifa  string `json:"valorTarifa"`
	Motivo       string `json:"motivo"`
	CodigoMotivo string `json:"codigoMotivo,omitempty"`
}

type Token struct {
	ClientId     string
	ClientSecret string
	Token        string
}

func init() {
	// Removing rand.Seed as it's not needed with crypto/rand
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "../interface/pix.db?_timeout=5000&_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pix_transfers (
			id INTEGER PRIMARY KEY,
			id_transacao TEXT,
			e2e TEXT,
			pagador_chave_pix TEXT,
			pagador_agencia TEXT,
			pagador_conta TEXT,
			recebedor_cpf_cnpj TEXT,
			recebedor_tipo_chave TEXT,
			recebedor_chave_pix TEXT,
			recebedor_nome_favorecido TEXT,
			valor TEXT,
			descricao TEXT,
			data_criacao TEXT,
			status TEXT,
			valor_tarifa TEXT,
			motivo TEXT,
			codigo_motivo TEXT
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			id INTEGER PRIMARY KEY,
			client_id TEXT,
			client_secret TEXT,
			token TEXT
		);
	`)
	return err
}

func getRandomStatus() string {
	// Using crypto/rand instead of math/rand
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		// In case of error, default to CONCLUIDO
		return "CONCLUIDO"
	}

	// If the byte is even, return CONCLUIDO, if odd return EM_PROCESSAMENTO
	if b[0]%2 == 0 {
		return "CONCLUIDO"
	}
	return "EM_PROCESSAMENTO"
}

func getMotivo(status string) string {
	if status == "CONCLUIDO" {
		return "Transação realizada com sucesso"
	}
	return "Transação em processamento"
}

func insertData(db *sql.DB, data PixTransferRequest, e2e string) string {
	status := getRandomStatus()
	_, err := db.Exec(`
		INSERT INTO pix_transfers (
			id_transacao,
			e2e,
			pagador_chave_pix,
			pagador_agencia,
			pagador_conta,
			recebedor_cpf_cnpj,
			recebedor_tipo_chave,
			recebedor_chave_pix,
			recebedor_nome_favorecido,
			valor,
			descricao,
			data_criacao,
			status,
			valor_tarifa,
			motivo,
			codigo_motivo
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`,
		data.IdTransacao,
		e2e,
		data.Pagador.ChavePix,
		data.Pagador.Agencia,
		data.Pagador.Conta,
		data.Recebedor.CpfCnpj,
		data.Recebedor.TipoChave,
		data.Recebedor.ChavePix,
		data.Recebedor.NomeFavorecido,
		data.Valor,
		data.Descricao,
		time.Now().Format("2006-01-02T15:04:05.999Z"),
		status,
		"0.00",
		getMotivo(status),
		"00",
	)
	if err != nil {
		log.Fatal(err)
	}
	return status
}

func validateTransferRequest(req PixTransferRequest) (*ErrorResponse, error) {
	// Validar campos obrigatórios
	if req.Valor == "" {
		return &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "A requisição que busca alterar ou criar um(a) transferência não respeita o schema ou está semanticamente errada.",
			Violacoes: []struct {
				Razao       string `json:"razao"`
				Propriedade string `json:"propriedade"`
				Valor       string `json:"valor,omitempty"`
			}{
				{
					Razao:       "O campo valor é obrigatório",
					Propriedade: "valor",
				},
			},
		}, nil
	}

	// Validar formato do valor
	valorRegex := regexp.MustCompile(`^\d*[0-9\.]*(\.)([0-9]{2})$`)
	if !valorRegex.MatchString(req.Valor) {
		return &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "Formato do valor inválido",
			Violacoes: []struct {
				Razao       string `json:"razao"`
				Propriedade string `json:"propriedade"`
				Valor       string `json:"valor,omitempty"`
			}{
				{
					Razao:       "O campo valor deve estar no formato 0.00",
					Propriedade: "valor",
					Valor:       req.Valor,
				},
			},
		}, nil
	}

	// Validar tipo de chave do pagador
	switch req.Pagador.TipoChave {
	case TIPO_CHAVE_TELEFONE, TIPO_CHAVE_EMAIL, TIPO_CHAVE_CPFCNPJ, TIPO_CHAVE_EVP:
		if req.Pagador.ChavePix == "" {
			return &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
				Title:  "Transferência inválida",
				Status: http.StatusBadRequest,
				Detail: "Chave Pix do pagador é obrigatória para o tipo de chave informado",
			}, nil
		}
	case TIPO_CHAVE_AGENCIACONTA:
		if req.Pagador.Agencia == "" || req.Pagador.Conta == "" {
			return &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
				Title:  "Transferência inválida",
				Status: http.StatusBadRequest,
				Detail: "Agência e conta são obrigatórias para transferência por dados bancários",
			}, nil
		}
	default:
		return &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "Tipo de chave do pagador inválido",
		}, nil
	}

	// Validar tipo de chave do recebedor
	switch req.Recebedor.TipoChave {
	case TIPO_CHAVE_TELEFONE, TIPO_CHAVE_EMAIL, TIPO_CHAVE_CPFCNPJ, TIPO_CHAVE_EVP:
		if req.Recebedor.ChavePix == "" {
			return &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
				Title:  "Transferência inválida",
				Status: http.StatusBadRequest,
				Detail: "Chave Pix do recebedor é obrigatória para o tipo de chave informado",
			}, nil
		}
	case TIPO_CHAVE_AGENCIACONTA:
		if req.Recebedor.Agencia == "" || req.Recebedor.Conta == "" || req.Recebedor.Ispb == "" {
			return &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
				Title:  "Transferência inválida",
				Status: http.StatusBadRequest,
				Detail: "Agência, conta e ISPB são obrigatórios para transferência por dados bancários",
			}, nil
		}
		if req.Recebedor.TipoConta == "" {
			return &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
				Title:  "Transferência inválida",
				Status: http.StatusBadRequest,
				Detail: "Tipo de conta é obrigatório para transferência por dados bancários",
			}, nil
		}
		// Validar tipo de conta
		switch req.Recebedor.TipoConta {
		case TIPO_CONTA_CORRENTE, TIPO_CONTA_POUPANCA, TIPO_CONTA_PAGAMENTO, TIPO_CONTA_SALARIO:
		default:
			return &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
				Title:  "Transferência inválida",
				Status: http.StatusBadRequest,
				Detail: "Tipo de conta inválido",
			}, nil
		}
	default:
		return &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "Tipo de chave do recebedor inválido",
		}, nil
	}

	return nil, nil
}

func validatePixKey(key string, keyType string) bool {
	switch keyType {
	case TIPO_CHAVE_TELEFONE:
		// Formato: +55DDNNNNNNNNN
		return regexp.MustCompile(`^\+55[1-9][0-9]{10}$`).MatchString(key)
	case TIPO_CHAVE_EMAIL:
		// Formato básico de email
		return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(key)
	case TIPO_CHAVE_CPFCNPJ:
		// CPF (11 dígitos) ou CNPJ (14 dígitos)
		return regexp.MustCompile(`^\d{11}$|^\d{14}$`).MatchString(key)
	case TIPO_CHAVE_EVP:
		// UUID v4
		return regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(key)
	default:
		return false
	}
}

func getErrorResponse(code string) *ErrorResponse {
	// Obter a mensagem de rejeição usando a nova função
	detail := GetMotivoRejeicao(code)

	// Determinar o status HTTP apropriado
	// Para códigos de rejeição específicos, usamos StatusAccepted (202)
	// Para outros casos, usamos StatusInternalServerError (500)
	status := http.StatusInternalServerError
	if strings.HasPrefix(code, "200") {
		status = http.StatusAccepted
	}

	// Determinar o tipo de erro com base no status
	errorType := "https://pix.bcb.gov.br/api/v2/error/ErroInternoDoServidor"
	title := "Erro interno do servidor"

	if status == http.StatusAccepted {
		errorType = "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida"
		title = "Transferência inválida"
	}

	return &ErrorResponse{
		Type:   errorType,
		Title:  title,
		Status: status,
		Detail: detail,
	}
}

func writeErrorResponse(w http.ResponseWriter, errResp *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.Status)
	json.NewEncoder(w).Encode(errResp)
}

func handlePixTransfer(w http.ResponseWriter, r *http.Request) {
	var transferRequest PixTransferRequest
	err := json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		errResp := &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "A requisição que busca alterar ou criar um(a) transferência não respeita o schema ou está semanticamente errada.",
		}
		writeErrorResponse(w, errResp)
		return
	}

	// Validar a requisição
	if errResp, err := validateTransferRequest(transferRequest); errResp != nil || err != nil {
		if errResp != nil {
			writeErrorResponse(w, errResp)
		} else {
			writeErrorResponse(w, &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/ErroInternoDoServidor",
				Title:  "Erro interno do servidor",
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			})
		}
		return
	}

	// Validar chaves PIX se fornecidas
	if transferRequest.Pagador.TipoChave != TIPO_CHAVE_AGENCIACONTA && !validatePixKey(transferRequest.Pagador.ChavePix, transferRequest.Pagador.TipoChave) {
		writeErrorResponse(w, &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "Formato da chave PIX do pagador inválido",
		})
		return
	}

	if transferRequest.Recebedor.TipoChave != TIPO_CHAVE_AGENCIACONTA && !validatePixKey(transferRequest.Recebedor.ChavePix, transferRequest.Recebedor.TipoChave) {
		writeErrorResponse(w, &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "Formato da chave PIX do recebedor inválido",
		})
		return
	}

	// Gerar E2E ID
	uuid, err := generateUUIDV4X()
	if err != nil {
		writeErrorResponse(w, &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/ErroInternoDoServidor",
			Title:  "Erro interno do servidor",
			Status: http.StatusInternalServerError,
			Detail: "Erro ao gerar identificador da transação",
		})
		return
	}

	// Validar se o ID da transação já foi usado
	if transferRequest.IdTransacao == "" {
		transferRequest.IdTransacao, err = generateUUIDV4X()
		if err != nil {
			writeErrorResponse(w, &ErrorResponse{
				Type:   "https://pix.bcb.gov.br/api/v2/error/ErroInternoDoServidor",
				Title:  "Erro interno do servidor",
				Status: http.StatusInternalServerError,
				Detail: "Erro ao gerar identificador da transação",
			})
			return
		}
	}

	// Conectar ao banco de dados
	db, err := openDB()
	if err != nil {
		writeErrorResponse(w, &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/ErroInternoDoServidor",
			Title:  "Erro interno do servidor",
			Status: http.StatusInternalServerError,
			Detail: "Erro de conexão com o banco de dados: " + err.Error(),
		})
		return
	}
	defer db.Close()

	// Criar tabelas se necessário
	err = createTable(db)
	if err != nil {
		writeErrorResponse(w, &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/ErroInternoDoServidor",
			Title:  "Erro interno do servidor",
			Status: http.StatusInternalServerError,
			Detail: "Erro ao criar tabelas: " + err.Error(),
		})
		return
	}

	// Verificar se o ID da transação já existe
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pix_transfers WHERE id_transacao = ?", transferRequest.IdTransacao).Scan(&count)
	if err != nil {
		writeErrorResponse(w, &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/ErroInternoDoServidor",
			Title:  "Erro interno do servidor",
			Status: http.StatusInternalServerError,
			Detail: "Erro ao verificar ID da transação: " + err.Error(),
		})
		return
	}
	if count > 0 {
		writeErrorResponse(w, &ErrorResponse{
			Type:   "https://pix.bcb.gov.br/api/v2/error/PgtoTransferenciaOperacaoInvalida",
			Title:  "Transferência inválida",
			Status: http.StatusBadRequest,
			Detail: "ID da transação já utilizado",
		})
		return
	}

	// Inserir a transação
	status := insertData(db, transferRequest, uuid)

	// Preparar resposta
	transferResponse := PixTransferResponse{
		Pagador: struct {
			CpfCnpj   string `json:"cpfCnpj"`
			Agencia   string `json:"agencia"`
			Conta     string `json:"conta"`
			TipoConta string `json:"tipoConta"`
		}{
			CpfCnpj:   transferRequest.Pagador.ChavePix,
			Agencia:   transferRequest.Pagador.Agencia,
			Conta:     transferRequest.Pagador.Conta,
			TipoConta: TIPO_CONTA_CORRENTE,
		},
		Recebedor: struct {
			CpfCnpj        string `json:"cpfCnpj"`
			TipoChave      string `json:"tipoChave"`
			TipoConta      string `json:"tipoConta,omitempty"`
			ChavePix       string `json:"chavePix,omitempty"`
			Ispb           string `json:"ispb,omitempty"`
			Agencia        string `json:"agencia,omitempty"`
			Conta          string `json:"conta,omitempty"`
			Banco          string `json:"banco,omitempty"`
			NomeFavorecido string `json:"nomeFavorecido"`
		}{
			CpfCnpj:        transferRequest.Recebedor.CpfCnpj,
			TipoChave:      transferRequest.Recebedor.TipoChave,
			TipoConta:      transferRequest.Recebedor.TipoConta,
			ChavePix:       transferRequest.Recebedor.ChavePix,
			Ispb:           transferRequest.Recebedor.Ispb,
			Agencia:        transferRequest.Recebedor.Agencia,
			Conta:          transferRequest.Recebedor.Conta,
			Banco:          transferRequest.Recebedor.Banco,
			NomeFavorecido: transferRequest.Recebedor.NomeFavorecido,
		},
		Valor:       transferRequest.Valor,
		E2e:         uuid,
		IdTransacao: transferRequest.IdTransacao,
		Descricao:   transferRequest.Descricao,
		DataCriacao: time.Now().Format("2006-01-02T15:04:05.999Z"),
		Status:      status,
		ValorTarifa: "0.00",
		Motivo:      getMotivo(status),
	}

	// Definir código HTTP com base no status
	httpStatus := http.StatusOK
	if status == "EM_PROCESSAMENTO" {
		httpStatus = http.StatusAccepted
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(transferResponse)
}

func handleGetPixTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) != 7 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	idTransacao := parts[5]
	e2e := parts[6]

	if idTransacao == "" && e2e == "" {
		http.Error(w, "idTransacao or e2e parameter is required", http.StatusBadRequest)
		return
	}

	db, err := openDB()
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Transaction start error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Will rollback if not committed

	var transfer PixTransferResponse
	var rows *sql.Rows

	if idTransacao != "" {
		rows, err = tx.Query(`
            SELECT 
                pagador_chave_pix, pagador_agencia, pagador_conta,
                recebedor_cpf_cnpj, recebedor_tipo_chave, recebedor_chave_pix, recebedor_nome_favorecido,
                valor, descricao, data_criacao, status, valor_tarifa, motivo, id_transacao, e2e, codigo_motivo
            FROM pix_transfers 
            WHERE id_transacao = ?
            ORDER BY id DESC LIMIT 1`, idTransacao)
	} else {
		rows, err = tx.Query(`
            SELECT 
                pagador_chave_pix, pagador_agencia, pagador_conta,
                recebedor_cpf_cnpj, recebedor_tipo_chave, recebedor_chave_pix, recebedor_nome_favorecido,
                valor, descricao, data_criacao, status, valor_tarifa, motivo, id_transacao, e2e, codigo_motivo
            FROM pix_transfers 
            WHERE e2e = ?
            ORDER BY id DESC LIMIT 1`, e2e)
	}

	if err != nil {
		http.Error(w, "Query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	var (
		pagadorChavePix, pagadorAgencia, pagadorConta                                    string
		recebedorCpfCnpj, recebedorTipoChave, recebedorChavePix, recebedorNomeFavorecido string
		valor, descricao, dataCriacao, status, valorTarifa, motivo, codigoMotivo         string
		dbIdTransacao, dbE2e                                                             string
	)

	err = rows.Scan(
		&pagadorChavePix, &pagadorAgencia, &pagadorConta,
		&recebedorCpfCnpj, &recebedorTipoChave, &recebedorChavePix, &recebedorNomeFavorecido,
		&valor, &descricao, &dataCriacao, &status, &valorTarifa, &motivo,
		&dbIdTransacao, &dbE2e, &codigoMotivo,
	)
	if err != nil {
		http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if status == "EM_PROCESSAMENTO" {
		status = "CONCLUIDO"
		motivo = "Transação realizada com sucesso"

		_, err = tx.Exec(`
            UPDATE pix_transfers 
            SET status = ?, motivo = ? 
            WHERE id = (
                SELECT id FROM pix_transfers 
                WHERE id_transacao = ? OR e2e = ?
                ORDER BY id DESC LIMIT 1
            )`,
			status, motivo, dbIdTransacao, dbE2e)

		if err != nil {
			http.Error(w, "Update error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Commit error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	transfer = PixTransferResponse{
		Pagador: struct {
			CpfCnpj   string `json:"cpfCnpj"`
			Agencia   string `json:"agencia"`
			Conta     string `json:"conta"`
			TipoConta string `json:"tipoConta"`
		}{
			CpfCnpj:   pagadorChavePix,
			Agencia:   pagadorAgencia,
			Conta:     pagadorConta,
			TipoConta: "CONTA_CORRENTE",
		},
		Recebedor: struct {
			CpfCnpj        string `json:"cpfCnpj"`
			TipoChave      string `json:"tipoChave"`
			TipoConta      string `json:"tipoConta,omitempty"`
			ChavePix       string `json:"chavePix,omitempty"`
			Ispb           string `json:"ispb,omitempty"`
			Agencia        string `json:"agencia,omitempty"`
			Conta          string `json:"conta,omitempty"`
			Banco          string `json:"banco,omitempty"`
			NomeFavorecido string `json:"nomeFavorecido"`
		}{
			CpfCnpj:        recebedorCpfCnpj,
			TipoChave:      recebedorTipoChave,
			TipoConta:      "",
			ChavePix:       recebedorChavePix,
			Ispb:           "",
			Agencia:        "",
			Conta:          "",
			Banco:          "",
			NomeFavorecido: recebedorNomeFavorecido,
		},
		Valor:        valor,
		E2e:          dbE2e,
		IdTransacao:  dbIdTransacao,
		Descricao:    descricao,
		DataCriacao:  dataCriacao,
		Status:       status,
		ValorTarifa:  valorTarifa,
		Motivo:       motivo,
		CodigoMotivo: codigoMotivo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transfer)
}

func handleOAuthToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
		http.Error(w, "Basic authentication required", http.StatusUnauthorized)
		return
	}

	// Extract the base64 encoded credentials
	encodedCreds := strings.TrimPrefix(authHeader, "Basic ")
	decodedCreds, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		http.Error(w, "Invalid authentication credentials", http.StatusUnauthorized)
		return
	}

	// Split decoded string into client_id and client_secret
	credentials := strings.Split(string(decodedCreds), ":")
	if len(credentials) != 2 {
		http.Error(w, "Invalid authentication format", http.StatusUnauthorized)
		return
	}

	clientId := credentials[0]
	clientSecret := credentials[1]

	if clientId == "" || clientSecret == "" {
		http.Error(w, "client_id and client_secret are required", http.StatusUnauthorized)
		return
	}

	db, err := openDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var token string
	err = db.QueryRow("SELECT token FROM tokens WHERE client_id = ? AND client_secret = ?", clientId, clientSecret).Scan(&token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenType := "Bearer"
	expiresIn := 3600

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token": token,
		"token_type":   tokenType,
		"expires_in":   expiresIn,
	})
}

func authenticateBearer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Bearer authentication required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Validate token from database
		db, err := openDB()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM tokens WHERE token = ?)", token).Scan(&exists)
		if err != nil || !exists {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func generateUUIDV4X() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant (2) bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	// Convert to hex string without hyphens
	return strings.ToUpper(hex.EncodeToString(uuid)), nil
}

// GetMotivoRejeicao retorna a mensagem correspondente ao código de rejeição
func GetMotivoRejeicao(codigo string) string {
	// Verifica se é um dos códigos que compartilham prefixo (10-15)
	if strings.HasPrefix(codigo, REJEITADO_SEM_LIMITE) {
		suffix := codigo[len(REJEITADO_SEM_LIMITE):]
		if suffix >= "10" && suffix <= "15" {
			return motivosRejeicao[REJEITADO_SEM_LIMITE+"10"]
		}
	}

	// Verifica códigos que compartilham a mesma mensagem
	switch codigo {
	case "20016519", "20016520", "20016524", "20016525", "20016526", "20016529",
		"20016531", "20016532", "20016533", "20016534", "20016535", "20016536":
		return motivosRejeicao[REJEITADO_DEVOLUCAO_NAO_PERMITIDA]
	case "20016527", "20016528", "20016530", "20016537":
		return motivosRejeicao[REJEITADO_SEGURANCA2]
	case "20016540", "20016541", "20016542", "20016544":
		return motivosRejeicao[REJEITADO_CONTA_DEBITO_INVALIDA]
	case "20016543", "20016545":
		return motivosRejeicao[REJEITADO_CONTA_CREDITO_INVALIDA]
	}

	if msg, ok := motivosRejeicao[codigo]; ok {
		return msg
	}
	return "Não foi possível efetuar a transação"
}

func main() {
	fmt.Println("Iniciando servidor...")

	// OAuth endpoint uses Basic auth (handled internally in handleOAuthToken)
	http.HandleFunc("/oauth/token", handleOAuthToken)

	// Protected endpoints use Bearer auth
	http.HandleFunc("/v1/spi/solicitar-transferencia", authenticateBearer(handlePixTransfer))
	http.HandleFunc("/v1/spi/consulta/transferencia/", authenticateBearer(handleGetPixTransfer))

	errChan := make(chan error, 2)

	go func() {
		fmt.Println("Iniciando servidor HTTP na porta 9090...")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			fmt.Printf("Erro ao iniciar servidor HTTP: %v\n", err)
			errChan <- err
		}
	}()

	go func() {
		fmt.Println("Iniciando servidor HTTPS na porta 9093...")
		if err := http.ListenAndServeTLS(":9093", "server.crt", "server.key", nil); err != nil {
			fmt.Printf("Erro ao iniciar servidor HTTPS: %v\n", err)
			errChan <- err
		}
	}()

	err := <-errChan
	log.Fatal(err)
}
