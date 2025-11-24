## 开发事项
1.本项目使用了 wire 依赖注入框架，请先安装wire，并添加 gopath/bin到环境变量    
`go install github.com/google/wire/cmd/wire@latest`   
2.开发中有添加或修改了组件，请在项目目录中使用wire命令来生成新的wire_gen.go文件，如下：      
```
wire
```   
3.生成swagger文档   
`swag init --parseDependency --parseInternal -g cmd/main.go`

