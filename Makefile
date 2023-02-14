all: build proto

build:
	go build -o main .

proto:
	protoc -I ./proto --go-grpc_out=./proto ./proto/healthcheck.proto