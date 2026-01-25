.PHONY: dev build run migrate seed clean

dev:
	air

build:
	go build -o bin/main main.go

run: build
	./bin/main

migrate:
	go run main.go -migrate

seed:
	go run main.go -seed

clean:
	rm -rf bin tmp
