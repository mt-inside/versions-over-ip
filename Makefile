setup:
	./setup.sh

.PHONY: generate
generate:
	./generate.sh


validate:
	go vet ./...
	golangci-lint run ./...
	go test ./...

client-build: validate
	go build cmd/client/client.go

client-run: validate
	go run cmd/client/client.go


server-run: validate
	go run cmd/server/server.go
