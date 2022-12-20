package main

import (
	"bytes"
	"strings"
	"text/template"
)

// rpc service 信息
type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/helloworld/helloworld.proto
	GenValidate bool
	Methods     []*methodDesc
}

// rpc 方法信息
type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

// execute 方法实现也其实不复杂，总起来就是 go 的 temple 包的使用
// 提前写好模板文件，然后拿到所有需要的变量，进行模板渲染，写入文件
func (s *serviceDesc) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http_server").Parse(strings.TrimSpace(httpCodeTmpl))
	if err != nil {
		panic(err)
	}
	if err = tmpl.Execute(buf, s); err != nil {
		panic(err)
	}

	return strings.Trim(buf.String(), "\r\n")
}
