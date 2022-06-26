package ldb

import (
	"log"
	"testing"
	"time"
)

func Test_all(t *testing.T) {

	// go func() {

	engine := NewDefaultEngine()
	for i := 1; i <= 100; i++ {
		engine.AddTextLog(`   java.sql.SQLException:  ddduse them aalav`)
	}
	time.Sleep(time.Duration(5) * time.Second)

	for i := 1; i <= 10000; i++ {
		engine.AddTextLog(`   java.sql.SQLException:   them aalav`)
	}
	time.Sleep(time.Duration(5) * time.Second)

	for i := 1; i <= 10000; i++ {
		engine.AddTextLog(`   java.sql.SQLException: them`)
	}
	time.Sleep(time.Duration(5) * time.Second)

	rs := engine.Search(`     them java      `, 20, 0, true)
	log.Println("共查到", rs.TotalCount, "件")
	for _, v := range rs.Data {
		log.Println(v.Content)
	}
	//time.Sleep(time.Duration(5) * time.Second)
	// }()
	// time.Sleep(time.Duration(30) * time.Second)

}
