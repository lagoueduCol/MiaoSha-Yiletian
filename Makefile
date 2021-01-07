all: build

build:
	go build -o seckill main.go

clean:
	rm seckill

.PHONY: clean build all
