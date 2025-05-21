.PHONY: host
host:
	mkdir -p host/cmd
	go run ./scripts/build/build.go > host/cmd/main.go
	go build -o bin/host host/cmd/main.go
