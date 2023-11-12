/**
 * RabbitMQ日志接收封装
 */
package consume

import (
	"glc/conf"
	"glc/ldb"
	"glc/ldb/storage/logdata"
	"glc/www/controller"
	"sync"

	"github.com/gotoeasy/glang/cmn"
)

type RabbitMQConsume struct {
	rabbitMQ *RabbitMQ
	mu       sync.Mutex // 锁
}

var rabbitMQConsume *RabbitMQConsume

func init() {
	rabbitMQConsume = new(RabbitMQConsume)
	cmn.OnExit(onExit)
}

func StartRabbitMQConsume() error {

	if rabbitMQConsume.rabbitMQ != nil && !rabbitMQConsume.rabbitMQ.closing {
		return nil
	}

	rabbitMQConsume.mu.Lock()
	defer rabbitMQConsume.mu.Unlock()

	if rabbitMQConsume.rabbitMQ != nil && !rabbitMQConsume.rabbitMQ.closing {
		return nil
	}

	mq, err := NewSimpleRabbitMQ()
	if err != nil {
		return err
	} else {
		rabbitMQConsume.rabbitMQ = mq
		rabbitMQConsume.rabbitMQ.StartConsume(fnAmqpJsonLogHandle)
	}
	return nil
}

func StopRabbitMQConsume() {
	if rabbitMQConsume.rabbitMQ != nil {
		rabbitMQConsume.rabbitMQ.Close()
	}
}

func fnAmqpJsonLogHandle(jsonLog string, err error) bool {
	if err != nil {
		cmn.Error(err)
		return false
	}

	// cmn.Debug("接收到rabbitmq的日志", jsonLog)

	md := &logdata.LogDataModel{}
	if conf.IsAmqpJsonFormat() {
		if md.LoadJson(jsonLog) != nil {
			md.Text = jsonLog // 错误的json字符串？当文本吃掉
		}
	} else {
		md.Text = jsonLog
	}

	engine := ldb.NewDefaultEngine()
	md.Text = cmn.Trim(md.Text)
	engine.AddLogDataModel(md)

	if conf.IsClusterMode() {
		go controller.TransferGlc(conf.LogTransferAdd, md.ToJson()) // 转发其他GLC服务
	}

	return true
}

func onExit() {
	rabbitMQConsume.mu.Lock()
	defer rabbitMQConsume.mu.Unlock()
	rabbitMQConsume.rabbitMQ.Close()
}
