package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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
func respondWithError(w http.ResponseWriter, status int, error Error) {
	// respond with an error
	// Send status bad request

	w.WriteHeader(status) // 400
	json.NewEncoder(w).Encode(error)
	return
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	var error Error
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is missing."
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is missing."
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}

	user.Password = string(hash)

	stmt := "insert into users (email, password) values($1, $2) RETURNING id;"

	err = db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)

	if err != nil {
		error.Message = "Server error."
		respondWithError(w, http.StatusInternalServerError, error)
		return
	}

	user.Password = ""

	w.Header().Set("Content-Type", "application/json")

	responseJSON(w, user)

	// spew.Dump(user)
}

func GenerateToken(user User) (string, error) {
	var err error
	secret := "secret"

	// a JWT
	// 3 parts: header, payload, secret

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	spew.Dump(token)

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

func login(w http.ResponseWriter, r *http.Request) {

	var user User
	// var jwt JWT
	var error Error

	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	// password := user.Password

	row := db.QueryRow("select * from users where email=$1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			error.Message = "The user does not exist!"
			respondWithError(w, http.StatusBadRequest, error)
			return
		} else {
			log.Fatal(err)
		}
	}

	spew.Dump(user)

	// token, err := GenerateToken(user)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(token)
	// w.Write([]byte("Successfully called login"))
}

func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup invoked")
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyMiddleware invoked")
	return nil
}
