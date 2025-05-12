package main

import (
	"errors"
	"fmt"
	"os"
)

func CheckEnv() (pqConnectionSting string, error error) {
	DB_URL := os.Getenv("DB_URL")
	DB_PORT := os.Getenv("DB_PORT")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")
	if len(DB_URL) == 0 {
		return "", errors.New("DB_URL is not set")
	}
	if len(DB_PASSWORD) == 0 {
		return "", errors.New("DB_PASSWORD is not set")
	}
	if len(DB_USER) == 0 {
		return "", errors.New("DB_USER is not set")
	}
	if len(DB_NAME) == 0 {
		return "", errors.New("DB_NAME is not set")
	}

	pqConnectionSting = fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", DB_URL, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	return pqConnectionSting, nil
}

/*
NOTE: This is from the secrets manager docs if the above doesn't work

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)
	secretName := "rds-db-credentials/cluster-X2PRYBL4NAZ43D7SEN526PSMJI/mst_admin/1745813016480"
	region := "us-west-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}
	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		log.Fatal(err.Error())
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
*/
