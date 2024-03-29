package main

import (
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	errorsPackage  = protogen.GoImportPath("errors")
	httpPackage    = protogen.GoImportPath("net/http")
	contextPackage = protogen.GoImportPath("context")
	ginPackage     = protogen.GoImportPath("github.com/gin-gonic/gin")
	phttpPackage   = protogen.GoImportPath("github.com/yusank/protoc-gen-go-http/http")
	restyv2Package = protogen.GoImportPath("github.com/go-resty/resty/v2")

	deprecationComment = "// Deprecated: Do Not Use."
)

var methodSets = make(map[string]int)

// generateFile generates a _http.pb.go file containing gin/iris handler.
func generateFile(gen *protogen.Plugin, file *protogen.File, gp *GenParam) *protogen.GeneratedFile {
	if len(file.Services) == 0 || (*gp.Omitempty && !hasHTTPRule(file.Services)) {
		return nil
	}
	// 这里我们可以自定义文件名
	filename := file.GeneratedFilenamePrefix + "_http.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	// 写入一些警告之类的 告诉用户不要修改
	g.P("// Code generated by protoc-gen-go-http. DO NOT EDIT.")
	g.P("// versions:")
	g.P(fmt.Sprintf("// protoc-gen-go-http %s", Version))
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	generateFileContent(gen, file, g, gp)
	return g
}

// generateFileContent generates the _http.pb.go file content, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, gp *GenParam) {
	if len(file.Services) == 0 {
		return
	}
	//// import
	//// 这里有个插曲：其实 import 相关的代码我们这么不需要特殊指定，protogen 包会帮我们处理，
	//// 但是import 的 path 前的别名默认取 path 最后一个 `/` 之后的字符，
	//// 比如：github.com/kataras/iris/v12 被处理成 v12 "github.com/kataras/iris/v12"
	//// 这个我不太愿意接受 所以自己写入 import
	g.P("// This imports are custom by go-http.")
	g.P("import (")
	g.P("restyv2", " ", restyv2Package)
	g.P("phttp", " ", phttpPackage)
	g.P(")")

	// 注： 我们难免有一些 _ "my/package" 这种需求，这其实不用自己写 直接调 g.Import("my/package") 就可以

	// 这里定义一堆变量是为了程序编译的时候确保这些包是正确的，如果包不存在或者这些定义的包变量不存在都会编译失败
	g.P("// This is a compile-time assertion to ensure that generated files are safe and compilable.")
	// 只要调用这个 Ident 方法 就会自动写入到 import 中 ，所以如果对 import 的包名没有特殊要求，那就直接使用 Ident
	g.P("var _ ", contextPackage.Ident("Context"))
	g.P("var _ ", httpPackage.Ident("Client"))
	g.P("var _ ", "phttp.CallOption")
	g.P("var _ =", errorsPackage.Ident("New"))
	g.P("const _ = ", ginPackage.Ident("Version"))
	g.P("const _ = ", "restyv2.Version")
	g.P()

	// 到这里我们就把包名 import 和变量写入成功了，剩下的就是针对 rpc service 生成对应的 handler
	for _, service := range file.Services {
		genService(gen, file, g, service, gp)
	}
}

// 生成 service 相关代码
func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, gp *GenParam) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}

	// HTTP Server.
	// 服务的主要变量，比如服务名 服务类型等
	sd := &serviceDesc{
		ServiceType: service.GoName,
		ServiceName: string(service.Desc.FullName()),
		Metadata:    file.Desc.Path(),
		GenValidate: *gp.GenValidateCode,
	}
	// 开始遍历服务的方法
	for _, method := range service.Methods {
		// 不处理
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue
		}
		// annotations 这个就是我们在 rpc 方法里 option 里定义的 http 路由
		rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule != nil && ok {
			for _, bind := range rule.AdditionalBindings {
				// 拿到 option里定义的路由， http method等信息
				sd.Methods = append(sd.Methods, buildHTTPRule(g, method, bind))
			}
			sd.Methods = append(sd.Methods, buildHTTPRule(g, method, rule))
		} else if !*gp.Omitempty {
			path := fmt.Sprintf("/%s/%s", service.Desc.FullName(), method.Desc.Name())
			sd.Methods = append(sd.Methods, buildMethodDesc(g, method, "POST", path))
		}
	}

	// 拿到了 n 个 rpc 方法，开始生成了
	if len(sd.Methods) != 0 {
		// 渲染
		g.P(sd.execute())
	}
}

// 检查是否有 http 规则 即
//
//	option (google.api.http) = {
//	     get: "/user/query"
//	   };
func hasHTTPRule(services []*protogen.Service) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				continue
			}
			rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
			if rule != nil && ok {
				return true
			}
		}
	}
	return false
}

// 解析 http 规则，读取内容
func buildHTTPRule(g *protogen.GeneratedFile, m *protogen.Method, rule *annotations.HttpRule) *methodDesc {
	var (
		path         string
		method       string
		body         string
		responseBody string
	)
	// 读取 路由和方法
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = "GET"
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = "PUT"
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = "POST"
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = "DELETE"
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = "PATCH"
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}
	body = rule.Body
	responseBody = rule.ResponseBody
	md := buildMethodDesc(g, m, method, path)
	if method == "GET" {
		md.HasBody = false
	} else if body == "*" {
		md.HasBody = true
		md.Body = ""
	} else if body != "" {
		md.HasBody = true
		md.Body = "." + camelCaseVars(body)
	} else {
		md.HasBody = false
	}
	if responseBody == "*" {
		md.ResponseBody = ""
	} else if responseBody != "" {
		md.ResponseBody = "." + camelCaseVars(responseBody)
	}
	return md
}

// 构建 每个方法的基础信息
// 到这里我们拿到了 我们需要生成一个 handler 的所有信息
// 名称，输入，输出，方法类型，路由
func buildMethodDesc(g *protogen.GeneratedFile, m *protogen.Method, method, path string) *methodDesc {
	defer func() { methodSets[m.GoName]++ }()
	return &methodDesc{
		Name:    m.GoName,
		Num:     methodSets[m.GoName],
		Request: g.QualifiedGoIdent(m.Input.GoIdent),  // rpc 方法中的 request
		Reply:   g.QualifiedGoIdent(m.Output.GoIdent), // rpc 方法中的 response
		Path:    path,
		Method:  method,
		HasVars: len(buildPathVars(m, path)) > 0,
	}
}

// 处理 路由中 /api/user/{name} 这种情况
func buildPathVars(method *protogen.Method, path string) (res []string) {
	for _, v := range strings.Split(path, "/") {
		if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			name := strings.TrimRight(strings.TrimLeft(v, "{"), "}")
			res = append(res, name)
		}
	}
	return
}

func camelCaseVars(s string) string {
	var (
		vars []string
		subs = strings.Split(s, ".")
	)
	for _, sub := range subs {
		vars = append(vars, camelCase(sub))
	}
	return strings.Join(vars, ".")
}

// camelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
