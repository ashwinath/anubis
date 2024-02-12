PHONY: fedora
fedora:
	sudo go run cmd/anubis.go --target fedora --config ./config.yaml

PHONY: debug-fedora
debug-fedora:
	sudo dlv debug cmd/anubis.go -- --target fedora --config ./config.yaml

PHONY: clean
clean:
	rm build/*

build: clean
	go build -o build/anubis cmd/anubis.go

PHONY: fedora-github
fedora-github:
	sudo go run cmd/anubis.go --target fedora
