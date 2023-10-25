LISTEN_PORT?=9999

server:
	cd src; LISTEN_PORT=${LISTEN_PORT} go run main.go server

.PHONY: wire
wire: mockery
	cd src/cmd/server; wire
	cd src/server/tests; wire

.PHONY: protoc
protoc:
	cd src/proto; protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative message.proto

.PHONY: mockery
mockery:
	cd src; mockery

.PHONY: test
test:
	cd src/configs; go test ./...
	cd src/server/tests; go test ./...
	cd src/notifier/kafka; go test -race ./...