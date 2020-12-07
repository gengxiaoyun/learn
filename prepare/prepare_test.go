package prepare

import "testing"

const(
	filename = "/home/gengxy/Goproject/src/learn/testfile/bbb"
	pathTmp = "/home/gengxy/Goproject/src/learn/testfile/bbb001"
	file = "/home/gengxy/Goproject/src/learn/testfile/bbbb"
)

func TestChangeConfFile(t *testing.T) {
	err := ChangeConfFile(filename,pathTmp)
	if err!=nil{
		t.Fatal("failed")
	}
}

func TestCopyConfFile(t *testing.T) {
	err := CopyConfFile(filename,file)
	if err!=nil{
		t.Fatal("failed")
	}
}