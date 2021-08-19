package main

// TODO: support validate
var ginTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{$validate := .GenValidate}}

// 这里定义 handler interface
type {{.ServiceType}}HTTPHandler interface {
{{- range .MethodSets}}
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
            c.AbortWithStatusJSON(400, gin.H{"err": err.Error()})
            return
        }

        // execute
        out, err = srv.{{.Name}}(ctx, in)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"err": err.Error()})
            return
        }
        
        c.JSON(200, out)
    }
}
{{end}}
`
