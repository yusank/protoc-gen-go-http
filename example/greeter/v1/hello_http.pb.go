// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v0.0.2

package v1

import (
	context "context"
	errors "errors"
	gin "github.com/gin-gonic/gin"
	http "net/http"
)

// This imports are custom by go-http.
import (
	restyv2 "github.com/go-resty/resty/v2"
	phttp "github.com/yusank/protoc-gen-go-http/http"
)

// This is a compile-time assertion to ensure that generated files are safe and compilable.
var _ context.Context
var _ http.Client
var _ phttp.CallOption
var _ = errors.New

const _ = gin.Version
const _ = restyv2.Version

// HelloHTTPHandler defines HelloServer http handler
type HelloHTTPHandler interface {
	Add(context.Context, *AddRequest) (*AddResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
}

// RegisterHelloHTTPHandler define http router handle by gin.
func RegisterHelloHTTPHandler(g *gin.RouterGroup, srv HelloHTTPHandler) {
	g.POST("/api/hello/service/v1/add", _Hello_Add0_HTTP_Handler(srv))
	g.GET("/api/hello/service/v1/get", _Hello_Get0_HTTP_Handler(srv))
}

type Validator interface {
	Validate() error
}

// _Hello_Add0_HTTP_Handler is gin http handler to handle
// http request [POST] /api/hello/service/v1/add.
func _Hello_Add0_HTTP_Handler(srv HelloHTTPHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
			in  = new(AddRequest)
			out = new(AddResponse)
			ctx = context.TODO()
		)

		if err = c.ShouldBind(in); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		v, ok := interface{}(in).(Validator)
		if ok {
			if err = v.Validate(); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
				return
			}
		}

		// execute
		out, err = srv.Add(ctx, in)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, out)
	}
}

// _Hello_Get0_HTTP_Handler is gin http handler to handle
// http request [GET] /api/hello/service/v1/get.
func _Hello_Get0_HTTP_Handler(srv HelloHTTPHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
			in  = new(GetRequest)
			out = new(GetResponse)
			ctx = context.TODO()
		)

		if err = c.ShouldBind(in); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		v, ok := interface{}(in).(Validator)
		if ok {
			if err = v.Validate(); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
				return
			}
		}

		// execute
		out, err = srv.Get(ctx, in)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, out)
	}
}

// HelloHTTPClient defines call HelloServer client
type HelloHTTPClient interface {
	Add(ctx context.Context, req *AddRequest, opts ...phttp.CallOption) (*AddResponse, error)
	Get(ctx context.Context, req *GetRequest, opts ...phttp.CallOption) (*GetResponse, error)
}

// HelloHTTPClientImpl implement HelloHTTPClient
type HelloHTTPClientImpl struct {
	cli        *restyv2.Client
	clientOpts []phttp.ClientOption
}

func NewHelloHTTPClient(baseUrl string, cli *http.Client, opts ...phttp.ClientOption) (HelloHTTPClient, error) {
	if baseUrl == "" {
		return nil, errors.New("base url is empty")
	}

	c := &HelloHTTPClientImpl{
		clientOpts: opts,
	}

	hc := cli
	if hc == nil {
		hc = http.DefaultClient
	}

	c.cli = restyv2.NewWithClient(hc)
	c.cli.SetBaseURL(baseUrl)
	for _, opt := range opts {
		if err := opt.Apply(c.cli); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Add is call [POST] /api/hello/service/v1/add api.
func (c *HelloHTTPClientImpl) Add(ctx context.Context, req *AddRequest, opts ...phttp.CallOption) (rsp *AddResponse, err error) {
	rsp = new(AddResponse)

	r := c.cli.R()
	for _, opt := range opts {
		if err = opt.Before(r); err != nil {
			return
		}
	}
	// set response data struct.
	r.SetResult(rsp)
	// do request
	restyResp, err := r.Execute("POST", "/api/hello/service/v1/add")
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err = opt.After(restyResp); err != nil {
			return
		}
	}

	return
}

// Get is call [GET] /api/hello/service/v1/get api.
func (c *HelloHTTPClientImpl) Get(ctx context.Context, req *GetRequest, opts ...phttp.CallOption) (rsp *GetResponse, err error) {
	rsp = new(GetResponse)

	r := c.cli.R()
	for _, opt := range opts {
		if err = opt.Before(r); err != nil {
			return
		}
	}
	// set response data struct.
	r.SetResult(rsp)
	// do request
	restyResp, err := r.Execute("GET", "/api/hello/service/v1/get")
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err = opt.After(restyResp); err != nil {
			return
		}
	}

	return
}