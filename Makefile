build:
	go build ./cmd/phpdocker-gen

test-all:
	go test -coverprofile=coverage.txt -v ./internal/... ./pkg/...