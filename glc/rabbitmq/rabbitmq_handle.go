/**
 * RabbitMQ日志接收的启停
 */
package rabbitmq

import (
	"glc/conf"
	"glc/rabbitmq/consume"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

func Start() {
	go func() {
		if conf.IsEnableAmqpConsume() {
			err := consume.StartRabbitMQConsume()
			if err != nil {
				cmn.Error("RabbitMQ连接失败（10秒后再试）", err)
				timer := time.NewTimer(time.Second * 10)
				<-timer.C
				Start() // 10秒后再试
			} else {
				cmn.Info("启动RabbitMQ日志消费")
			}
		}
	}()
}

func Stop() {
	if conf.IsEnableAmqpConsume() {
		consume.StopRabbitMQConsume()
		cmn.Info("停止RabbitMQ日志消费")
	}
}
