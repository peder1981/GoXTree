package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type PixTransfer struct {
	PagadorChavePix         string
	PagadorAgencia          string
	PagadorConta            string
	RecebedorCpfCnpj        string
	RecebedorTipoChave      string
	RecebedorChavePix       string
	RecebedorNomeFavorecido string
	Valor                   string
	Descricao               string
	DataCriacao             string
	Status                  string
	ValorTarifa             string
	Motivo                  string
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

func insertData(db *sql.DB, data PixTransfer) {
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
		data.PagadorChavePix,
		data.PagadorAgencia,
		data.PagadorConta,
		data.RecebedorCpfCnpj,
		data.RecebedorTipoChave,
		data.RecebedorChavePix,
		data.RecebedorNomeFavorecido,
		data.Valor,
		data.Descricao,
		data.DataCriacao,
		data.Status,
		data.ValorTarifa,
		data.Motivo,
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

func generateRandomData() PixTransfer {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	pagadorChavePix := string(letters[rand.Intn(len(letters))])
	pagadorAgencia := string(letters[rand.Intn(len(letters))])
	pagadorConta := string(letters[rand.Intn(len(letters))])
	recebedorCpfCnpj := string(letters[rand.Intn(len(letters))])
	recebedorTipoChave := string(letters[rand.Intn(len(letters))])
	recebedorChavePix := string(letters[rand.Intn(len(letters))])
	recebedorNomeFavorecido := string(letters[rand.Intn(len(letters))])
	valor := string(letters[rand.Intn(len(letters))])
	descricao := string(letters[rand.Intn(len(letters))])
	dataCriacao := time.Now().Format("2006-01-02 15:04:05")
	status := string(letters[rand.Intn(len(letters))])
	valorTarifa := string(letters[rand.Intn(len(letters))])
	motivo := string(letters[rand.Intn(len(letters))])
	return PixTransfer{
		PagadorChavePix:         pagadorChavePix,
		PagadorAgencia:          pagadorAgencia,
		PagadorConta:            pagadorConta,
		RecebedorCpfCnpj:        recebedorCpfCnpj,
		RecebedorTipoChave:      recebedorTipoChave,
		RecebedorChavePix:       recebedorChavePix,
		RecebedorNomeFavorecido: recebedorNomeFavorecido,
		Valor:                   valor,
		Descricao:               descricao,
		DataCriacao:             dataCriacao,
		Status:                  status,
		ValorTarifa:             valorTarifa,
		Motivo:                  motivo,
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./pix.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)

	for i := 0; i < 100; i++ {
		data := generateRandomData()
		insertData(db, data)
	}

	fmt.Println("Enter client id:")
	var clientId string
	fmt.Scanln(&clientId)
	fmt.Println("Enter client secret:")
	var clientSecret string
	fmt.Scanln(&clientSecret)
	//fmt.Scanln(&token)
	var token string
	token = base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
	fmt.Println("Enter token:", &token)

	tokenData := Token{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Token:        token,
	}
	insertToken(db, tokenData)

	for {
		fmt.Println("1. Insert data")
		fmt.Println("2. Get data")
		fmt.Println("3. Exit")
		var choice string
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			var pagadorChavePix, pagadorAgencia, pagadorConta, recebedorCpfCnpj, recebedorTipoChave, recebedorChavePix, recebedorNomeFavorecido, valor, descricao, status, valorTarifa, motivo string
			fmt.Println("Enter pagador chave pix:")
			fmt.Scanln(&pagadorChavePix)
			fmt.Println("Enter pagador agencia:")
			fmt.Scanln(&pagadorAgencia)
			fmt.Println("Enter pagador conta:")
			fmt.Scanln(&pagadorConta)
			fmt.Println("Enter recebedor cpf cnpj:")
			fmt.Scanln(&recebedorCpfCnpj)
			fmt.Println("Enter recebedor tipo chave:")
			fmt.Scanln(&recebedorTipoChave)
			fmt.Println("Enter recebedor chave pix:")
			fmt.Scanln(&recebedorChavePix)
			fmt.Println("Enter recebedor nome favorecido:")
			fmt.Scanln(&recebedorNomeFavorecido)
			fmt.Println("Enter valor:")
			fmt.Scanln(&valor)
			fmt.Println("Enter descricao:")
			fmt.Scanln(&descricao)
			fmt.Println("Enter status:")
			fmt.Scanln(&status)
			fmt.Println("Enter valor tarifa:")
			fmt.Scanln(&valorTarifa)
			fmt.Println("Enter motivo:")
			fmt.Scanln(&motivo)
			data := PixTransfer{
				PagadorChavePix:         pagadorChavePix,
				PagadorAgencia:          pagadorAgencia,
				PagadorConta:            pagadorConta,
				RecebedorCpfCnpj:        recebedorCpfCnpj,
				RecebedorTipoChave:      recebedorTipoChave,
				RecebedorChavePix:       recebedorChavePix,
				RecebedorNomeFavorecido: recebedorNomeFavorecido,
				Valor:                   valor,
				Descricao:               descricao,
				DataCriacao:             time.Now().Format("2006-01-02 15:04:05"),
				Status:                  status,
				ValorTarifa:             valorTarifa,
				Motivo:                  motivo,
			}
			insertData(db, data)
		case "2":
			rows, err := db.Query("SELECT * FROM pix_transfers")
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				var id int
				var pagadorChavePix, pagadorAgencia, pagadorConta, recebedorCpfCnpj, recebedorTipoChave, recebedorChavePix, recebedorNomeFavorecido, valor, descricao, dataCriacao, status, valorTarifa, motivo string
				err = rows.Scan(&id, &pagadorChavePix, &pagadorAgencia, &pagadorConta, &recebedorCpfCnpj, &recebedorTipoChave, &recebedorChavePix, &recebedorNomeFavorecido, &valor, &descricao, &dataCriacao, &status, &valorTarifa, &motivo)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(id, pagadorChavePix, pagadorAgencia, pagadorConta, recebedorCpfCnpj, recebedorTipoChave, recebedorChavePix, recebedorNomeFavorecido, valor, descricao, dataCriacao, status, valorTarifa, motivo)
			}
			if err = rows.Err(); err != nil {
				log.Fatal(err)
			}
		case "3":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please select 1, 2, or 3.")
		}
	}
}
