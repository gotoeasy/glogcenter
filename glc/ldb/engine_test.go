package ldb

import (
	"log"
	"testing"
)

func Test_all(t *testing.T) {

	engine := NewDefaultEngine()
	// for i := 1; i <= 10000; i++ {
	// 	engine.AddTextLog(`   java.sql.SQLException:  ddduse them aalav`)
	// }
	// time.Sleep(time.Duration(5) * time.Second)

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
