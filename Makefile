BIN=$(CURDIR)/bin

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build -o $(BIN)/vvm ./cmd/cli

build-lambda:
	GOOS=linux go build -o $(BIN)/main ./cmd/aws/; cd $(BIN); zip function.zip main; rm main

build-lambda-chatops:
	GOOS=linux go build -o $(BIN)/main ./cmd/aws/chatops/; cd $(BIN); zip chatops_function.zip main; rm main

.PHONY: clean
clean:
	rm -rf $(BIN)