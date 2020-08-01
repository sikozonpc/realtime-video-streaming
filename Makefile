.PHONY: run-api

PORT := 7777
ADDR := "localhost"
DB_HOST := "localhost"
DB_USER := "postgres"
DB_NAME := "streaming_server"
DB_PORT := "5432"
DB_PASSWORD := "batata123"

run-api:
	GO111MODULE=on go build -o goproject && ./goproject -port=${PORT} -addr=${ADDR} -db-host=${DB_HOST} -db-user=${DB_USER} -db-name=${DB_NAME} -db-port=${DB_PORT} -db-password=${DB_PASSWORD}