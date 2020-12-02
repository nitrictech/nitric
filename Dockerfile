FROM golang as build

RUN apt update
RUN apt install -y protobuf-compiler

RUN go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

WORKDIR /

COPY . .

RUN ./build.sh

FROM alpine

COPY --from=build /bin/membrane /membrane
RUN chmod +xr /membrane

ENTRYPOINT [ "/membrane" ]

