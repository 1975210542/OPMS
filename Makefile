.PHONY: main mod run

mod:
	go mod tidy
	go mod vendor

build:
	go build -o main main.go
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go
build-win:
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o main main.go

run:
	./main
run-linux:
	nohup ./main &

stop-linux:
	lsof -i :8088
