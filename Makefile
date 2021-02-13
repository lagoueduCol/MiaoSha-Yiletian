all: build

proto: application/api/rpc/event.proto
	protoc --go_out=plugins=grpc:./ application/api/rpc/event.proto

build:
	go build -o bin/seckill main.go

clean:
	rm bin/seckill

runApi:
	./bin/seckill api -c ./config/seckill.toml

.PHONY: clean build all
