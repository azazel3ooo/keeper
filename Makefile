.PHONY: build
build:
	go build ./cmd/client && go build ./cmd/server

.PHONY: client_windows
client_windows: # на windows проблемы с переменными окружения, поэтому просто выполняем стандартный билд
	GOOS=windows go build -o "client.exe" ./cmd/client

.PHONY: client_linux
client_linux:
	GOOS=linux go build -o "client" ./cmd/client

.PHONY: client_mac
client_mac:
	GOOS=darwin go build -o "client" ./cmd/client

.PHONY: server_linux
server_linux:
	GOOS=linux go build -o "server" ./cmd/server

.PHONY: server_windows
server_windows: # на windows проблемы с переменными окружения, поэтому просто выполняем стандартный билд
	GOOS=windows go build -o "server.exe" ./cmd/server

.PHONY: server_mac
server_mac:
	GOOS=darwin go build -o "server" ./cmd/server

.DEFAULT_GOAL := build