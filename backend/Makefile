.PHONY: all test clean clean-containers

REQUESTS ?= 100

all:
	go run ./pkg/main.go

test:
	go run ./test/main.go -requests=$(REQUESTS)

clean:
	rm -rf ./logs

clean-containers:
	docker stop $$(docker ps -q)