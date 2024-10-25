go_files = $(shell go list ./... | grep -v /vendor/)

.PHONY: test format

all: serve-swagger

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

