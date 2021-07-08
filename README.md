### 用 Gin 搭建博客  

教程来源：
https://eddycjy.com/go-categories/

踩坑总结：

1、部署应用到 docker 时，不同架构需要进行跨平台编译
`RUN CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -o go-gin-example .
CGO_ENABLED -- CGO 工具是否可用，交叉编译时不可用
GOOS -- 目标环境的系统，windows linux darwin
GOARCH -- 目标环境的架构 386 amd64
-a 强制重新编译，简单来说，就是不利用缓存或已编译好的部分文件，直接所有包都是最新的代码重新编译和关联
-installsuffix 在软件包安装的目录中增加后缀标识，以保持输出与默认版本分开
-o 指定编译后的可执行文件名称
`