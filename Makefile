.PHONY:

build:
	go build -o ./.bin/bot/main.go cmd/bot/main.go

CFLAGS = -token wrong-test-token

run: build
	./.bin/bot/main.go $(CFLAGS)