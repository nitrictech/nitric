mkdir -p ./interfaces/

# Generate the client/server go code for the gcp ambassador
protoc --go_out=./interfaces \
    --go-grpc_out=./interfaces \
    -I ./contracts/proto/ \
    ./contracts/proto/**/*.proto

# build the go project
CGO_ENABLED=0 GOOS=linux go build -o bin/ambassador -ldflags="-extldflags=-static" main.go