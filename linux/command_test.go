package linux

import (
	"testing"
	"fmt"
	"os"
	"os/exec"
	"bytes"
)

func TestCmd_root(t *testing.T) {
	group := "groupadd testUser"
	Cmd_root(group)
	command:="cat /etc/group | grep testUser"
	cmd:=exec.Command("/bin/bash","-c",command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout=&out
	cmd.Stderr=&stderr
	err:=cmd.Run()
	if err != nil{
		fmt.Println(fmt.Sprint(err)+": "+stderr.String())
	}
	if out.String()==""{
		t.Errorf("failed")
	}
}

func TestCmd(t *testing.T) {
	cmd:="sudo touch 1.txt"
	dir:=""
	Cmd(cmd,dir)
	_,err:=os.Stat("1.txt")
	if os.IsNotExist(err) {
		t.Errorf("failed")
	}
}