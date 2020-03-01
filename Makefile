generate:
	protoc --go_out=plugins=grpc:internal/grpc api/api.proto -I .
lint:
	golangci-lint run
test:
	go test -v ./...
	go test -v -cover ./...
	go test -v -race ./...
build-server:
	go build -o abf-srv github.com/Lefthander/otus-go-antibruteforce/cmd/server
build-client:
	go build -o abf-ctl github.com/Lefthander/otus-go-antibruteforce/cmd/client