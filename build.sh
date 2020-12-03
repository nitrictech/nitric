# build the go project
CGO_ENABLED=0 GOOS=linux go build -o bin/membrane -ldflags="-extldflags=-static" main.go