package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"crypto/rand"
	"encoding/hex"

	_ "github.com/mattn/go-sqlite3"
)

type PixTransferRequest struct {
	Pagador struct {
		TipoChave string `json:"tipoChave"`
		ChavePix  string `json:"chavePix"`
		Agencia   string `json:"agencia"`
		Conta     string `json:"conta"`
	} `json:"pagador"`
	Recebedor struct {
		CpfCnpj        string `json:"cpfCnpj"`
		TipoChave      string `json:"tipoChave"`
		ChavePix       string `json:"chavePix"`
		NomeFavorecido string `json:"nomeFavorecido"`
	} `json:"recebedor"`
	IdTransacao string `json:"idTransacao"`
	Valor       string `json:"valor"`
	Descricao   string `json:"descricao"`
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
		ChavePix       string `json:"chavePix"`
		NomeFavorecido string `json:"nomeFavorecido"`
	} `json:"recebedor"`
	Valor       string `json:"valor"`
	E2e         string `json:"e2e"`
	IdTransacao string `json:"idTransacao"`
	Descricao   string `json:"descricao"`
	DataCriacao string `json:"dataCriacao"`
	Status      string `json:"status"`
	ValorTarifa string `json:"valorTarifa"`
	Motivo      string `json:"motivo"`
}

type Token struct {
	ClientId     string
	ClientSecret string
	Token        string
}

func createTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pix_transfers (
			id INTEGER PRIMARY KEY,
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
			motivo TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tokens (
			id INTEGER PRIMARY KEY,
			client_id TEXT,
			client_secret TEXT,
			token TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func insertData(db *sql.DB, data PixTransferRequest) {
	_, err := db.Exec(`
		INSERT INTO pix_transfers (
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
			motivo
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`,
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
		"CONCLUIDO",
		"0.00",
		"Transação realizada com sucesso",
	)
	if err != nil {
		log.Fatal(err)
	}
}

func handlePixTransfer(w http.ResponseWriter, r *http.Request) {
	var transferRequest PixTransferRequest
	err := json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "../interface/pix.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	createTable(db)
	insertData(db, transferRequest)
	
	uuid, err := generateUUIDV4()

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
			TipoConta: "CONTA_CORRENTE",
		},
		Recebedor: struct {
			CpfCnpj        string `json:"cpfCnpj"`
			TipoChave      string `json:"tipoChave"`
			ChavePix       string `json:"chavePix"`
			NomeFavorecido string `json:"nomeFavorecido"`
		}{
			CpfCnpj:        transferRequest.Recebedor.CpfCnpj,
			TipoChave:      transferRequest.Recebedor.TipoChave,
			ChavePix:       transferRequest.Recebedor.ChavePix,
			NomeFavorecido: transferRequest.Recebedor.NomeFavorecido,
		},
		Valor:       transferRequest.Valor,
		//E2e:         "E60746948202211301715L2856oOXfn8",
		E2e:         uuid,
		IdTransacao: transferRequest.IdTransacao,
		Descricao:   transferRequest.Descricao,
		DataCriacao: time.Now().Format("2006-01-02T15:04:05.999Z"),
		Status:      "CONCLUIDO",
		ValorTarifa: "0.00",
		Motivo:      "Transação realizada com sucesso",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transferResponse)
}

func handleOAuthToken(w http.ResponseWriter, r *http.Request) {
	clientId := r.Header.Get("client_id")
	clientSecret := r.Header.Get("client_secret")

	if clientId == "" || clientSecret == "" {
		http.Error(w, "client_id and client_secret are required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "../interface/pix.db")
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

func generateUUIDV4() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return hex.EncodeToString(uuid), nil
}

func main() {
	fmt.Println("Iniciando servidor...")
	
	http.HandleFunc("/v1/spi/solicitar-transferencia", handlePixTransfer)
	http.HandleFunc("/oauth/token", handleOAuthToken)

	// Canal para sincronização de erros
	errChan := make(chan error, 2)

	// Iniciar servidor HTTP em uma goroutine
	go func() {
		fmt.Println("Iniciando servidor HTTP na porta 9090...")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			fmt.Printf("Erro ao iniciar servidor HTTP: %v\n", err)
			errChan <- err
		}
	}()

	// Iniciar servidor HTTPS em outra goroutine
	go func() {
		fmt.Println("Iniciando servidor HTTPS na porta 9093...")
		if err := http.ListenAndServeTLS(":9093", "server.crt", "server.key", nil); err != nil {
			fmt.Printf("Erro ao iniciar servidor HTTPS: %v\n", err)
			errChan <- err
		}
	}()

	// Aguardar por erros de qualquer um dos servidores
	err := <-errChan
	log.Fatal(err)
}
