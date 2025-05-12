package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var db *sql.DB

type Response struct {
	Valid    bool   `json:"valid"`
	Messages string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}
type proxyResponseWriter struct {
	headers http.Header
	body    []byte
	status  int
}

func init() {
	dbConnectionString, err := CheckEnv()
	if err != nil {
		createErrorResponse(http.StatusInternalServerError, err.Error())
		panic(err)
	}

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		createErrorResponse(http.StatusInternalServerError, err.Error())
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

func HandleRequest(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	router := setupRouter()

	if db == nil {
		return createErrorResponse(http.StatusInternalServerError, "Database connection is not established")
	}
	w := &proxyResponseWriter
	events.APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{}
}

func createErrorResponse(httpStatusCode int, body string) events.APIGatewayProxyResponse {
	responseBody := Response{
		Valid: false,
		Error: body,
	}
	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Error marshalling response: %v\n", err)
	}
	return events.APIGatewayProxyResponse{
		StatusCode: httpStatusCode,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseJSON),
	}
}

func main() {
	lambda.Start(HandleRequest)
}
