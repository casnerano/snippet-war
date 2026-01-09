LOCAL_BIN := $(CURDIR)/bin

.PHONY: download-bin-deps
download-bin-deps:
	ls $(LOCAL_BIN)/buf &> /dev/null || GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@latest
	ls $(LOCAL_BIN)/protoc-gen-go &> /dev/null || GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	ls $(LOCAL_BIN)/protoc-gen-go-grpc &> /dev/null || GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

.PHONY: generate-rbac-proto
generate-rbac-proto:
	$(LOCAL_BIN)/buf generate --config ${CURDIR}/buf.yaml --template ${CURDIR}/buf.gen.yaml

.PHONY: generate
generate: download-bin-deps generate-rbac-proto
	go mod tidy

.PHONY: build
build:
	go build -o ${LOCAL_BIN}/snippet-war ./cmd/snippet-war

.PHONY: run
run:
	go run ./cmd/snippet-war/main.go
