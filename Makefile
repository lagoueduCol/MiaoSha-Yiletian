all: build

proto: application/api/rpc/event.proto
	protoc --go_out=plugins=grpc:./ application/api/rpc/event.proto

build:
	go build -o bin/seckill main.go

clean:
	rm bin/seckill

runApi:
	ulimit -n 1000000
	./bin/seckill api -c ./config/seckill.toml

ab:
	ulimit -n 1000000
	ab -n 1000000 -c 50 http://localhost:8080/event/list

ab-k:
	ulimit -n 1000000
	ab -n 1000000 -c 50 -k http://localhost:8080/event/list

bench:
	ulimit -n 1000000
	./bin/seckill bench -C 16 -r 1000000 -u http://localhost:8080/event/list

bench-k:
	ulimit -n 1000000
	./bin/seckill bench -C 16 -k -r 1000000 -u http://localhost:8080/event/list

test:
	go test -bench=. -cover ./...

.PHONY: clean build all
