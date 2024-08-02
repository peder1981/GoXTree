package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Criar ou abrir o banco de dados
	db, err := sql.Open(
		"postgres",
		"host=localhost port=5432 user=golang password=g0l4ng dbname=golang sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Criação das tabelas
	tables := []string{
		`CREATE TABLE IF NOT EXISTS ZC1510 (
            ZC1_FILIAL INTEGER,
            ZC1_CODIGO INTEGER,
            ZC1_DESCRI TEXT,
            ZC1_TIPO TEXT,
            ZC1_ATIVO INTEGER,
            PRIMARY KEY (ZC1_CODIGO)
        );`,
		`CREATE TABLE IF NOT EXISTS ZC2510 (
            ZC2_FILIAL INTEGER,
            ZC2_CODIGO INTEGER,
            ZC2_DESAGR TEXT,
            ZC2_TIPAGR TEXT,
            ZC2_ATIVO INTEGER,
            PRIMARY KEY (ZC2_CODIGO)
        );`,
		`CREATE TABLE IF NOT EXISTS ZC3510 (
            ZC3_FILIAL INTEGER,
            ZC3_DESAGR TEXT,
            ZC3_TIPAGR TEXT,
            ZC3_ATIVO INTEGER,
            PRIMARY KEY (ZC3_TIPAGR)
        );`,
	}

	// Executar as instruções de criação
	for _, table := range tables {
		_, err = db.Exec(table)
		if err != nil {
			log.Fatalf("Erro ao criar tabela: %v", err)
		}
	}

	fmt.Println("Tabelas criadas com sucesso!")
}
