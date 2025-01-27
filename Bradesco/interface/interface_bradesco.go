package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
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

type PixTransfer struct {
	Pagador struct {
		CpfCnpj   string
		TipoChave string
		ChavePix  string
		Agencia   string
		Conta     string
		TipoConta string
	}
	Recebedor struct {
		CpfCnpj        string
		TipoChave      string
		TipoConta      string
		ChavePix       string
		Ispb           string
		Agencia        string
		Conta          string
		Banco          string
		NomeFavorecido string
	}
	IdTransacao  string
	Valor        string
	E2e          string
	Descricao    string
	DataCriacao  string
	Status       string
	ValorTarifa  string
	Motivo       string
	CodigoMotivo string
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
			pagador_cpf_cnpj TEXT,
			pagador_tipo_chave TEXT,
			pagador_chave_pix TEXT,
			pagador_agencia TEXT,
			pagador_conta TEXT,
			pagador_tipo_conta TEXT,
			recebedor_cpf_cnpj TEXT,
			recebedor_tipo_chave TEXT,
			recebedor_tipo_conta TEXT,
			recebedor_chave_pix TEXT,
			recebedor_ispb TEXT,
			recebedor_agencia TEXT,
			recebedor_conta TEXT,
			recebedor_banco TEXT,
			recebedor_nome_favorecido TEXT,
			id_transacao TEXT,
			valor TEXT,
			e2e TEXT,
			descricao TEXT,
			data_criacao TEXT,
			status TEXT,
			valor_tarifa TEXT,
			motivo TEXT,
			codigo_motivo TEXT
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

func insertData(db *sql.DB, data PixTransfer) {
	_, err := db.Exec(`
		INSERT INTO pix_transfers (
			pagador_cpf_cnpj,
			pagador_tipo_chave,
			pagador_chave_pix,
			pagador_agencia,
			pagador_conta,
			pagador_tipo_conta,
			recebedor_cpf_cnpj,
			recebedor_tipo_chave,
			recebedor_tipo_conta,
			recebedor_chave_pix,
			recebedor_ispb,
			recebedor_agencia,
			recebedor_conta,
			recebedor_banco,
			recebedor_nome_favorecido,
			id_transacao,
			valor,
			e2e,
			descricao,
			data_criacao,
			status,
			valor_tarifa,
			motivo,
			codigo_motivo
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`,
		data.Pagador.CpfCnpj,
		data.Pagador.TipoChave,
		data.Pagador.ChavePix,
		data.Pagador.Agencia,
		data.Pagador.Conta,
		data.Pagador.TipoConta,
		data.Recebedor.CpfCnpj,
		data.Recebedor.TipoChave,
		data.Recebedor.TipoConta,
		data.Recebedor.ChavePix,
		data.Recebedor.Ispb,
		data.Recebedor.Agencia,
		data.Recebedor.Conta,
		data.Recebedor.Banco,
		data.Recebedor.NomeFavorecido,
		data.IdTransacao,
		data.Valor,
		data.E2e,
		data.Descricao,
		data.DataCriacao,
		data.Status,
		data.ValorTarifa,
		data.Motivo,
		data.CodigoMotivo,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func insertToken(db *sql.DB, token Token) {
	_, err := db.Exec(`
		INSERT INTO tokens (
			client_id,
			client_secret,
			token
		) VALUES (?, ?, ?);
	`,
		token.ClientId,
		token.ClientSecret,
		token.Token,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func generateSampleData() PixTransfer {
	var data PixTransfer

	// Pagador
	data.Pagador.CpfCnpj = "18241818000161"
	data.Pagador.TipoChave = TIPO_CHAVE_CPFCNPJ
	data.Pagador.ChavePix = "18241818000161"
	data.Pagador.Agencia = "2856"
	data.Pagador.Conta = "565"
	data.Pagador.TipoConta = TIPO_CONTA_CORRENTE

	// Recebedor
	data.Recebedor.CpfCnpj = "09999902291969"
	data.Recebedor.TipoChave = TIPO_CHAVE_CPFCNPJ
	data.Recebedor.ChavePix = "09999902291969"
	data.Recebedor.NomeFavorecido = "KLEBER ADILSON"

	// Transação
	data.IdTransacao = fmt.Sprintf("TransfenciaAPI%d", time.Now().UnixNano())
	data.Valor = "50.00"
	data.E2e = fmt.Sprintf("E60746948%s", time.Now().Format("20060102150405"))
	data.Descricao = "Pagamento teste"
	data.DataCriacao = time.Now().Format("2006-01-02T15:04:05.000Z")
	data.Status = "CONCLUIDO"
	data.ValorTarifa = "0.00"

	return data
}

func main() {
	db, err := sql.Open("sqlite3", "./pix.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

	fmt.Println("Enter client id:")
	var clientId string
	fmt.Scanln(&clientId)
	fmt.Println("Enter client secret:")
	var clientSecret string
	fmt.Scanln(&clientSecret)

	token := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	fmt.Printf("Generated token: %s\n", token)

	tokenData := Token{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Token:        token,
	}
	insertToken(db, tokenData)

	for {
		fmt.Println("\n=== Menu ===")
		fmt.Println("1. Insert new PIX transfer")
		fmt.Println("2. View all PIX transfers")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option: ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			var data PixTransfer

			// Pagador
			fmt.Print("Enter pagador CPF/CNPJ: ")
			fmt.Scanln(&data.Pagador.CpfCnpj)
			fmt.Print("Enter pagador tipo chave (TELEFONE/EMAIL/CPFCNPJ/EVP/AGENCIACONTA): ")
			fmt.Scanln(&data.Pagador.TipoChave)
			fmt.Print("Enter pagador chave PIX: ")
			fmt.Scanln(&data.Pagador.ChavePix)
			fmt.Print("Enter pagador agência: ")
			fmt.Scanln(&data.Pagador.Agencia)
			fmt.Print("Enter pagador conta: ")
			fmt.Scanln(&data.Pagador.Conta)
			fmt.Print("Enter pagador tipo conta (CONTA_CORRENTE/CONTA_POUPANCA/CONTA_PAGAMENTO/CONTA_SALARIO): ")
			fmt.Scanln(&data.Pagador.TipoConta)

			// Recebedor
			fmt.Print("Enter recebedor CPF/CNPJ: ")
			fmt.Scanln(&data.Recebedor.CpfCnpj)
			fmt.Print("Enter recebedor tipo chave (TELEFONE/EMAIL/CPFCNPJ/EVP/AGENCIACONTA): ")
			fmt.Scanln(&data.Recebedor.TipoChave)
			fmt.Print("Enter recebedor chave PIX: ")
			fmt.Scanln(&data.Recebedor.ChavePix)
			fmt.Print("Enter recebedor nome favorecido: ")
			fmt.Scanln(&data.Recebedor.NomeFavorecido)

			// Transação
			fmt.Print("Enter valor: ")
			fmt.Scanln(&data.Valor)
			fmt.Print("Enter descrição: ")
			fmt.Scanln(&data.Descricao)

			// Campos automáticos
			data.IdTransacao = fmt.Sprintf("TransfenciaAPI%d", time.Now().UnixNano())
			data.E2e = fmt.Sprintf("E60746948%s", time.Now().Format("20060102150405"))
			data.DataCriacao = time.Now().Format("2006-01-02T15:04:05.000Z")
			data.Status = "CONCLUIDO"
			data.ValorTarifa = "0.00"

			insertData(db, data)
			fmt.Println("PIX transfer inserted successfully!")

		case "2":
			rows, err := db.Query("SELECT * FROM pix_transfers")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			fmt.Println("\n=== PIX Transfers ===")
			for rows.Next() {
				var id int
				var data PixTransfer
				err = rows.Scan(&id,
					&data.Pagador.CpfCnpj, &data.Pagador.TipoChave, &data.Pagador.ChavePix,
					&data.Pagador.Agencia, &data.Pagador.Conta, &data.Pagador.TipoConta,
					&data.Recebedor.CpfCnpj, &data.Recebedor.TipoChave, &data.Recebedor.TipoConta,
					&data.Recebedor.ChavePix, &data.Recebedor.Ispb, &data.Recebedor.Agencia,
					&data.Recebedor.Conta, &data.Recebedor.Banco, &data.Recebedor.NomeFavorecido,
					&data.IdTransacao, &data.Valor, &data.E2e, &data.Descricao, &data.DataCriacao,
					&data.Status, &data.ValorTarifa, &data.Motivo, &data.CodigoMotivo)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("\nTransfer ID: %d\n", id)
				fmt.Printf("Pagador: %s (Chave: %s)\n", data.Pagador.CpfCnpj, data.Pagador.ChavePix)
				fmt.Printf("Recebedor: %s (%s)\n", data.Recebedor.NomeFavorecido, data.Recebedor.CpfCnpj)
				fmt.Printf("Valor: R$ %s\n", data.Valor)
				fmt.Printf("Status: %s\n", data.Status)
				if data.Motivo != "" {
					fmt.Printf("Motivo: %s (Código: %s)\n", data.Motivo, data.CodigoMotivo)
				}
				fmt.Printf("Data: %s\n", data.DataCriacao)
				fmt.Println("------------------------")
			}

		case "3":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice. Please select 1, 2, or 3.")
		}
	}
}
