package main

import (
	"testing"
	"io/ioutil"
	"fmt"
	"strings"
	"os"
	"database/sql"
	"os/exec"
	"bytes"
)

func TestUntargz(t *testing.T) {
	Untargz(srcfile,"unzipfile/")
	_,err:=os.Stat("unzipfile/mysql-5.7.31-linux-glibc2.12-x86_64")
	if os.IsNotExist(err) {
		t.Fatal("failed")
	}
	t.Log("succeeded")
}

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
	pd:=get_pd(mysql_dir+"data/mysql.err")
	if pd==""{
		t.Fatal("failed to get temporary password")
	}
	t.Log("get temporary password successfully")
}

func TestReadline(t *testing.T) {
	Readline("aaa")
	bytes,err:=ioutil.ReadFile("aaa.mdf")
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


func TestDbconnect(t *testing.T) {
	Dbconnect()
	path:=strings.Join([]string{uname,":",pwd,"@tcp(",ip,":",port,")/",dbname,"?charset=utf8"},"")
	db,_:=sql.Open("mysql",path)
	defer db.Close()
	err:=db.Ping()
	if err!=nil{
		fmt.Println(err)
	}
	rows,err:=db.Query("select * from `test`;")
	defer rows.Close()
	if err!=nil{
		t.Fatal("failed")
	}
	for rows.Next() {
		var id int
		var name string
		err:=rows.Scan(&id,&name)
		if err!=nil{
			t.Fatal("failed")
		}
	}
	t.Log("succeeded")
}

