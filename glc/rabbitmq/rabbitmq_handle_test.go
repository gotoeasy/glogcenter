package rabbitmq

import (
	"glc/rabbitmq/consume"
	"testing"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

func Test_all(t *testing.T) {

	r, err := consume.NewSimpleRabbitMQ()
	if err != nil {
		cmn.Error(err)
	}

	for i := 0; i < 1456; i++ {
		r.SimplePublish("{\"text\":\"aaa bbb ccc ddd  eee\",\"system\":\"rabbitmq\"}")
	}

	// Start()
	time.Sleep(time.Duration(10) * time.Second)
}
