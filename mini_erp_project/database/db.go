package database

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

func InitDB(user string) *sql.DB {
    dbPath := "C:/Users/" + user + "/OneDrive/mini_erp/" + user + "/database.db"
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Fatal(err)
    }
    return db
}
