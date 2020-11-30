setup:
	./setup.sh

generate:
	./generate.sh
	@echo now fix import in cmd/client/versionsclient/versions_client.go

client-run:
	go run cmd/client/client.go

server-run:
	go run cmd/server/server.go
