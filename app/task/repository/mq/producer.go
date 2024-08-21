package mq

import (
	"fmt"
	"micro-todoList/consts"

	"github.com/streadway/amqp"
)

// SendMessage2MQ 发送消息到  MQ
func SendMessage2MQ(body []byte) (err error) {
	ch, err := RabbitMq.Channel()
	if err != nil {
		return
	}

	// 消息队列
	q, _ := ch.QueueDeclare(consts.RabbitMqTaskQueue, true, false, false, false, nil)
	// 发布
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return
	}

	fmt.Println("发送MQ成功...")
	return
}
