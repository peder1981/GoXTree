package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var tokens = map[string]string{}
var accounts = map[string]float64{
	"12345": 1000.0,
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PaymentRequest struct {
	ReceiverAccount string  `json:"receiver_account"`
	Amount          float64 `json:"amount"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type Transaction struct {
	SenderAccount   string  `json:"sender_account"`
	ReceiverAccount string  `json:"receiver_account"`
	Amount          float64 `json:"amount"`
}

type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := base64.StdEncoding.EncodeToString([]byte(loginRequest.Username + ":" + loginRequest.Password))
	tokens[token] = loginRequest.Username

	var loginResponse LoginResponse
	loginResponse.Token = token

	json.NewEncoder(w).Encode(loginResponse)
}

func authenticate(token string) (string, error) {
	username, ok := tokens[token]
	if !ok {
		return "", errors.New("invalid token")
	}
	return username, nil
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "token is required", http.StatusUnauthorized)
		return
	}

	token = strings.Replace(token, "Basic ", "", 1)
	username, err := authenticate(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	balance, ok := accounts[username]
	if !ok {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	var balanceResponse BalanceResponse
	balanceResponse.Balance = balance

	json.NewEncoder(w).Encode(balanceResponse)
}

func makePayment(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "token is required", http.StatusUnauthorized)
		return
	}

	token = strings.Replace(token, "Basic ", "", 1)
	username, err := authenticate(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var paymentRequest PaymentRequest
	err = json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, ok := accounts[username]
	if !ok {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	if balance < paymentRequest.Amount {
		http.Error(w, "insufficient balance", http.StatusConflict)
		return
	}

	accounts[username] -= paymentRequest.Amount
	accounts[paymentRequest.ReceiverAccount] += paymentRequest.Amount

	w.WriteHeader(http.StatusCreated)
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "token is required", http.StatusUnauthorized)
		return
	}

	token = strings.Replace(token, "Basic ", "", 1)
	username, err := authenticate(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var transactionsResponse TransactionsResponse
	// Add transactions to the response
	transactionsResponse.Transactions = []Transaction{
		{
			SenderAccount:   username,
			ReceiverAccount: "12345",
			Amount:          100.0,
		},
	}

	json.NewEncoder(w).Encode(transactionsResponse)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/balance", getBalance).Methods("GET")
	router.HandleFunc("/payment", makePayment).Methods("POST")
	router.HandleFunc("/transactions", getTransactions).Methods("GET")

	fmt.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))

	// Clean up tokens every 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			tokens = map[string]string{}
		}
	}()
}
