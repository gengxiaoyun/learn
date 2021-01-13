package dbsql

import (
	"testing"
)

func TestChangeSql(t *testing.T) {
	var (
		err error
		reportIp string
		reportPort string
	)
	reportIp = "192.168.186.137"
	reportPort = "3306"

	err = ChangeSql(reportIp,reportPort)
	if err != nil{
		t.Fatal("failed")
	}
}
