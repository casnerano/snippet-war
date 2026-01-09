LOCAL_BIN := $(CURDIR)/bin
VENDOR_PROTO_PATH := $(CURDIR)/vendor.protogen

.PHONY: download-bin-deps
download-bin-deps:
	ls $(LOCAL_BIN)/buf &> /dev/null || GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@latest
	ls $(LOCAL_BIN)/protoc-gen-go &> /dev/null || GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	ls $(LOCAL_BIN)/protoc-gen-go-grpc &> /dev/null || GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	ls $(LOCAL_BIN)/protoc-gen-grpc-gateway &> /dev/null || GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.18.1
	ls $(LOCAL_BIN)/protoc-gen-openapiv2 &> /dev/null || GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.18.1

.PHONY: vendor-proto/google/api
vendor-proto/google/api:
	rm -rf $(VENDOR_PROTO_PATH)/google
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis $(VENDOR_PROTO_PATH)/temp && \
 		cd $(VENDOR_PROTO_PATH)/temp && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
	mkdir -p  $(VENDOR_PROTO_PATH)/google
	mv $(VENDOR_PROTO_PATH)/temp/google/api $(VENDOR_PROTO_PATH)/google
	rm -rf $(VENDOR_PROTO_PATH)/temp

.PHONY: vendor-proto/protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	rm -rf $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/grpc-ecosystem/grpc-gateway $(VENDOR_PROTO_PATH)/temp && \
 		cd $(VENDOR_PROTO_PATH)/temp && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	mv $(VENDOR_PROTO_PATH)/temp/protoc-gen-openapiv2/options $(VENDOR_PROTO_PATH)/protoc-gen-openapiv2
	rm -rf $(VENDOR_PROTO_PATH)/temp

.PHONY: vendor-proto/validate
vendor-proto/validate:
	rm -rf $(VENDOR_PROTO_PATH)/validate
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate $(VENDOR_PROTO_PATH)/temp && \
		cd $(VENDOR_PROTO_PATH)/temp && \
		git sparse-checkout set --no-cone validate && \
		git checkout
	mkdir -p $(VENDOR_PROTO_PATH)/validate
	mv $(VENDOR_PROTO_PATH)/temp/validate $(VENDOR_PROTO_PATH)/
	rm -rf $(VENDOR_PROTO_PATH)/temp

.PHONY: vendor-proto
vendor-proto: vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options vendor-proto/validate

.PHONY: generate-proto
generate-proto:
	mkdir -p api/openapi
	mkdir -p pkg
	protoc \
	  --proto_path=. \
	  --proto_path=vendor.protogen \
	  \
	  --plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	  --go_out=./pkg \
	  --go_opt=paths=source_relative \
	  \
	  --plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	  --go-grpc_out=./pkg \
	  --go-grpc_opt=paths=source_relative \
	  \
	  --plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	  --grpc-gateway_out=./pkg \
	  --grpc-gateway_opt=logtostderr=true \
	  --grpc-gateway_opt=paths=source_relative \
	  \
	  --plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	  --openapiv2_out=./api/openapi \
	  --openapiv2_opt=logtostderr=true \
	  --openapiv2_opt=allow_merge=true \
	  --openapiv2_opt=merge_file_name=quiz \
	  \
	  ./api/v1/quiz/service.proto

.PHONY: generate
generate: download-bin-deps generate-proto
	go mod tidy

.PHONY: build
build:
	go build -o ${LOCAL_BIN}/snippet-war ./cmd/snippet-war

.PHONY: run
run:
	go run ./cmd/snippet-war/main.go