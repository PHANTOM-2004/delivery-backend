go_files = $(shell go list ./... | grep -v /vendor/)

.PHONY: test format rundt run build

# 检查环境变量，设置默认目标
ifeq ($(GO_SERVICE_ENV),DOCKER_TEST)
    DEFAULT_GOAL := rundt
else ifeq ($(GO_SERVICE_ENV),DOCKER_DEPLOY)
    DEFAULT_GOAL := runserver
else
    DEFAULT_GOAL := run
endif

# 如果没有指定目标，则将默认目标设置为 DEFAULT_GOAL
ifeq ($(MAKECMDGOALS),)
    .DEFAULT_GOAL := $(DEFAULT_GOAL)
endif

check-swagger:
	@which swagger || (go install github.com/go-swagger/go-swagger/cmd/swagger@latest)

swagger: check-swagger
	@swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: swagger
	@swagger serve -F=swagger swagger.yaml

format:
	@go fmt $(go_files)

test:
	@go vet $(go_files)
	@go test -race $(go_files)

build:
	@mkdir -p target
	@go build -o target ./...

run:
	@go run -v . --LocalDebug

rundt:
	@go run -v . --DockerTest

runserver:
	@go run -v . --DockerDeploy
