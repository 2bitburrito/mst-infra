package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

type JsonErrReturn struct {
	Error string `json:"error"`
	Msg   string `json:"msg,omitempty"`
}

func logError(message string, e error, extraData ...any) {
	err := fmt.Errorf("ERROR: %v - %w \nDATA: %+v", message, e, extraData)
	sentry.CaptureException(err)
	log.Println(err)
}

func returnJsonError(w http.ResponseWriter, e string, statusCode int, msg ...string) {
	if statusCode > 499 {
		log.Printf("Responding with 5XX error message: %s\n", msg)
	}
	sentry.CaptureException(fmt.Errorf("ERROR: %s \nCLIENT MSG: %s", e, msg))
	log.Println(e)
	var message string
	if len(msg) > 0 {
		message = msg[0]
	}
	rtnMap := JsonErrReturn{
		Error: e,
		Msg:   message,
	}
	respondWithJSON(w, statusCode, rtnMap)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
