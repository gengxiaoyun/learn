package main

//import (
//	"archive/tar"
//	"os"
//	"fmt"
//	"compress/gzip"
//	"io"
//	"strings"
//	"os/exec"
//	"bytes"
//	"bufio"
//	"database/sql"
//	_"github.com/go-sql-driver/mysql"
//	"io/ioutil"
//)
//
//
//var(
//	destfile = "/home/gengxy/mysql01/"
//	srcfile = "/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"
//	mysql_dir = destfile+"mysql-5.7.31-linux-glibc2.12-x86_64/"
//	password = "mysql"
//	src1 = "my.cnf"
//	dest1 = "/etc/my.cnf"
//	src2 = mysql_dir+"support-files/mysql.server"
//	dest2 = "/etc/init.d/mysql"
//	src3 = "/etc/init.d/mysql.mdf"
//)
//
//const(
//	uname="root"
//	pwd="mysql"
//	ip="127.0.0.1"
//	port="3306"
//	dbname="mysql"
//)
//
//func Createfile(name string) (*os.File,error) {
//	err:=os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name,"/")]),0755)
//	if err!=nil{
//		return nil,err
//	}
//	return os.Create(name)
//}
//// unzip
//func Untargz(srcfile string, destfile string) error {
//	fr, err := os.Open(srcfile)
//	if err != nil {
//		return err
//	}
//	defer fr.Close()
//	// gzip read
//	gr, err := gzip.NewReader(fr)
//	if err != nil {
//		return err
//	}
//	defer gr.Close()
//	// tar read
//	tr := tar.NewReader(gr)
//	// read file
//	for {
//		h, err := tr.Next()
//		if err!=nil{
//			if err == io.EOF {
//				break
//			}else {
//				return err
//			}
//		}
//		filename := destfile + h.Name
//		if h.Typeflag == tar.TypeDir {
//			if err:=os.MkdirAll(filename,os.FileMode(h.Mode));err!=nil{
//				return err
//			}
//		}else{
//			file,err := Createfile(filename)
//			if err != nil{
//				return err
//			}
//			_, err = io.Copy(file, tr)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
//// os.exec command
//func Cmd(command,password string) {
//	cmd:=exec.Command("/bin/bash","-c",command)
//	err:=cmd.Start()
//	if err != nil{
//		fmt.Println(err.Error())
//	}
//	cmd.Wait()
//	ps:=exec.Command("echo",password)
//	grep:=exec.Command("passwd","-stdin",password)
//	r,w:=io.Pipe()
//	defer r.Close()
//	defer w.Close()
//	ps.Stdout=w
//	grep.Stdin=r
//	var buffer bytes.Buffer
//	grep.Stdout=&buffer
//	_=ps.Start()
//	_=grep.Start()
//	ps.Wait()
//	w.Close()
//	grep.Wait()
//	io.Copy(os.Stdout,&buffer)
//
//}
//
////change file --/etc/init.d/mysql
//func Readline(filename string) {
//	f,err:=os.Open(filename)
//	if err!=nil{
//		fmt.Println(err.Error())
//	}
//	defer f.Close()
//
//	out,err:=os.OpenFile(filename+".mdf", os.O_RDWR, 0777)
//	if err!=nil{
//		fmt.Println(err.Error())
//	}
//	defer out.Close()
//	buf:=bufio.NewReader(f)
//	newline:=""
//	for {
//		line,_,err:=buf.ReadLine()
//		if err==io.EOF{
//			break
//		}
//		if err!=nil{
//			fmt.Println(err.Error())
//		}
//		newline = string(line)
//		if newline=="basedir=" {
//			newline = strings.Replace(newline, "basedir=", "basedir=/home/gengxy/mysql01/mysql-5.7.31-linux-glibc2.12-x86_64/", 1)
//
//		}
//		if newline=="datadir="{
//			newline = strings.Replace(newline,"datadir=","datadir=/home/gengxy/mysql01/mysql-5.7.31-linux-glibc2.12-x86_64/data/",1)
//		}
//		_,err1:=out.WriteString(newline+"\n")
//		if err1!=nil{
//			fmt.Println(err1.Error())
//		}
//
//
//	}
//
//	//err1:=os.Remove(filename)
//	//if err1!=nil{
//	//	fmt.Println(err1.Error())
//	//}
//	//err2:=os.Rename(filename+".mdf",filename)
//	//if err2!=nil{
//	//	fmt.Println(err2.Error())
//	//}
//
//}
//
////out,err:=cmd.Output()
////if err != nil{
////	fmt.Println(err.Error())
////}
////return string(out),err
//
//// mysql initialize
//func Cmdinit(command,password string) {
//	cmd:=exec.Command("/bin/bash","-c",command)
//	cmd.Dir="/home/gengxy/mysql01/mysql-5.7.31-linux-glibc2.12-x86_64/bin/"
//	err:=cmd.Start()
//	if err != nil{
//		fmt.Println(err.Error())
//	}
//	cmd.Wait()
//	ps:=exec.Command("echo",password)
//	grep:=exec.Command("passwd","-stdin",password)
//	r,w:=io.Pipe()
//	defer r.Close()
//	defer w.Close()
//	ps.Stdout=w
//	grep.Stdin=r
//	var buffer bytes.Buffer
//	grep.Stdout=&buffer
//	_=ps.Start()
//	_=grep.Start()
//	ps.Wait()
//	w.Close()
//	grep.Wait()
//	io.Copy(os.Stdout,&buffer)
//
//}
//
//// output error message
//func Checkerr(command string) {
//	cmd:=exec.Command("/bin/bash","-c",command)
//	var out bytes.Buffer
//	var stderr bytes.Buffer
//	cmd.Stdout=&out
//	cmd.Stderr=&stderr
//	cmd.Dir="/home/gengxy/mysql01/mysql-5.7.31-linux-glibc2.12-x86_64/bin/"
//	err:=cmd.Run()
//	if err != nil{
//		fmt.Println(fmt.Sprint(err)+": "+stderr.String())
//		return
//	}
//	fmt.Println("Result: "+out.String())
//}
//
//func Dbconnect() {
//	path:=strings.Join([]string{uname,":",pwd,"@tcp(",ip,":",port,")/",dbname,"?charset=utf8&multiStatements=true"},"")
//	db,_:=sql.Open("mysql",path)
//	defer db.Close()
//	db.SetConnMaxLifetime(100)
//	db.SetMaxIdleConns(10)
//	err:=db.Ping()
//	if err!=nil{
//		fmt.Println("open database fail")
//		return
//	}
//	fmt.Println("connect success")
//	sqlbytes,err := ioutil.ReadFile("test.sql")
//	if err!=nil{
//		fmt.Println(err.Error())
//		return
//	}
//	sqltable := string(sqlbytes)
//	fmt.Println(sqltable)
//	_,err1 := db.Exec(sqltable)
//	if err1!=nil{
//		fmt.Println(err1.Error())
//		return
//	}
//}
//
//
//
//
//func main() {
//	// unzip
//	//Untargz(srcfile,destfile)
//	//fmt.Println("un tar.gz ok")
//
//	// add_user_chown_chmod
//	//user := "sudo useradd mysql"
//	//chown1 := "sudo chown -R mysql:mysql "+mysql_dir
//	//chown2 := "sudo chown -R mysql "+mysql_dir
//	//chmod := "sudo chmod -R 777 "+mysql_dir
//	//Cmd(user,password)
//	//Cmd(chown1,password)
//	//Cmd(chown2,password)
//	//Cmd(chmod,password)
//	//
//	//// install_libaio  --ubuntu
//	//ilibaio := "sudo apt-get install libaio"
//	//Cmd(ilibaio,password)
//	//
//	//// my.cnf
//	//cp_cmd1 := fmt.Sprintf(`sudo cp "%s" "%s"`,src1,dest1)
//	//Cmd(cp_cmd1,password)
//	//// mysql.server->mysql
//	//cp_cmd2 := fmt.Sprintf(`sudo cp "%s" "%s"`,src2,dest2)
//	//Cmd(cp_cmd2,password)
//	//
//	//// add_path
//	//cmdstr1 := "sudo chmod -R 777 /etc/init.d/mysql"
//	//cmdstr2 := "sudo touch /etc/init.d/mysql.mdf"
//	//cmdstr3 := "sudo chmod -R 777 /etc/init.d/mysql.mdf"
//	//Cmd(cmdstr1,password)
//	//Cmd(cmdstr2,password)
//	//Cmd(cmdstr3,password)
//	//Readline(dest2)
//	//cmdremove := fmt.Sprintf(`sudo mv "%s" "%s"`,src3,dest2)
//	//Cmd(cmdremove,password)
//	//
//	//// initialize
//	//cmdstr4:="./mysqld --initialize --user=mysql --basedir="+mysql_dir+" --datadir="+mysql_dir+"data"
//	//////checkerr(cmdstr4)
//	//Cmdinit(cmdstr4,password)
//	//
//	//cmdstr5:="./mysql_ssl_rsa_setup --datadir="+mysql_dir+"data"
//	//////checkerr(cmdstr5)
//	//Cmdinit(cmdstr5,password)
//	//
//	//cmdstr6:="./mysqld_safe --user=mysql &"
//	//////checkerr(cmdstr6)
//	//Cmdinit(cmdstr6,password)
//	//
//	//cmdstr7:="ps -ef|grep mysql"
//	////checkerr(cmdstr7)
//	//Cmd(cmdstr7,password)
//
//	//login
//	cmdstr8:="mysql -uroot -p"
//	Cmd(cmdstr8,password)
//
//	// connect database and create table
//	//Dbconnect()
//}


