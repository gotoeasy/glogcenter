package ldb

import (
	"fmt"
	"glc/com"
	"glc/ldb/sysmnt"
	"glc/ldb/tokenizer"
	"testing"
	"time"

	"github.com/gotoeasy/glang/cmn"
)

func Test_CutWords(t *testing.T) {
	ws := tokenizer.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造，Java和Go都学得不错，Java和Go都不错的")
	cmn.Println(ws)
}

func Test_GetIP(t *testing.T) {
	cmn.Println(com.GetLocalIp())
}
func Test_GetSubDirs(t *testing.T) {
	rs := sysmnt.GetStorageList()

	for i := 0; i < len(rs.Data); i++ {
		cmn.Println(rs.Data[i])
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

	rs := engine.Search(`              them java     `, "", "", 5, 0, true)
	cmn.Println("共查到", rs.Total, "件")
	for _, v := range rs.Data {
		cmn.Println(v.Id, v.Text)
	}

}
