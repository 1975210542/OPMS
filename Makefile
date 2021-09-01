.PHONY: main mod run

mod:
	go mod tidy
	go mod vendor

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

run:
	./main