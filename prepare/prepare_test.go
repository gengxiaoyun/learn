package prepare

import "testing"

const(
	filename = "github.com/gengxiaoyun/learn/testfile/bbb"
	pathtmp = "github.com/gengxiaoyun/learn/testfile/bbb001"
)

func TestChangefile(t *testing.T) {
	err:=Changefile(filename,pathtmp)
	if err!=nil{
		t.Fatal("failed")
	}
	t.Log("succeeded")
}