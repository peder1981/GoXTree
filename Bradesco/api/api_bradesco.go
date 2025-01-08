package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

func init() {
	// Removing rand.Seed as it's not needed with crypto/rand
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../interface/pix.db?_timeout=5000&_journal_mode=WAL&_busy_timeout=5000")
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
			motivo TEXT
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
			motivo
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
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
	)
	if err != nil {
		log.Fatal(err)
	}
	return status
}

func handlePixTransfer(w http.ResponseWriter, r *http.Request) {
	var transferRequest PixTransferRequest
	err := json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := generateUUIDV4X()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if transferRequest.IdTransacao == "" {
		transferRequest.IdTransacao, err = generateUUIDV4X()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	db, err := openDB()
	if err != nil {
		http.Error(w, "Database connection error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		http.Error(w, "Table creation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	status := insertData(db, transferRequest, uuid)

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
		E2e:         uuid,
		IdTransacao: transferRequest.IdTransacao,
		Descricao:   transferRequest.Descricao,
		DataCriacao: time.Now().Format("2006-01-02T15:04:05.999Z"),
		Status:      status,
		ValorTarifa: "0.00",
		Motivo:      getMotivo(status),
	}

	w.Header().Set("Content-Type", "application/json")
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
                valor, descricao, data_criacao, status, valor_tarifa, motivo, id_transacao, e2e
            FROM pix_transfers 
            WHERE id_transacao = ?
            ORDER BY id DESC LIMIT 1`, idTransacao)
	} else {
		rows, err = tx.Query(`
            SELECT 
                pagador_chave_pix, pagador_agencia, pagador_conta,
                recebedor_cpf_cnpj, recebedor_tipo_chave, recebedor_chave_pix, recebedor_nome_favorecido,
                valor, descricao, data_criacao, status, valor_tarifa, motivo, id_transacao, e2e
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
		valor, descricao, dataCriacao, status, valorTarifa, motivo                       string
		dbIdTransacao, dbE2e                                                             string
	)

	err = rows.Scan(
		&pagadorChavePix, &pagadorAgencia, &pagadorConta,
		&recebedorCpfCnpj, &recebedorTipoChave, &recebedorChavePix, &recebedorNomeFavorecido,
		&valor, &descricao, &dataCriacao, &status, &valorTarifa, &motivo,
		&dbIdTransacao, &dbE2e,
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
			ChavePix       string `json:"chavePix"`
			NomeFavorecido string `json:"nomeFavorecido"`
		}{
			CpfCnpj:        recebedorCpfCnpj,
			TipoChave:      recebedorTipoChave,
			ChavePix:       recebedorChavePix,
			NomeFavorecido: recebedorNomeFavorecido,
		},
		Valor:       valor,
		E2e:         dbE2e,
		IdTransacao: dbIdTransacao,
		Descricao:   descricao,
		DataCriacao: dataCriacao,
		Status:      status,
		ValorTarifa: valorTarifa,
		Motivo:      motivo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transfer)
}

func handleOAuthToken(w http.ResponseWriter, r *http.Request) {
	clientId := r.Header.Get("client_id")
	clientSecret := r.Header.Get("client_secret")

	if clientId == "" || clientSecret == "" {
		http.Error(w, "client_id and client_secret are required", http.StatusBadRequest)
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

func main() {
	fmt.Println("Iniciando servidor...")

	http.HandleFunc("/v1/spi/solicitar-transferencia", handlePixTransfer)
	http.HandleFunc("/v1/spi/consulta/transferencia/", handleGetPixTransfer)
	http.HandleFunc("/oauth/token", handleOAuthToken)

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
