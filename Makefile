.PHONY: run

gateway:
	API_TOKEN=dev-secret go run ./cmd/gateway

backend:
	go run ./cmd/backend