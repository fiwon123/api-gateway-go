FROM golang:1.26.7-alpine AS builder

RUN apk add --no-cache \
    ca-certificates \
    git

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG SERVICE=gateway

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -o /out/service ./cmd/${SERVICE}

FROM alpine:3.20

RUN apk add --no-cache ca-certificates \
    && adduser -D appuser

USER appuser

COPY --from=builder /out/service /service

ENTRYPOINT ["/service"]
