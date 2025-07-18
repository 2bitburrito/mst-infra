package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/2bitburrito/mst-infra/utils"
)

func returnJsonError(w http.ResponseWriter, e string, statusCode int) {
	log.Println(e)
	rtnMap := utils.JsonReturn{
		Success: false,
		Error:   e,
		Value:   nil,
	}
	dat, err := json.Marshal(rtnMap)
	if err != nil {
		log.Println("error Marshalling error response")
		http.Error(w, `{"success":false,"error":"internal error"}`, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(dat)
}
