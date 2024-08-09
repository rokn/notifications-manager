# Notifications Manager

## Requirements

- Go `1.22`
- Docker


## Development
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/swaggo/swag/cmd/swag@latest

# Generate proto files
protoc --go_out=. --go-grpc_out=. api/channels/api.proto

# Generate docs
swag init --dir pkg/ingress/ -g server.go --output pkg/ingress/docs
```
