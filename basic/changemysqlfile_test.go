package basic

import "testing"

const(
	srcServer = "learn/testfile/aaa"
	destServer = "learn/testfile/aaa2"
	srcMdf = "learn/testfile/aaa.mdf"
	strCommand = "sudo touch "+ srcMdf
)


func TestChangeFile(t *testing.T) {
	err := ChangeFile(srcServer,destServer,strCommand,srcMdf,dir)
	if err != nil{
		t.Fatal("failed")
	}
}
