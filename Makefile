bin/rcc: main.go cmd/*.go pkg/badge/*.go
	@mkdir -p bin
	go build -o bin/rcc main.go

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test
test:
	@go test ./pkg/...
