package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// Version protoc-gen-go-http 工具版本
const Version = "v0.0.2"

func main() {
	// 1. 传参定义
	// 即 插件是支持自定义参数的，这样我们可以更加灵活，针对不同的场景生成不同的代码
	var flags flag.FlagSet
	// 是否忽略没有指定 google.api 的方法
	omitempty := flags.Bool("omitempty", true, "omit if google.api is empty")
	// 我这里同时支持了 gin 和 iris 可以通过参数指定生成
	// 是否生校验代码块
	// 发现了一个很有用的插件 github.com/envoyproxy/protoc-gen-validate
	// 可以在 pb 的 message 中设置参数规则，然后会生成一个 validate.go 的文件 针对每个 message 生成一个 Validate() 方法
	// 我在每个 handler 处理业务前做了一次参数校验判断，通过这个 flag 控制是否生成这段校验代码
	genValidateCode := flags.Bool("validate", false, "add validate request params in handler")
	// 生成代码时参数 这么传：--go-http_out=router=iris,validate=true:.

	gp := &GenParam{
		Omitempty:       omitempty,
		GenValidateCode: genValidateCode,
	}
	// 这里就是入口，指定 option 后执行 Run 方法 ，我们的主逻辑就是在 Run 方法
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			// 这里是我们的生成代码方法
			generateFile(gen, f, gp)
		}
		return nil
	})
}

type GenParam struct {
	Omitempty       *bool
	GenValidateCode *bool
}
