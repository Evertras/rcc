GO_FILES=$(shell find . -name "*.go")

bin/rcc: $(GO_FILES)
	@mkdir -p bin
	go build -o bin/rcc ./cmd/rcc

bin/lambda: $(GO_FILES)
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -o bin/lambda ./cmd/lambda

.PHONY: fmt
fmt: node_modules
	go fmt ./...
	npx prettier --write .

.PHONY: test
test:
	@go test ./pkg/...

.PHONY: docker
docker:
	docker build -t evertras/rcc:latest .

node_modules: package.json package-lock.json
	npm install
	@touch node_modules

################################################################################
# Local tooling
#
# This section contains tools to download to the local ./bin directory for easy
# local use.  The .envrc file makes adding the local ./bin directory to our path
# simple, so we can use tools here without having to install them globally as if
# they actually were global.
#
# For now we only support Linux 64 bit and MacOS for simplicity
ifeq ($(shell uname), Darwin)
OS_URL := darwin
else
OS_URL := linux
endif

bin/terraform:
	mkdir -p bin
	curl -Lo bin/terraform.zip https://releases.hashicorp.com/terraform/1.3.2/terraform_1.3.2_$(OS_URL)_amd64.zip
	cd bin && unzip terraform.zip
	rm bin/terraform.zip
