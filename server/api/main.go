package main

import (
	"api/config"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Response struct {
	Valid    bool   `json:"valid"`
	Messages string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

var db *sql.DB

var API_KEY = os.Getenv("API_KEY")

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

	router.HandleFunc("POST /api/create-login-code", checkKeysMiddleware(createLoginCode))

	return router
}

func checkHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func checkKeysMiddleware(next http.HandlerFunc) http.HandlerFunc {
	cfg, _ := config.LoadConfig()
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		fmt.Println("Our api key:", cfg.ApiKey)
		fmt.Println("Received api key:", apiKey)
		if apiKey != cfg.ApiKey {
			http.Error(w, "Unauthorized X-API-Key", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func main() {
	router := setupRouter()
	cfg, _ := config.LoadConfig()
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
		Addr:    ":" + cfg.Port,
	}
	fmt.Println("Running on port:", cfg.Port)
	server.ListenAndServe()
}
