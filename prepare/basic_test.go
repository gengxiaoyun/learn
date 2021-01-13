package prepare

import "testing"

func TestStartMysql(t *testing.T) {
	str := []string{"192.168.186.137:3306","192.168.186.137:3307"}
	user := "root"
	pass := "root"
	err = StartMysql(str,user,pass)
	if err != nil{
		t.Fatal("failed")
	}
}
