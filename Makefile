LISTEN_PORT?=9999

server:
	cd src; LISTEN_PORT=${LISTEN_PORT} go run main.go server

.PHONY: wire
wire: mockery
	cd src/cmd/server; wire
	cd src/tests; wire

.PHONY: protoc
protoc:
	cd src/proto; protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative message.proto

.PHONY: mockery
mockery:
	cd src; mockery

.PHONY: test
test:
	go clean -testcache
	cd src; go test -race ./...

.PHONY: lint
lint:
	cd src; golangci-lint run ./...

.PHONY: benthos
benthos:
	REDIS_URL=redis://localhost:6379 KAFKA_BROKERS=localhost:9094 benthos test .benthos/cache_users_benthos_test.yaml
	REDIS_URL=redis://localhost:6379 KAFKA_BROKERS=localhost:9094 benthos test .benthos/hydrate_messages_benthos_test.yaml 