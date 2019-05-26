.PHONY: deps build

deps:
	go get -u ./...

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/lambda ./lambda
