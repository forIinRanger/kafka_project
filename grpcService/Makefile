.PHONY: proto
proto:
	protoc --go_out=./proto --go_opt=paths=source_relative \
	       --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
	       ./grpc.proto