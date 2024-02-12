PHONY: fedora
fedora:
	sudo go run cmd/anubis.go --target fedora --config ./config.yaml

PHONY: debug-fedora
debug-fedora:
	sudo dlv debug cmd/anubis.go -- --target fedora --config ./config.yaml

PHONY: clean
clean:
	rm build/*

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/anubis-linux-amd64 cmd/anubis.go

build-darwin: 
	GOOS=darwin GOARCH=arm64 go build -o build/anubis-darwin-arm64 cmd/anubis.go

PHONY: fedora-github
fedora-github:
	sudo go run cmd/anubis.go --target fedora

PHONY: test
test:
	go test ./...
