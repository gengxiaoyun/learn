package prepare

import "testing"

func TestStartMysql(t *testing.T) {
	err = StartMysql()
	if err != nil{
		t.Fatal("failed")
	}
}
