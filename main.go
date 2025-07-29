package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/2bitburrito/mst-infra/config"
	Queries "github.com/2bitburrito/mst-infra/db/sqlc"
	"github.com/2bitburrito/mst-infra/store"

	_ "github.com/lib/pq"
)

type API struct {
	db                *sql.DB
	queries           Queries.Queries
	verificationStore *store.VerificationStore
	config            config.Config
}

func (api *API) setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/healthz/", api.checkHealth)

	router.Handle("POST /api/user", api.apiMiddleware(http.HandlerFunc(api.postUser)))
	router.Handle("POST /api/cognito-user", api.apiMiddleware(http.HandlerFunc(api.postCognitoUser)))
	router.Handle("PATCH /api/user/{id}", api.apiMiddleware(http.HandlerFunc(api.patchUser)))
	router.Handle("GET /api/user/{id}", api.apiMiddleware(http.HandlerFunc(api.getUser)))
	router.Handle("DELETE /api/user", api.apiMiddleware(http.HandlerFunc(api.deleteUser)))
	router.Handle("GET /api/user/is-beta/{email}", api.apiMiddleware(http.HandlerFunc(api.checkUserIsBeta)))

	router.Handle("PUT /api/email-all-beta-users", api.apiMiddleware(http.HandlerFunc(api.emailAllBetaUsers)))
	router.Handle("PUT /api/email-select-beta-users", api.apiMiddleware(http.HandlerFunc(api.emailSelectBetaUsers)))
	router.Handle("PUT /api/test-email-beta-users", api.apiMiddleware(http.HandlerFunc(api.testEmails)))

	router.Handle("POST /api/license", api.apiMiddleware(http.HandlerFunc(api.postLicense)))
	router.Handle("PATCH /api/license/{id}", api.apiMiddleware(http.HandlerFunc(api.patchLicense)))
	router.Handle("GET /api/license/{id}", api.apiMiddleware(http.HandlerFunc(api.getLicense)))
	router.Handle("POST /api/license/check", api.desktopAppRouterMiddleware(http.HandlerFunc(api.checkLicense)))

	router.Handle("POST /api/create-login-code", api.apiMiddleware(http.HandlerFunc(api.createLoginCode)))
	router.Handle("POST /api/check-login-code", http.HandlerFunc(api.checkLoginCodeAndCreateJWT))

	router.Handle("POST /api/latest-binaries", api.apiMiddleware(http.HandlerFunc(api.insertLatestBinaries)))
	router.Handle("GET /api/latest-binaries/{platform}/{arch}", api.apiMiddleware(http.HandlerFunc(api.getLatestBinaries)))

	return router
}

func (api *API) checkHealth(w http.ResponseWriter, req *http.Request) {
	w.Header().Del("Content-Type")
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (api *API) desktopAppRouterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (api *API) apiMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := map[string]bool{
			"http://localhost:3000":      true,
			"https://metasoundtools.com": true,
		}

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
		if apiKey != api.config.ApiKey {
			returnJsonError(w, "Unauthorized Api Key from: "+r.RemoteAddr, http.StatusInternalServerError)
			return
		}
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
		queries:           *queries,
		verificationStore: verificationStore,
		config:            *cfg,
	}
	if api.db == nil {
		panic("DB IS NIL")
	}

	router := api.setupRouter()

	server := &http.Server{
		Handler: router,
		Addr:    ":" + cfg.Port,
	}

	log.Println("Running on port:", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
