BIN=$(CURDIR)/bin

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build-cli
build-cli:
	go build -o $(BIN)/vvw ./cmd/cli

.PHONY: clean
clean:
	rm -rf $(BIN)

.PHONY: build-lambda
build-lambda:
	echo Building and compressing Lambda function
	cd cmd/aws; \
    GOOS=linux go build -o main main.go; \
    zip "function.zip" main