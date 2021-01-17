all: build

proto: application/api/rpc/event.proto
	protoc --go_out=plugins=grpc:./ application/api/rpc/event.proto

build:
	go build -o bin/seckill main.go

clean:
	rm bin/seckill

.PHONY: clean build all
