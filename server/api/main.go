package main

import (
	"api/config"
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
	cfg, _ := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DB.URL)
	if err != nil {
		fmt.Println("error establishing db connection", err.Error())
		panic(err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(3)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DB Ping Successful")
	}
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

	router.HandleFunc("POST /api/create-login-code", createLoginCode)

	return router
}

func checkHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func middlewareSetup(next http.Handler) http.Handler {
	cfg, _ := config.LoadConfig()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS config
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check API Key:
		apiKey := r.Header.Get("X-API-Key")
		log.Println("API KEY MATCH")
		if apiKey != cfg.ApiKey {
			http.Error(w, "Unauthorized X-API-Key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg, _ := config.LoadConfig()
	router := setupRouter()
	handler := middlewareSetup(router)

	server := &http.Server{
		Handler: handler,
		Addr:    ":" + cfg.Port,
	}
	fmt.Println("Running on port:", cfg.Port)
	server.ListenAndServe()
}
