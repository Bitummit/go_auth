PHONY: generate-structs
generate-structs:
	mkdir -p pkg/auth_proto_gen
	protoc --go_out=pkg/auth_proto_gen --go_opt=paths=source_relative \
		proto/auth.proto

PHONY: generate
generate:
	mkdir -p pkg/auth_proto_gen
	protoc --go_out=pkg/auth_proto_gen --go_opt=paths=source_relative \
	--go-grpc_out=pkg/auth_proto_gen --go-grpc_opt=paths=source_relative \
	./**/*.proto