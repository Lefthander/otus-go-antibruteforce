generate:
	protoc --go_out=plugins=grpc:internal/grpc api/api.proto -I .
lint:
	golangci-lint run
test:
	go test -v ./...
	go test -v -cover ./...
	go test -v -race ./...

.PHONY: integration_test
integration_test:
	go test -v ./tests	
build-server:
	go build -o ./bin/abf-srv github.com/Lefthander/otus-go-antibruteforce/cmd/server
build-client:
	go build -o ./bin/abf-ctl github.com/Lefthander/otus-go-antibruteforce/cmd/client

undeploy:
	docker-compose -f build/compose/docker-compose.yml down

deploy:
	docker-compose -f build/compose/docker-compose.yml up -d --build

status:
	docker-compose -f build/compose/docker-compose.yml ps

run:
	docker-compose -f build/compose/docker-compose.yml up

.PHONY: tests
tests:
	docker-compose -f build/compose/docker-compose.itest.yml up  -d --build
	docker-compose -f build/compose/docker-compose.itest.yml logs --follow integration_test
.PHONY: untests
untests:
	docker-compose -f build/compose/docker-compose.itest.yml down