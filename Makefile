LISTEN_PORT?=9999

server:
	cd src; LISTEN_PORT=${LISTEN_PORT} go run main.go server

.PHONY: wire
wire:
	cd src/cmd/server; wire
	cd src/server/tests; wire