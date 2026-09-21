# Simple API Gateway

This project is a simple API Gateway using, authentiction, routing, reverse proxy, shutdown gracefully, etc.

## Pre-Requisites
- Docker Engine 29.8.0
- Docker Compose v5.5.1
- GNU Make 4.4.1
- curl 8.15.0 
- Go 1.26.7

## How To Use
### Start
```
docker compose up --build
```
### Stop
```
docker compose down
```

## Make
I've used curl for requests and I've automated using Make to keep it simple and easy to visualize. You can use any of these programs: Postman, Insomia, or Bruno for requests.

### Health
Check if gateway is working
```
make health
```

### Ready
Check if gateway, user backend and order backend is working
```
make ready
```

### Request User for testing
```
make req_user
```


### Request Order for testing
```
make req_order
```