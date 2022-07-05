package rabbitmq

import (
	"glc/rabbitmq/consume"
	"log"
	"testing"
	"time"
)

func Test_all(t *testing.T) {

	r, err := consume.NewSimpleRabbitMQ()
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < 1456; i++ {
		r.SimplePublish("{\"text\":\"aaa bbb ccc ddd  eee\",\"system\":\"rabbitmq\"}")
	}

	// Start()
	time.Sleep(time.Duration(10) * time.Second)
}
