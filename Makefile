LOCAL_BIN := $(CURDIR)/bin

.PHONY: build
build:
	go build -o ${LOCAL_BIN}/snippet-war ./cmd/snippet-war

.PHONY: run
run:
	go run ./cmd/snippet-war/main.go

