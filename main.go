package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

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
	fmt.Println("TokenVerify invoked")
	return nil
}
