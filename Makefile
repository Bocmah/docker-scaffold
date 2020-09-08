build:
	go build ./cmd/phpdocker-gen

test-all:
	go test -v ./internal/... ./pkg/...