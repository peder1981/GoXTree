package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/peder1981/mini_erp/cloud"
	"github.com/peder1981/mini_erp/database"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Aqui você deve implementar a lógica de verificação de usuário e senha
	// Para fins de exemplo, vamos apenas logar o usuário
	log.Printf("Usuário %s logado com sucesso!", username)

	// Criar sessão
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Aqui você deve implementar a lógica de registro de usuário
	log.Printf("Registrando usuário: %s", username)

	// Inicializar o banco de dados
	db := database.InitDB(username)
	defer db.Close()

	// Criar pasta no OneDrive
	cloud.CreateFolder(username)

	// Redirecionar após o registro
	http.Redirect(w, r, "/", http.StatusFound)
}
