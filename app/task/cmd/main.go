package main

import (
	"context"
	"fmt"
	"micro_todoList/app/task/repository/mq"
	"micro_todoList/app/task/script"
	"micro_todoList/app/task/service"
	"micro_todoList/app/user/repository/dao"
	"micro_todoList/config"
	"micro_todoList/idl/pb"
	log "micro_todoList/pkg/logger"

	"go-micro.dev/v4"

	"go-micro.dev/v4/registry"
)

func main() {
	config.Init()
	dao.InitDB()
	mq.InitRabbitMQ()
	log.InitLog()

	// 启动一些脚本
	loadingScript()

	// etcd注册件
	etcdReg := registry.NewRegistry(
		registry.Addrs(fmt.Sprintf("%s:%s", config.EtcdHost, config.EtcdPort)),
	)
	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcTaskService"), // 微服务名字
		micro.Address(config.TaskServiceAddress),
		micro.Registry(etcdReg), // etcd注册件
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterTaskServiceHandler(microService.Server(), service.GetTaskSrv())
	// 启动微服务
	_ = microService.Run()
}

func loadingScript() {
	ctx := context.Background()
	go script.TaskCreateSync(ctx)
}
