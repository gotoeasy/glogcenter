package ldb

import (
	"log"
	"testing"
)

func Test_all(t *testing.T) {

	// go func() {

	engine := NewDefaultEngine()
	// for i := 1; i <= 5000000; i++ {
	// 	engine.AddTextLog(`  java.sql.SQLException: View szy-xdqttest.sys_user references inv definer/invoker of vi ddduse them　aalav"`)
	// }
	// time.Sleep(time.Duration(180) * time.Second)

	rs := engine.Search(` java.sql.SQLException: View szy-xdqttest.sys_user references inv definer/invoker of vi ddduse them　aalav  `)
	log.Println("共查到", rs.TotalCount, "件")
	// log.Println(rs.Result)
	// for _, v := range rs.Result {
	// 	log.Println(v.Content)
	// }

	// }()
	// time.Sleep(time.Duration(5) * time.Second)
}
