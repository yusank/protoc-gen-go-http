package greeter

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	restyv2 "github.com/go-resty/resty/v2"
	"google.golang.org/grpc"

	v1 "github.com/yusank/protoc-gen-go-http/example/greeter/v1"
	phttp "github.com/yusank/protoc-gen-go-http/http"
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

func callAsClient() {
	opt := func(c *restyv2.Client) error {
		c.SetRetryCount(3)
		c.SetRetryWaitTime(time.Millisecond * 100)
		return nil
	}
	cli, err := v1.NewHelloHTTPClient("http://127.0.0.1:8080", nil, phttp.ApplyToClient(opt))
	if err != nil {
		return
	}

	reqOpt := func(r *restyv2.Request) error {
		r.SetHeader("custom-header", "proto-gen-go-http")
		return nil
	}
	respOpt := func(r *restyv2.Response) error {
		v := r.Header().Get("custom-header")
		if v != "proto-gen-go-http" {
			return errors.New("invalid custom-header")
		}

		return nil
	}
	rsp, err := cli.Add(context.Background(), &v1.AddRequest{}, phttp.Before(reqOpt), phttp.After(respOpt))
	if err != nil {
		return
	}

	// handle response
	_ = rsp.String()
}
