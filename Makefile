.PHONY: gateway backend orders health req_user req_order

JWT_SECRET ?= local-development-secret
JWT_ISSUER ?= local-auth
JWT_AUDIENCE ?= api-gateway
USERS_URL ?= http://localhost:8081
ORDERS_URL ?= http://localhost:8082


gateway:
	JWT_SECRET=$(JWT_SECRET) \
	JWT_ISSUER=$(JWT_ISSUER) \
	JWT_AUDIENCE=$(JWT_AUDIENCE) \
	USERS_URL=$(USERS_URL) \
	ORDERS_URL=$(ORDERS_URL) \
	go run ./cmd/gateway

backend:
	go run ./cmd/backend

orders:
	go run ./cmd/orders

health:
	curl -i http://localhost:8080/health

ready:
	curl -i http://localhost:8080/ready 

fmt:
	gofmt -w .

metrics:
	curl http://localhost:8080/metrics

tests:
	go test ./...

req_user:
	@TOKEN=$$( \
		JWT_SECRET=$(JWT_SECRET) \
		JWT_ISSUER=$(JWT_ISSUER) \
		JWT_AUDIENCE=$(JWT_AUDIENCE) \
		go run ./cmd/token \
	); \
	curl -i \
		-H "Authorization: Bearer $$TOKEN" \
		http://localhost:8080/api/users/42

req_order:
	@TOKEN=$$( \
		JWT_SECRET=$(JWT_SECRET) \
		JWT_ISSUER=$(JWT_ISSUER) \
		JWT_AUDIENCE=$(JWT_AUDIENCE) \
		go run ./cmd/token \
	); \
	curl -i \
		-H "Authorization: Bearer $$TOKEN" \
		http://localhost:8080/api/orders/99
