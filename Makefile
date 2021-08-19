
proto: ## protoc
	protoc -I/usr/local/include -I$(GOPATH)/src/github.com/googleapis/googleapis\
 	--proto_path=$(GOPATH)/src:. --go_out=. --go-http_out=. --go-grpc_out=.\
 	 example/greeter/v1/hello.proto
