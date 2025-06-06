package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	Queries "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/server/api/config"
	"github.com/2bitburrito/mst-infra/server/api/store"

	_ "github.com/lib/pq"
)

type API struct {
	db                *sql.DB
	ctx               context.Context
	queries           Queries.Queries
	verificationStore *store.VerificationStore
}

func (api *API) setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/api/healthz", api.checkHealth)

	router.HandleFunc("POST /api/user", api.postUser)
	router.HandleFunc("POST /api/cognito-user", api.postCognitoUser)
	router.HandleFunc("PATCH /api/user/{id}", api.patchUser)
	router.HandleFunc("GET /api/user/{id}", api.getUser)
	router.HandleFunc("DELETE /api/user", api.deleteUser)

	router.HandleFunc("POST /api/license", api.postLicense)
	router.HandleFunc("PATCH /api/license/{id}", api.patchLicense)
	router.HandleFunc("GET /api/license/{id}", api.getLicense)

	router.HandleFunc("POST /api/license/check/", api.checkLicense)

	router.HandleFunc("POST /api/create-login-code", api.createLoginCode)
	router.HandleFunc("POST /api/check-login-code", api.checkLoginCodeAndCreateJWT)

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
		allowedOrigins := map[string]bool{
			"http://localhost:3000":      true,
			"http://localhost:5173":      true,
			"https://metasoundtools.com": true,
		}
		// CORS config
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		w.Header().Add("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Check API Key:
		apiKey := r.Header.Get("X-API-Key")
		log.Println("API Key received:", apiKey)
		log.Println("API Key in config:", cfg.ApiKey)
		if apiKey != cfg.ApiKey {
			returnJsonError(w, "Unauthorixed Api Key", http.StatusInternalServerError)
			return
		}
		log.Println("Successful Api Key Match")

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg, _ := config.LoadConfig()
	verificationStore := store.CreateVerificationStore(1*time.Minute, 10*time.Minute)

	var db *sql.DB
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

	queries := Queries.New(db)

	api := &API{
		db:                db,
		ctx:               context.TODO(),
		queries:           *queries,
		verificationStore: verificationStore,
	}
	if api.db == nil {
		panic("DB IS NIL")
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
