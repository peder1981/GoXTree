package modules

import (
    "net/http"
    "github.com/gorilla/sessions"
    "github.com/peder1981/my-go-projects/mini_erp_project/database"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func loginHandler(w http.ResponseWriter, r *http.Request) {
    // Lógica de autenticação
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    // Lógica de registro de usuário
    user := r.FormValue("username")
    db := database.InitDB(user)
    defer db.Close()
    // Criar banco de dados e pastas no OneDrive
}
