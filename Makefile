.PHONY: run-client
run-client:
	go run client/*.go 

.PHONY: run-server
run-server:
	go run server/*.go 

.PHONY: proto-gen
proto-gen:
	protoc --go_out=./proto --go-grpc_out=./proto proto/*.proto 
