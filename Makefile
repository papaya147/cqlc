tidy:
	go mod tidy

test:
	go test -v -cover ./...

build:
	env GOOS=linux CGO_ENABLED=0 go build -o go-cassandra-codegen-app main.go

.PHONY: test