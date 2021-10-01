29 lines (23 sloc)  788 Bytes

BIN=$(CURDIR)/bin

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build -o $(BIN)/vvw

.PHONY: clean
clean:
	rm -rf $(BIN)