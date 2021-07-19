protoc \
    --proto_path pkg/protos \
    --go_out=pkg/test \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/test \
    --go-grpc_opt=paths=source_relative \
    pkg/**/*.proto