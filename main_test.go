package main

import (
	"testing"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"os/exec"
	"bytes"
)

var(
	file1="testfile/unzipfile/"
	file2="testfile/unzipfile/mysql-5.7.31-linux-glibc2.12-x86_64"
	file3="testfile/aaa"
	file4="testfile/aaa.mdf"

)

func TestUntargz(t *testing.T) {
	err:=os.Mkdir(file1,os.ModePerm)
	if err!=nil{
		fmt.Println(err)
	}
	Untargz(srcfile,file1)
	_,err1:=os.Stat(file2)
	if os.IsNotExist(err1) {
		t.Fatal("failed")
	}
	t.Log("succeeded")
}

func TestAdduser(t *testing.T) {
	user2 := "sudo useradd -r -g testUser testUser"
	Adduser(user2,"testUser","test")
	command:="cat /etc/passwd|grep testUser"
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

func TestGet_pd(t *testing.T) {
	pd:=Get_pd(mysql_dir+"data/mysql.err")
	if pd==""{
		t.Fatal("failed to get temporary password")
	}
	t.Log("get temporary password successfully")
}

func TestReadline(t *testing.T) {
	out,err:=os.Create(file4)
	if err!=nil{
		fmt.Println(err)
	}
	defer out.Close()
	Readline(file3)
	bytes,err:=ioutil.ReadFile(file4)
	if err!=nil{
		fmt.Println(err)
	}
	f:=string(bytes)
	strings.TrimSpace(f)
	if f==""{
		t.Fatal("failed")
	}
	t.Log("succeeded")
}

