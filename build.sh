# build the go project
# Normally we would make this a static build for container compatibility
# However we want to use plugins
CGO_ENABLED=1 GOOS=linux go build -o bin/membrane main.go