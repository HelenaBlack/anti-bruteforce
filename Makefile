BIN_SRV := "./bin/antibruteforce"
BIN_CLI := "./bin/admin"

.PHONY: build
build:
	go build -o $(BIN_SRV) ./cmd/antibruteforce
	go build -o $(BIN_CLI) ./cmd/admin

.PHONY: run
run:
	docker-compose up --build

.PHONY: stop
stop:
	docker-compose down

.PHONY: test
test:
	go test -v -race -count 100 ./internal/...

.PHONY: test-integration
test-integration:
	go test -v -count=1 ./tests/integration/...

.PHONY: generate
generate:
	protoc --proto_path=api/proto --go_out=api/gen --go_opt=paths=source_relative --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative antibruteforce.proto

.PHONY: lint
lint:
	golangci-lint run
