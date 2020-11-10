package linux

import (
	"os/exec"
	"bytes"
	"fmt"
)

var(
	destfile = "/home/gengxy/mysql01/"
	mysql_dir = destfile+"mysql-5.7.31-linux-glibc2.12-x86_64/"
)

// show output
func Command(cmd *exec.Cmd){
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout=&out
	cmd.Stderr=&stderr
	err:=cmd.Run()
	if err != nil{
		fmt.Println(fmt.Sprint(err)+": "+stderr.String())
		return
	}
	fmt.Println(out.String())
}

// linux command  as root
func Cmd_root(command string){
	cmd:=exec.Command("sudo","su","root","-c",command)
	Command(cmd)
	fmt.Println("succeeded")
}

// linux command
func Cmd(command,dir string){
	cmd:=exec.Command("/bin/bash","-c",command)
	cmd.Dir=dir
	Command(cmd)
	fmt.Println("succeeded")
}
