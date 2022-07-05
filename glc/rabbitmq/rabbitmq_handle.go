package rabbitmq

import (
	"glc/conf"
	"glc/rabbitmq/consume"
	"log"
	"time"
)

func Start() {
	go func() {
		if conf.IsEnableAmqpConsume() {
			err := consume.StartRabbitMQConsume()
			if err != nil {
				log.Println("RabbitMQ连接失败（10秒后再试）", err)
				timer := time.NewTimer(time.Second * 10)
				<-timer.C
				Start() // 10秒后再试
			} else {
				log.Println("启动RabbitMQ日志消费")
			}
		}
	}()
}

func Stop() {
	if conf.IsEnableAmqpConsume() {
		consume.StopRabbitMQConsume()
		log.Println("停止RabbitMQ日志消费")
	}
}
