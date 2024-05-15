# 单独使用protoc，用于在protoc中新增字段
protoc login.proto --go_out=./types --go-grpc_out=./types