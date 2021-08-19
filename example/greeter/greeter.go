package greeter

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	v1 "github.com/yusank/protoc-gen-go-http/example/greeter/v1"
)

// Greeter implement v1.HelloServer and v1.HelloHTTPHandler at same time
type Greeter struct {
	v1.UnimplementedHelloServer
}

func (g *Greeter) Add(ctx context.Context, in *v1.AddRequest) (*v1.AddResponse, error) {
	panic("implement me")
}

func (g *Greeter) Get(ctx context.Context, in *v1.GetRequest) (*v1.GetResponse, error) {
	panic("implement me")
}

func register() {
	// can call method by rpc or http
	v1.RegisterHelloServer(&grpc.Server{}, &Greeter{})
	v1.RegisterHelloHTTPHandler(gin.Default().Group("/"), &Greeter{})
}
