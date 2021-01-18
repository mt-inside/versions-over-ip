setup:
	./setup.sh

.PHONY: generate
generate:
	./generate.sh
	@echo now fix import in cmd/client/versionsclient/versions_client.go


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
