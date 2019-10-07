## Fetch dependencies
fetch:
	GO111MODULE=on go get -v ./...

run-worker:
	GO111MODULE=on go run worker/new_product/*.go -consumers 2

run-server:
	GO111MODULE=on go run server/main.go

## Run tests
test:
	GO111MODULE=on go test -race -v ./...

## Run tests with coverage
test-cover:
	GO111MODULE=on go test -coverprofile=cover.out -race -v ./...

# .PHONY: clean
# ## Remove binary
# clean:
# 	if [ -f $(BINARY_NAME) ]; then rm -f $(BINARY_NAME); fi
