.PHONY: run

gateway:
	JWT_SECRET=local-development-secret go run ./cmd/gateway

backend:
	go run ./cmd/backend