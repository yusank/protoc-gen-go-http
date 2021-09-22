// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v0.0.1

package v1

import (
	context "context"
	gin "github.com/gin-gonic/gin"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the galaxy package it is being compiled against.
var _ context.Context

const _ = gin.Version

// 这里定义 handler interface
type HelloHTTPHandler interface {
	Add(context.Context, *AddRequest) (*AddResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
}

// RegisterHelloHTTPHandler define http router handle by gin.
// 注册路由 handler
func RegisterHelloHTTPHandler(g *gin.RouterGroup, srv HelloHTTPHandler) {
	g.POST("/api/hello/service/v1/add", _Hello_Add0_HTTP_Handler(srv))
	g.GET("/api/hello/service/v1/get", _Hello_Get0_HTTP_Handler(srv))
}

// 定义 handler
// 遍历之前解析到所有 rpc 方法信息

func _Hello_Add0_HTTP_Handler(srv HelloHTTPHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
			in  = new(AddRequest)
			out = new(AddResponse)
			ctx = context.TODO()
		)

		if err = c.ShouldBind(in); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"err": err.Error()})
			return
		}

		// execute
		out, err = srv.Add(ctx, in)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
			return
		}

		c.JSON(200, out)
	}
}

func _Hello_Get0_HTTP_Handler(srv HelloHTTPHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
			in  = new(GetRequest)
			out = new(GetResponse)
			ctx = context.TODO()
		)

		if err = c.ShouldBind(in); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"err": err.Error()})
			return
		}

		// execute
		out, err = srv.Get(ctx, in)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
			return
		}

		c.JSON(200, out)
	}
}

// Client defines call remote server client and implement selector
type Client interface {
	Call(ctx context.Context, req, rsp interface{}) error
}

// HelloHTTPClient defines call HelloServer client
type HelloHTTPClient interface {
	Add(context.Context, *AddRequest) (*AddResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
}

// HelloHTTPClientImpl implement HelloHTTPClient
type HelloHTTPClientImpl struct {
	cli Client
}

func NewHelloHTTPClient(cli Client) HelloHTTPClient {
	return &HelloHTTPClientImpl{
		cli: cli,
	}
}

func (c *HelloHTTPClientImpl) Add(ctx context.Context, req *AddRequest) (resp *AddResponse, err error) {
	resp = new(AddResponse)
	err = c.cli.Call(ctx, req, resp)

	return
}

func (c *HelloHTTPClientImpl) Get(ctx context.Context, req *GetRequest) (resp *GetResponse, err error) {
	resp = new(GetResponse)
	err = c.cli.Call(ctx, req, resp)

	return
}
