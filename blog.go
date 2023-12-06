//极简博客

package main

import (
	"blog/core"
	"flag"
)

var (
	//可定制化参数
	configPath = flag.String("c", "./", "配置文件路径")
	port = flag.String("p", "8081", "运行端口")
	authInfo   = flag.String("a", "", "开启账号密码登录验证, -a user:password的格式传参")
)

func init() {
	//参数初始化
	flag.Parse()
}

func main() {
	//启动服务
	core.StartServer(*configPath,*port,*authInfo)
}
