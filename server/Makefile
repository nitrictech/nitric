.PHONY: host
host:
	mkdir -p host/cmd
	go run ./scripts/build/build.go > host/cmd/main.go
	CGO_ENABLED=0 go build -o bin/host host/cmd/main.go
