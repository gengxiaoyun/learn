package mylinux

import (
	"os/exec"
	"bytes"
	"fmt"
)

// show output
func Command(cmd *exec.Cmd) error {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout=&out
	cmd.Stderr=&stderr
	err:=cmd.Run()
	if err != nil{
		//fmt.Println(fmt.Sprint(err)+": "+stderr.String())
		return err
	}
	fmt.Println(out.String())
	return nil
}

// linux command  as root
func CmdRoot(command string) error {
	cmd:=exec.Command("sudo","su","root","-c",command)
	err:=Command(cmd)
	if err != nil{
		return err
	}
	fmt.Println("succeeded")
	return nil
}

// linux command
func Cmd(command,dir string) error {
	cmd:=exec.Command("/bin/bash","-c",command)
	cmd.Dir=dir
	err:=Command(cmd)
	if err != nil{
		return err
	}
	fmt.Println("succeeded")
	return nil
}
