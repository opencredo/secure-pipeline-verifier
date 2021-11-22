BIN=$(CURDIR)/bin

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build -o $(BIN)/main ./cmd/cli
build-lambda: build; zip $(BIN)/function.zip $(BIN)/main; rm $(BIN)/main

.PHONY: clean
clean:
	rm -rf $(BIN)