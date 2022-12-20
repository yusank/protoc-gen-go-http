package main

// TODO: support validate
var httpCodeTmpl = `
{{/*gotype: github.com/yusank/protoc-gen-go-http.serviceDesc*/}}
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{$validate := .GenValidate}}

// 这里定义 handler interface
type {{.ServiceType}}HTTPHandler interface {
{{- range .Methods}}
    {{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

// Register{{.ServiceType}}HTTPHandler define http router handle by gin.
// 注册路由 handler
func Register{{.ServiceType}}HTTPHandler(g *gin.RouterGroup, srv {{.ServiceType}}HTTPHandler) {
{{- range .Methods}}
    g.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
{{- end}}
}

{{if $validate}}
type Validator interface {
    Validate() error
}
{{end}}

// 定义 handler
// 遍历之前解析到所有 rpc 方法信息
{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPHandler) func(c *gin.Context) {
    return func(c *gin.Context) {
        var (
            err error
            in  = new({{.Request}})
            out = new({{.Reply}})
            ctx = context.TODO()
        )

        if err = c.ShouldBind(in{{.Body}}); err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
            return
        }
        {{if $validate}}
        v,ok := interface{}(in).(Validator)
        if ok {
            if err = v.Validate();err != nil {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
                return
            }
        }
        {{end}}
        // execute
        out, err = srv.{{.Name}}(ctx, in)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
            return
        }

        c.JSON(http.StatusOK, out)
    }
}
{{end}}
// {{.ServiceType}}HTTPClient defines call {{.ServiceType}}Server client
type {{.ServiceType}}HTTPClient interface {
{{- range .Methods}}
    {{.Name}}(ctx context.Context, req *{{.Request}}, opts ...phttp.CallOption) (*{{.Reply}}, error)
{{- end}}
}

// {{.ServiceType}}HTTPClientImpl implement {{.ServiceType}}HTTPClient
type {{.ServiceType}}HTTPClientImpl struct {
    cli *restyv2.Client
    clientOpts []phttp.ClientOption
}

func New{{.ServiceType}}HTTPClient(cli *http.Client, opts ...phttp.ClientOption) ({{.ServiceType}}HTTPClient, error) {
    c := &{{.ServiceType}}HTTPClientImpl{
        clientOpts: opts,
    }

    hc := cli
    if hc == nil {
        hc = http.DefaultClient
    }

    c.cli = restyv2.NewWithClient(hc)
    for _, opt := range opts {
        if err := opt.Apply(c.cli);err != nil {
            return nil, err
        }
    }

    return c, nil
}

{{range .Methods}}
// {{.Name}} is call [{{.Method}}] {{.Path}} api.
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, req *{{.Request}}, opts ...phttp.CallOption)(rsp *{{.Reply}} ,err error) {
    rsp = new({{.Reply}})

    r := c.cli.R()
    for _, opt := range opts {
        if err = opt.Before(r);err != nil {
            return
        }
    }
    // set response data struct.
    r.SetResult(rsp)
    // do request
    restyResp,err := r.Execute("{{.Method}}", "{{.Path}}")
    if err != nil {
        return nil, err
    }
    for _, opt := range opts {
        if err = opt.After(restyResp);err != nil {
            return
        }
    }

    return
}
{{end}}
`
