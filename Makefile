.PHONY: proto
proto: install ## protoc
	protoc -I/usr/local/include -I$(GOPATH)/src/github.com/googleapis/googleapis\
 	--proto_path=$(GOPATH)/src:. --go_out=. --go-http_out=validate=true:. --go-grpc_out=.\
 	--validate_out=lang=go,paths=source_relative:.\
 	 example/greeter/v1/hello.proto

install:
	go install .