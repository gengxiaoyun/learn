package prepare

import "testing"

const(
	filename = "/home/gengxy/Goproject/src/learn/testfile/bbb"
	pathtmp = "/home/gengxy/Goproject/src/learn/testfile/bbb001"
)

func TestChangefile(t *testing.T) {
	err:=Changefile(filename,pathtmp)
	if err!=nil{
		t.Fatal("failed")
	}
	t.Log("succeeded")
}