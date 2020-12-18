package dbsql

import (
	"testing"
)

func TestChangeSql(t *testing.T) {
	var (
		err error
		report_ip string
		report_port string
	)
	report_ip = "192.168.186.137"
	report_port = "3306"

	err = ChangeSql(report_ip,report_port)
	if err != nil{
		t.Fatal("failed")
	}
}
