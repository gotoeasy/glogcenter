package ldb

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_all(t *testing.T) {

	engine := NewDefaultEngine()
	for i := 1; i <= 300000; i++ {
		engine.AddTextLog("date", fmt.Sprintf(`DEBUG ==> Preparing: SELECT id,order_code,order_time,verification_code,tenant_code,tenant_name,credentials_type,credentials_number,mobilephone,license_plate_number,variety,origin_place,unit_code,unit_price,area,begin_time,end_time,lease_time,iz_pay_fees,pay_amount,pay_method,valid_until,order_use_status,member_id,create_by,create_time,update_by,update_time,sys_org_code FROM f1530007 WHERE (order_use_status%d = ?)		 lav%d`, i, i), "sssss_ssss")
	}
	time.Sleep(time.Duration(30) * time.Second)

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
