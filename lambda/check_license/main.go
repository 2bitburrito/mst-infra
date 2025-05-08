package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LicenceObj struct {
	LicenseKey string        `json:"license_key"`
	TokenObj   CognitoObject `json:"token"`
	MachineID  string        `json:"machine_id"`
}

type Response struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

var db *sql.DB

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

func HandleRequest(req events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	if db == nil {
		return createErrorResponse(http.StatusInternalServerError, "Database connection is not established")
	}

	var request LicenceObj
	if err := json.Unmarshal([]byte(req.Body), &request); err != nil {
		fmt.Printf("Error unmarshalling request body: %v", err)
		return createErrorResponse(http.StatusBadRequest, "Invalid request body")
	}
	var dbResponse LicenceObj

	err := db.QueryRow("SELECT licence_key, machine_id FROM licenses WHERE license_key = $1",
		request.LicenseKey).Scan(&dbResponse.LicenseKey, &dbResponse.MachineID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No matching license found")
			return createErrorResponse(http.StatusUnauthorized, "Invalid license key or machine ID")
		}
		return createErrorResponse(http.StatusInternalServerError, err.Error())
	}
	fmt.Printf("Found matching license: %v\n", dbResponse)

	res := Response{
		Valid: dbResponse.MachineID == request.MachineID,
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("Error marshalling response: %v", err)
		return createErrorResponse(http.StatusInternalServerError, err.Error())
	}

	// Return the response to prompt log out if needed
	if res.Valid {
		err = RevokeAccess(request.TokenObj)
		if err != nil {
			return createErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Error revoking access: %v", err))
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(resJSON),
	}
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
