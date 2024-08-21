//package task
//
//import (
//	"context"
//	"encoding/json"
//	"micro_todoList/app/task/repository/mq"
//	"micro_todoList/app/task/service"
//	"micro_todoList/consts"
//	"micro_todoList/idl/pb"
//	log "micro_todoList/pkg/logger"
//	"sync"
//)
//
//type SyncTask struct {
//}
//
//// 监听队列
//
//func (s *SyncTask) RunTaskCreate(ctx context.Context) error {
//	rabbitMqQueue := consts.RabbitMqTaskQueue
//	msgs, err := mq.ConsumeMessage(ctx, rabbitMqQueue)
//	if err != nil {
//		return err
//	}
//	// forever 进行阻塞，防止进程被销毁
//	var forever chan struct{}
//
//	// goroutine
//	go func() {
//		for d := range msgs {
//			log.LogrusObj.Infof("Received run Task: %s", d.Body)
//
//			// 落库
//			reqRabbitMQ := new(pb.TaskRequest)
//			// 反序列化
//			err = json.Unmarshal(d.Body, reqRabbitMQ)
//			if err != nil {
//				log.LogrusObj.Infof("Received run Task: %s", err)
//			}
//
//			err = service.TaskMQ2MySQL(ctx, reqRabbitMQ)
//			if err != nil {
//				log.LogrusObj.Infof("Received run Task: %s", err)
//			}
//
//			d.Ack(false)
//
//		}
//	}()
//
//	log.LogrusObj.Infoln(err)
//	<-forever
//
//	return nil
//}

package task

import (
	"context"
	"encoding/json"
	"micro-todoList/app/task/repository/mq"
	"micro-todoList/app/task/service"
	"micro-todoList/consts"
	"micro-todoList/idl/pb"
	log "micro-todoList/pkg/logger"

	"sync"
)

type SyncTask struct{}

// 监听队列
func (s *SyncTask) RunTaskCreate(ctx context.Context) error {
	rabbitMqQueue := consts.RabbitMqTaskQueue
	msgs, err := mq.ConsumeMessage(ctx, rabbitMqQueue)
	if err != nil {
		return err
	}

	// 使用 sync.WaitGroup 来管理 goroutines
	var wg sync.WaitGroup

	// goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for d := range msgs {
			log.LogrusObj.Infof("Received run Task: %s", d.Body)

			// 反序列化
			reqRabbitMQ := new(pb.TaskRequest)
			err := json.Unmarshal(d.Body, reqRabbitMQ)
			if err != nil {
				log.LogrusObj.Errorf("Failed to unmarshal task: %v", err)
				continue // 跳过错误的消息
			}

			// 落库
			err = service.TaskMQ2MySQL(ctx, reqRabbitMQ)
			if err != nil {
				log.LogrusObj.Errorf("Failed to save task to MySQL: %v", err)
				continue // 如果落库失败，也跳过这条消息
			}

			// 消息确认
			if err := d.Ack(false); err != nil {
				log.LogrusObj.Errorf("Failed to ack message: %v", err)
			}
		}
	}()

	// 等待 goroutine 结束
	wg.Wait()
	return nil
}
