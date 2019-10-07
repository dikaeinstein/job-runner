build-server:
	GO111MODULE=on go build -o server/server server/main.go

build-worker:
	GO111MODULE=on go build -o worker/new_product/worker worker/new_product/*.go

build: build-worker build-server

## Fetch dependencies
fetch:
	GO111MODULE=on go get -v ./...

## Run worker
run-worker:
	GO111MODULE=on go run worker/new_product/*.go -consumers 2

## Run server
run-server:
	GO111MODULE=on go run server/main.go

## Run tests
test:
	GO111MODULE=on go test -race -v ./...

## Run tests with coverage
test-cover:
	GO111MODULE=on go test -coverprofile=cover.out -race -v ./...

clean-server:
	rm server/server
clean-worker:
	rm worker/new_product/worker
.PHONY: clean
## Remove binary
clean: clean-server clean-worker

