package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

var db *sql.DB

func main() {
	router := mux.NewRouter()
	pgURL, err := pq.ParseURL("postgres://yhctyufj:DvSSi2tltLmRJ3OQKzRuNH15h_j9avy3@isilo.db.elephantsql.com:5432/yhctyufj")

	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", pgURL)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleWare(protectedEndpoint)).Methods("GET")

	log.Println("Listening on Port 3333...")
	log.Fatal(http.ListenAndServe(":3333", router))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Successfully called signup"))
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Successfully called login"))
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup invoked")
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyMiddleware invoked")
	return nil
}
