package dbsql

import (
	"testing"
)

const(
	port="3306"
	sql_file="dbsql/master.sql"
)

func TestDbConnect(t *testing.T) {
	err := DbConnect(port,sql_file)
	if err != nil{
		t.Fatal("failed")
	}
}
