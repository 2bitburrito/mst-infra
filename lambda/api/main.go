package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var db *sql.DB

type Response struct {
	Valid    bool   `json:"valid"`
	Messages string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

func init() {
	dbConnectionString, err := CheckEnv()
	if err != nil {
		fmt.Println("error getting environmental vars:", err.Error())
		panic(err)
	}

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		fmt.Println("error establishing db connection", err.Error())
		panic(err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(3)

	fmt.Println("database Connection established")
}

func setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/user", postUser)
	router.HandleFunc("PATCH /api/user/{id}", patchUser)
	router.HandleFunc("GET /api/user/{id}", getUser)
	router.HandleFunc("DELETE /api/user", deleteUser)

	router.HandleFunc("POST /api/license", postLicense)
	router.HandleFunc("PATCH /api/license/{id}", patchLicense)
	router.HandleFunc("GET /api/license/{id}", getLicense)

	router.HandleFunc("GET /api/license/check/", checkLicense)

	return router
}

func main() {
	router := setupRouter()
	if db == nil {
		log.Fatalf("database connection is not established")
	}
	err := db.Ping()
	if err != nil {
		log.Fatalf("database failed to ping")
	}
	http.ListenAndServe(":3000", router)
}
