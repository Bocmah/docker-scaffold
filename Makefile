build:
	@go build ./cmd/phpdocker-gen

test-all:
	@go test -coverprofile=coverage.txt -v ./cmd/... ./internal/... ./pkg/...

generate:
	@go generate ./...