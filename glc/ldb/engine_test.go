package ldb

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_all(t *testing.T) {

	engine := NewDefaultEngine()
	for i := 1; i <= 100; i++ {
		engine.AddTextLog("date", fmt.Sprintf(` java.sql.SQLException:  ddduse them aalav_%d`, i), "sssss_ssss")
	}
	time.Sleep(time.Duration(5) * time.Second)

	// for i := 1; i <= 10000; i++ {
	// 	engine.AddTextLog(`   java.sql.SQLException:   them aalav`)
	// }
	// time.Sleep(time.Duration(5) * time.Second)

	// for i := 1; i <= 10000; i++ {
	// 	engine.AddTextLog(`  java.sql.SQLException: them`)
	// }
	// time.Sleep(time.Duration(5) * time.Second)

	rs := engine.Search(`              them java    ddduse  `, 5, 0, true)
	log.Println("共查到", rs.Total, "件")
	for _, v := range rs.Data {
		log.Println(v.Id, v.Text)
	}

}
