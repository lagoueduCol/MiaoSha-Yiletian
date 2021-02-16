all: build

proto: application/api/rpc/event.proto
	protoc --go_out=plugins=grpc:./ application/api/rpc/event.proto

build:
	go build -o bin/seckill main.go

clean:
	rm bin/seckill

runApi:
	./bin/seckill api -c ./config/seckill.toml

test:
	go test -bench=. -cover ./...

bench:
	ulimit -n 10000000
	./bin/seckill bench -C 50 -r 1000000 -u http://localhost:8080/event/list

ab:
	ulimit -n 1000000
	ab -n 1000000 -c 50 -k -I -b 10240 http://localhost:8080/event/list

.PHONY: clean build all
