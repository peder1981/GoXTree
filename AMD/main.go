package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

var db *sql.DB

type Despesa struct {
	ID    int     `json:"id"`
	Tipo  string  `json:"tipo"`
	Valor float64 `json:"valor"`
	Mes   string  `json:"mes"`
}

func main() {
	var err error
	// Configure seu banco de dados PostgreSQL
	connStr := "user=yourusername dbname=yourdbname sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.POST("/despesa", cadastrarDespesa)
	router.GET("/relatorio", gerarRelatorio)
	router.Run("localhost:8080")
}

func cadastrarDespesa(c *gin.Context) {
	var despesa Despesa
	if err := c.ShouldBindJSON(&despesa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO despesas (tipo, valor, mes) VALUES ($1, $2, $3)", despesa.Tipo, despesa.Valor, despesa.Mes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "despesa cadastrada com sucesso!"})
}

func gerarRelatorio(c *gin.Context) {
	rows, err := db.Query("SELECT tipo, mes, SUM(valor) FROM despesas GROUP BY tipo, mes ORDER BY tipo, mes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	f := excelize.NewFile()
	index := f.NewSheet("Relatório")

	// Adicionando cabeçalho
	f.SetCellValue("Relatório", "A1", "Tipo")
	f.SetCellValue("Relatório", "B1", "Mês")
	f.SetCellValue("Relatório", "C1", "Total")

	row := 2
	for rows.Next() {
		var despesa Despesa
		if err := rows.Scan(&despesa.Tipo, &despesa.Mes, &despesa.Valor); err != nil {
			log.Fatal(err)
		}
		f.SetCellValue("Relatório", fmt.Sprintf("A%d", row), despesa.Tipo)
		f.SetCellValue("Relatório", fmt.Sprintf("B%d", row), despesa.Mes)
		f.SetCellValue("Relatório", fmt.Sprintf("C%d", row), despesa.Valor)
		row++
	}

	f.SetActiveSheet(index)

	// Salva o arquivo em formato .xlsx
	if err := f.SaveAs("Relatorio_despesas.xlsx"); err != nil {
		log.Fatal(err)
	}

	c.File("Relatorio_despesas.xlsx")
}
