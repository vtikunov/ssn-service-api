GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on

SERVICE_NAME=ssn-service-api
SERVICE_FACADE_NAME=ssn-service-facade
SERVICE_PATH=ozonmp/ssn-service-api

PGV_VERSION:="v0.6.1"
BUF_VERSION:="v0.56.0"

OS_NAME=$(shell uname -s)
OS_ARCH=$(shell uname -m)
GO_BIN=$(shell go env GOPATH)/bin
BUF_EXE=$(GO_BIN)/buf$(shell go env GOEXE)

ifeq ("NT", "$(findstring NT,$(OS_NAME))")
OS_NAME=Windows
endif

.PHONY: run
run:
	go run cmd/grpc-server/main.go

.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go test -v -race -timeout 30s -coverprofile cover.out ./...
	go tool cover -func cover.out | grep total | awk '{print $$3}'


# ----------------------------------------------------------------

.PHONY: generate
generate: .generate-install-buf .generate-go .generate-python .generate-finalize-go .generate-finalize-python

.PHONY: generate
generate-go: .generate-install-buf .generate-go .generate-finalize-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) --create-dirs -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")

.generate-go: .build-go-tools
	$(BUF_EXE) generate
	./bin/add-debug-header-to-swagger swagger/ozonmp/ssn_service_api/v1/ssn_service_api.swagger.json
	protoc -o api/ozonmp/ssn_service_api/v1/ssn_service_api.desc  -I api --include_imports api/ozonmp/ssn_service_api/v1/ssn_service_api.proto

.generate-python:
	$(BUF_EXE) generate --template buf.gen.python.yaml

.generate-finalize-go:
	mv pkg/$(SERVICE_NAME)/github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME)/* pkg/$(SERVICE_NAME)
	rm -rf pkg/$(SERVICE_NAME)/github.com/
	cd pkg/$(SERVICE_NAME) && ls go.mod || (go mod init github.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

.generate-finalize-python:
	find pypkg/ssn-service-api -type d -not -path "pypkg/ssn-service-api/dist" -exec touch {}/__init__.py \;

.PHONY: generate-facade
generate-facade: .generate-install-buf .generate-facade

.generate-facade:
	$(BUF_EXE) generate --template buf.gen.facade.yaml
	mv pkg/$(SERVICE_FACADE_NAME)/github.com/$(SERVICE_PATH)/pkg/$(SERVICE_FACADE_NAME)/* pkg/$(SERVICE_FACADE_NAME)
	rm -rf pkg/$(SERVICE_FACADE_NAME)/github.com/
	cd pkg/$(SERVICE_FACADE_NAME) && ls go.mod || (go mod init github.com/$(SERVICE_PATH)/pkg/$(SERVICE_FACADE_NAME) && go mod tidy)

# ----------------------------------------------------------------

.PHONY: deps
deps: deps-go .deps-python

.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

.deps-python:
	python -m pip install grpcio-tools grpclib protobuf

.PHONY: build
build: generate .build

.PHONY: build-go
build-go: generate-go .build

.build:
	go mod download && CGO_ENABLED=0  go build \
		-tags='no_mysql no_sqlite3' \
		-ldflags=" \
			-X 'github.com/$(SERVICE_PATH)/internal/config.version=$(VERSION)' \
			-X 'github.com/$(SERVICE_PATH)/internal/config.commitHash=$(COMMIT_HASH)' \
		" \
		-o ./bin/grpc-server$(shell go env GOEXE) ./cmd/grpc-server/main.go

.PHONY: build-go-retranslator
build-go-retranslator:
	go mod download && CGO_ENABLED=0  go build \
    		-tags='no_mysql no_sqlite3' \
    		-ldflags=" \
    			-X 'github.com/$(SERVICE_PATH)/internal/retranslator/config.version=$(VERSION)' \
    			-X 'github.com/$(SERVICE_PATH)/internal/retranslator/config.commitHash=$(COMMIT_HASH)' \
    		" \
    		-o ./bin/retranslator$(shell go env GOEXE) ./cmd/retranslator/main.go

.PHONY: build-go-facade-grpc
build-go-facade-grpc: generate-facade .build-go-facade-grpc

.PHONY: build-go-bot
build-go-bot:
	go mod download && CGO_ENABLED=0  go build \
    		-tags='no_mysql no_sqlite3' \
    		-ldflags=" \
    			-X 'github.com/$(SERVICE_PATH)/internal/bot/config.version=$(VERSION)' \
    			-X 'github.com/$(SERVICE_PATH)/internal/bot/config.commitHash=$(COMMIT_HASH)' \
    		" \
    		-o ./bin/bot$(shell go env GOEXE) ./cmd/bot/main.go

.PHONY: build-go-facade
build-go-facade:
	go mod download && CGO_ENABLED=0  go build \
    		-tags='no_mysql no_sqlite3' \
    		-ldflags=" \
    			-X 'github.com/$(SERVICE_PATH)/internal/facade/config.version=$(VERSION)' \
    			-X 'github.com/$(SERVICE_PATH)/internal/facade/config.commitHash=$(COMMIT_HASH)' \
    		" \
    		-o ./bin/facade$(shell go env GOEXE) ./cmd/facade/main.go

.build-go-facade-grpc:
	go mod download && CGO_ENABLED=0  go build \
    		-tags='no_mysql no_sqlite3' \
    		-ldflags=" \
    			-X 'github.com/$(SERVICE_PATH)/internal/facade/grpc/config.version=$(VERSION)' \
    			-X 'github.com/$(SERVICE_PATH)/internal/facade/grpc/config.commitHash=$(COMMIT_HASH)' \
    		" \
    		-o ./bin/facade-grpc$(shell go env GOEXE) ./cmd/facade-grpc-server/main.go

.build-go-tools:
	go mod download && CGO_ENABLED=0  go build \
    		-tags='no_mysql no_sqlite3' \
    		-ldflags=" \
    			-X 'github.com/$(SERVICE_PATH)/internal/config.version=$(VERSION)' \
    			-X 'github.com/$(SERVICE_PATH)/internal/config.commitHash=$(COMMIT_HASH)' \
    		" \
    		-o ./bin/add-debug-header-to-swagger$(shell go env GOEXE) ./cmd/tools/add-debug-header-to-swagger/main.go
