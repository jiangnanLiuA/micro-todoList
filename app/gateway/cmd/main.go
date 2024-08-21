package main

import (
	"fmt"
	"micro_todoList/app/gateway/router"
	"micro_todoList/app/gateway/rpc"
	"micro_todoList/config"
	log "micro_todoList/pkg/logger"
	"time"

	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
)

func main() {
	config.Init()
	rpc.InitRPC()
	log.InitLog()
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.EtcdHost, config.EtcdPort)),
	)

	// 创建微服务实例，使用gin暴露http接口并注册到etcd
	server := web.NewService(
		web.Name("httpService"),
		web.Address("127.0.0.1:4000"),
		// 将服务调用实例使用gin处理
		web.Handler(router.NewRouter()),
		web.Registry(etcdReg),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
		web.Metadata(map[string]string{"protocol": "http"}),
	)
	// 接收命令行参数
	_ = server.Init()
	_ = server.Run()
}
