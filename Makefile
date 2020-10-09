.PHONY: run-api, run-client

PORT := 7777
ADDR := "localhost"

run-api:
	GO111MODULE=on go build -o goproject && ./goproject -port=${PORT} -addr=${ADDR}
run-client:
	cd client && npm run start