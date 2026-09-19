.PHONY: run build test clean docker-build docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test -v -race -cover ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t hirelly-backend:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down
