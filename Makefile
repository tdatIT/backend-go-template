APP_NAME=app
CMD_PATH=./cmd/main.go

.PHONY: build run test lint fmt tidy docker-build docker-up docker-down

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

tidy:
	go mod tidy

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down
