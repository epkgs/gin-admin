.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')


SERVER_BIN  	= ginadmin
APP_VERSION     = v1.0.0
GIT_COUNT 		= $(shell git rev-list --all --count)
GIT_HASH        = $(shell git rev-parse --short HEAD)
RELEASE_TAG     = $(APP_VERSION).$(GIT_COUNT).$(GIT_HASH)

CONFIG_DIR       = ./configs
RUN_MODE         = dev
START_ARGS       = -c $(CONFIG_DIR) -m $(RUN_MODE)

all: start

start:
	@go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" main.go start $(START_ARGS)

build:
	@go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)

# build-linux:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux-musl" CXX="zig c++ -target x86_64-linux-musl" CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)_linux_amd64

# build-mac:
# 	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC="zig cc -target x86_64-macos-gnu" CXX="zig c++ -target x86_64-macos-gnu" CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)_macos_amd64

build-win:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" CXX="zig c++ -target x86_64-windows-gnu" CGO_CFLAGS="-D_LARGEFILE64_SOURCE" go build -ldflags "-w -s -X main.VERSION=$(RELEASE_TAG)" -o $(SERVER_BIN)_windows_amd64.exe

# go install github.com/swaggo/swag/cmd/swag@latest
swagger:
	@swag init --parseDependency --generalInfo ./main.go --output ./internal/swagger

# https://github.com/OpenAPITools/openapi-generator
openapi:
	docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli generate -i /local/internal/swagger/swagger.yaml -g openapi -o /local/internal/swagger/v3

clean:
	rm -rf data $(SERVER_BIN)

