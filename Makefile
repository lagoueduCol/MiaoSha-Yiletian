all: build

build:
	go build -o bin/seckill main.go

clean:
	rm seckill

.PHONY: clean build all
