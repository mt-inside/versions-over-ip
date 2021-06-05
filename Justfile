default:
	@just --list

setup:
	./build/setup.sh

generate:
	./build/generate.sh

quick:
	./build/quick.sh

validate:
	go fmt ./...
	go vet ./...
	golint -set_exit_status ./...
	golangci-lint run ./...
	go test ./...

build CMD: #generate validate
	go build ./cmd/{{CMD}}

run CMD: #generate validate
	go run ./cmd/{{CMD}}

