# protoc-gen-go-http
generate go http server code via proto buffer

## how to use

```shell
$ make proto
# or
$ protoc -I/usr/local/include -I$(GOPATH)/src/github.com/googleapis/googleapis\
 	--proto_path=$(GOPATH)/src:. --go_out=. --go-http_out=. --go-grpc_out=.\
 	 path/to/file.proto
```

## TODO

- [ ] support [validate](https://github.com/envoyproxy/protoc-gen-validate)
- [ ] gen client code.