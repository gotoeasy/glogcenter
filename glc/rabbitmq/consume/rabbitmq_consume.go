/**
 * RabbitMQ简单模式消费者封装
 */
package consume

import (
	"glc/conf"

	"github.com/gotoeasy/glang/cmn"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Path      string
	closing   bool
}

// 实例化(简单模式)
func NewSimpleRabbitMQ() (*RabbitMQ, error) {
	rabbitmq := &RabbitMQ{
		QueueName: conf.GetAmqpQueueName(),
		Path:      conf.GetAmqpAddr(),
	}

	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Path)
	if err != nil {
		return nil, err
	}
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		return nil, err
	}

	return rabbitmq, nil
}

// 关闭连接
func (r *RabbitMQ) Close() {
	if r == nil || r.closing {
		return
	}

	r.closing = true
	if r.channel != nil {
		r.channel.Close()
	}
	if r.channel != nil {
		r.conn.Close()
	}
}

// 简单模式生产者
func (r *RabbitMQ) SimplePublish(message string) {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // durable  是否持久化数据
		false, // autoDelete 是否自动删除
		false, // exclusive 排他性，权限私有
		false, // noWaite 是否阻塞
		nil,
	)

	if err != nil {
		cmn.Error(err)
	}
	r.channel.Publish(
		"",          // 交换机
		r.QueueName, // 队列名
		false,       // mandatory 如果是true，会根据exchange和routekey规则，如果无法找到符合条件的队列会把消息的消息返回给发送者
		false,       // immediate true 表示当exchange将消息发送到队列后发现队列没有绑定消费者，则会把消息发回给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// 简单模式消费者
func (r *RabbitMQ) StartConsume(fnJsonLogHandle func(string, error) bool) {
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名
		false,       // durable  是否持久化数据
		false,       // autoDelete 是否自动删除
		false,       // exclusive 排他性，权限私有
		false,       // noWaite 是否阻塞
		nil,         // arguments
	)

	if err != nil {
		fnJsonLogHandle("", err)
		return
	}

	mqDelivery, err := r.channel.Consume(
		r.QueueName, // 队列名
		"",          // 区分多个消费者
		true,        // autoAck 是否自动应答
		false,       // exclusive 是否排它性
		false,       // noLocal true表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,       // 不返回执行结果,但是如果排他开启的话,则必须需要等待结果的,如果两个一起开就会报错
		nil,         // 其他参数
	)

	if err != nil {
		fnJsonLogHandle("", err)
		return
	}

	go func() {
		for msg := range mqDelivery {
			// TODO 是否要考虑失败?
			fnJsonLogHandle(string(msg.Body), nil) // 处理接收到的日志
			if r.closing {
				break
			}
		}
	}()
}
