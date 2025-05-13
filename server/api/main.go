package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Response struct {
	Valid    bool   `json:"valid"`
	Messages string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

var db *sql.DB

func init() {
	dbUrl, err := CheckEnv()
	if err != nil {
		fmt.Println("error getting environmental vars:", err.Error())
		panic(err)
	}

	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("error establishing db connection", err.Error())
		panic(err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(3)
}

func setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/check-health", checkHealth)

	router.HandleFunc("POST /api/user", postUser)
	router.HandleFunc("PATCH /api/user/{id}", patchUser)
	router.HandleFunc("GET /api/user/{id}", getUser)
	router.HandleFunc("DELETE /api/user", deleteUser)

	router.HandleFunc("POST /api/license, port", postLicense)
	router.HandleFunc("PATCH /api/license/{id}", patchLicense)
	router.HandleFunc("GET /api/license/{id}", getLicense)

	router.HandleFunc("GET /api/license/check/", checkLicense)

	return router
}

func checkHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	router := setupRouter()
	const port = "2000"
	if db == nil {
		log.Fatalf("database connection is not established")
	}
	err := db.Ping()
	if err != nil {
		log.Fatalf("database failed to ping")
	}
	fmt.Println("db ping success")
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	fmt.Println("Running on port:", port)
	server.ListenAndServe()
}
