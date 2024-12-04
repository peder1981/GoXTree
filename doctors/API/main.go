package main

import (
    "log"
    "net/http"

    "doctors/db"
    "doctors/handlers"

    "github.com/gorilla/mux"
)

func main() {
    db.Init()

    r := mux.NewRouter()

    r.HandleFunc("/doctors", handlers.CreateDoctor).Methods("POST")
    r.HandleFunc("/doctors/{id}", handlers.GetDoctor).Methods("GET")
    r.HandleFunc("/doctors/{id}", handlers.UpdateDoctor).Methods("PUT")
    r.HandleFunc("/doctors/{id}", handlers.DeleteDoctor).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8081", r))
}
