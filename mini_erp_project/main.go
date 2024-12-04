package main

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", indexHandler)
    r.HandleFunc("/login", loginHandler).Methods("POST")
    r.HandleFunc("/register", registerHandler).Methods("POST")

    log.Println("Servidor rodando na porta 8080")
    http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.html")
}
