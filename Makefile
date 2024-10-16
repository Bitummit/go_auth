PHONY: generate-structs
generate-structs:
	mkdir -p pkg/auth_v1
	protoc --go_out=pkg/auth_v1 --go_opt=paths=source_relative \
		proto/auth.proto

PHONY: generate
generate:
	mkdir -p pkg/auth_proto_gen
	protoc --go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	./**/*.proto