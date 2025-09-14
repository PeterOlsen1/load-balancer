.PHONY: all test clean clean-containers test-work test-rps lint

REQUESTS ?= 100

SECONDS ?= 3
RPS ?= 100

all:
	go run ./pkg/main.go

test:
	go run ./test/main.go -requests=$(REQUESTS)

test-rps:
	go run ./test/main.go -seconds=${SECONDS} -rps=${RPS}

test-wrk:
	wrk -t10 -c1000 -d10s http://localhost:8080

up-fds:
	ulimit -n $((1 << 16))

clean:
	rm -rf ./logs

clean-containers:
	docker stop $$(docker ps -q)

# increase file descriptors
# ulimit -n 65536

# increase max current connections (mac os)
# sudo sysctl -w kern.ipc.somaxconn=1024

# Test with wrk
# wrk -t10 -c500 -d20s http://localhost:8080
