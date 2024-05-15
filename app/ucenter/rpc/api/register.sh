#goctl rpc protoc register.proto --go_out=./types --go-grpc_out=./types --zrpc_out=./register --style go_zero


# 单独使用protoc，用于在protoc中新增字段
protoc register.proto --go_out=./types --go-grpc_out=./types