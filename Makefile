.PHONY: gateway backend orders health req_user req_order

JWT_SECRET ?= local-development-secret

gateway:
	JWT_SECRET=$(JWT_SECRET) go run ./cmd/gateway

backend:
	go run ./cmd/backend

orders:
	go run ./cmd/orders

health:
	curl -i http://localhost:8080/health

req_user:
	@TOKEN=$$(JWT_SECRET=$(JWT_SECRET) go run ./cmd/token); \
	curl -i \
		-H "Authorization: Bearer $$TOKEN" \
		http://localhost:8080/api/users/42

req_order:
	@TOKEN=$$(JWT_SECRET=$(JWT_SECRET) go run ./cmd/token); \
	curl -i \
		-H "Authorization: Bearer $$TOKEN" \
		http://localhost:8080/api/orders/99
