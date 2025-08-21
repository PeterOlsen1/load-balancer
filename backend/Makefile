.PHONY: all test clean clean-containers

all:
	go run ./pkg/main.go

test:
	go run ./test/main.go

clean:
	rm -rf ./logs

clean-containers:
	docker stop $$(docker ps -q)