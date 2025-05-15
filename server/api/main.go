package main

import (
	"api/config"
	"api/store"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type API struct {
	db                *sql.DB
	verificationStore *store.VerificationStore
}

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

func (api *API) setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/check-health", api.checkHealth)

	router.HandleFunc("POST /api/user", api.postUser)
	router.HandleFunc("PATCH /api/user/{id}", api.patchUser)
	router.HandleFunc("GET /api/user/{id}", api.getUser)
	router.HandleFunc("DELETE /api/user", api.deleteUser)

	router.HandleFunc("POST /api/license, port", api.postLicense)
	router.HandleFunc("PATCH /api/license/{id}", api.patchLicense)
	router.HandleFunc("GET /api/license/{id}", api.getLicense)

	router.HandleFunc("GET /api/license/check/", api.checkLicense)

	router.HandleFunc("POST /api/create-login-code", api.createLoginCode)

	return router
}

func (api *API) checkHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Del("Content-Type")
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
		w.Header().Add("Content-Type", "application/json")
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
	verificationStore := store.CreateVerificationStore(2*time.Minute, 20*time.Minute)

	api := &API{
		db:                db,
		verificationStore: verificationStore,
	}
	router := api.setupRouter()
	handler := middlewareSetup(router)

	server := &http.Server{
		Handler: handler,
		Addr:    ":" + cfg.Port,
	}

	log.Println("Running on port:", cfg.Port)
	server.ListenAndServe()
}
