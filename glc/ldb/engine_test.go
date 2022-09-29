package ldb

import (
	"fmt"
	"glc/cmn"
	"glc/ldb/sysmnt"
	"log"
	"testing"
	"time"
)

func Test_GetIP(t *testing.T) {
	log.Println(cmn.GetLocalIp())
}

func Test_GetSubDirs(t *testing.T) {
	rs := sysmnt.GetStorageList()

	for i := 0; i < len(rs.Data); i++ {
		log.Println(rs.Data[i])
	}

}

func Test_all(t *testing.T) {

	engine := NewDefaultEngine()
	for i := 1; i <= 10000; i++ {
		engine.AddTextLog("date", fmt.Sprintf(`DEBUG ==> Preparing: SELECT id,aaa,bbb, ccc,ddd,eee,fff,ggg  FROM abcde WHERE (aaa%d = ?)		 lav%d`, i, i), "sssss_ssss")
	}
	time.Sleep(time.Duration(10) * time.Second)

	// for i := 1; i <= 10000; i++ {
	// 	engine.AddTextLog(`   java.sql.SQLException:   them aalav`)
	// }
	// time.Sleep(time.Duration(5) * time.Second)

	// for i := 1; i <= 10000; i++ {
	// 	engine.AddTextLog(`  java.sql.SQLException: them`)
	// }
	// time.Sleep(time.Duration(5) * time.Second)

	rs := engine.Search(`              them java     `, 5, 0, true)
	log.Println("共查到", rs.Total, "件")
	for _, v := range rs.Data {
		log.Println(v.Id, v.Text)
	}

}
