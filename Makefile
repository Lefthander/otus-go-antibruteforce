generate:
	protoc --go_out=plugins=grpc:internal/grpc api/api.proto -I .
lint:
	golangci-lint run
test:
	go test -v ./...
	go test -v -cover ./...
	go test -v -race ./...