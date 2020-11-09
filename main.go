package main

import (
	"archive/tar"
	"os"
	"fmt"
	"compress/gzip"
	"io"
	"strings"
	"os/exec"
	"bytes"
	"bufio"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"io/ioutil"
	"regexp"
)


var(
	destfile = "/home/gengxy/mysql01/"
	srcfile = "/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"
	mysql_dir = destfile+"mysql-5.7.31-linux-glibc2.12-x86_64/"
	src1 = "my.cnf"
	dest1 = "/etc/my.cnf"
	src2 = mysql_dir+"support-files/mysql.server"
	dest2 = "/etc/init.d/mysql"
	src3 = "/etc/init.d/mysql.mdf"
)

const(
	uname="root"
	pwd="mysql"
	ip="127.0.0.1"
	port="3306"
	dbname="mysql"
)

func Createfile(name string) (*os.File,error) {
	err:=os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name,"/")]),0755)
	if err!=nil{
		return nil,err
	}
	return os.Create(name)
}
// unzip
func Untargz(srcfile string, destfile string) error {
	fr, err := os.Open(srcfile)
	if err != nil {
		return err
	}
	defer fr.Close()
	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()
	// tar read
	tr := tar.NewReader(gr)
	// read file
	for {
		h, err := tr.Next()
		if err!=nil{
			if err == io.EOF {
				break
			}else {
				return err
			}
		}
		filename := destfile + h.Name
		if h.Typeflag == tar.TypeDir {
			if err:=os.MkdirAll(filename,os.FileMode(h.Mode));err!=nil{
				return err
			}
		}else{
			file,err := Createfile(filename)
			if err != nil{
				return err
			}
			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Cmd(command string) {
	cmd:=exec.Command("/bin/bash","-c",command)
	err:=cmd.Start()
	if err!=nil{
		fmt.Println(err)
	}
	err1:=cmd.Wait()
	if err1!=nil{
		fmt.Println(err1)
	}
}

func Adduser(command,username,password string) {
	cmd:=exec.Command("/bin/bash","-c",command)
	err:=cmd.Start()
	if err != nil{
		fmt.Println(err.Error())
	}
	cmd.Wait()
	ps:=exec.Command("echo",password)
	grep:=exec.Command("passwd","-stdin",username)
	r,w:=io.Pipe()
	defer r.Close()
	defer w.Close()
	ps.Stdout=w
	grep.Stdin=r
	var buffer bytes.Buffer
	grep.Stdout=&buffer
	_=ps.Start()
	_=grep.Start()
	ps.Wait()
	w.Close()
	grep.Wait()
	io.Copy(os.Stdout,&buffer)

}

func  Cmd_output(command string) {
	cmd:=exec.Command("/bin/bash","-c",command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout=&out
	cmd.Stderr=&stderr
	err:=cmd.Run()
	if err != nil{
		fmt.Println(fmt.Sprint(err)+": "+stderr.String())
		return
	}
	fmt.Println("Result: "+out.String())
}

func Cmd_root(command string) {
	cmd:=exec.Command("sudo","su","root","-c",command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout=&out
	cmd.Stderr=&stderr
	err:=cmd.Run()
	if err != nil{
		fmt.Println(fmt.Sprint(err)+": "+stderr.String())
		return
	}
	fmt.Println("Successful")
}


//change file --/etc/init.d/mysql
func Readline(filename string) {
	f,err:=os.Open(filename)
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer f.Close()

	out,err:=os.OpenFile(filename+".mdf", os.O_RDWR, 0777)
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer out.Close()
	buf:=bufio.NewReader(f)
	newline:=""
	for {
		line,_,err:=buf.ReadLine()
		if err==io.EOF{
			break
		}
		if err!=nil{
			fmt.Println(err.Error())
		}
		newline = string(line)
		if newline=="basedir=" {
			newline = strings.Replace(newline, "basedir=", "basedir=/home/gengxy/mysql01/mysql-5.7.31-linux-glibc2.12-x86_64/", 1)

		}
		if newline=="datadir="{
			newline = strings.Replace(newline,"datadir=","datadir=/home/gengxy/mysql01/mysql-5.7.31-linux-glibc2.12-x86_64/data/",1)
		}
		_,err1:=out.WriteString(newline+"\n")
		if err1!=nil{
			fmt.Println(err1.Error())
		}
	}
}

// mysql initialize
func Checkerr(command string) {
	cmd:=exec.Command("/bin/bash","-c",command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout=&out
	cmd.Stderr=&stderr
	cmd.Dir=mysql_dir+"bin/"
	err:=cmd.Run()
	if err != nil{
		fmt.Println(fmt.Sprint(err)+": "+stderr.String())
		return
	}
	fmt.Println("initialize successful")
}

func Get_pd(filename string) string{
	f,err:=os.Open(filename)
	if err!=nil{
		fmt.Println(err)
	}
	defer f.Close()
	buf:=bufio.NewReader(f)
	var pd_str []string
	for {
		line,err:=buf.ReadString('\n')
		if err==io.EOF{
			break
		}
		if err!=nil{
			fmt.Println(err.Error())
		}
		if strings.Contains(line,"root@localhost:"){
			Regexp:=regexp.MustCompile("(.*?)(root@localhost: )(.*?)\n$")
			pd_str=Regexp.FindStringSubmatch(line)
		}
	}
	return string(pd_str[3])
}

func Dbconnect() {
	path:=strings.Join([]string{uname,":",pwd,"@tcp(",ip,":",port,")/",dbname,"?charset=utf8&multiStatements=true"},"")
	db,_:=sql.Open("mysql",path)
	defer db.Close()
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err:=db.Ping()
	if err!=nil{
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connect success")
	sqlbytes,err := ioutil.ReadFile("test.sql")
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	sqltable := string(sqlbytes)
	fmt.Println(sqltable)
	_,err1 := db.Exec(sqltable)
	if err1!=nil{
		fmt.Println(err1.Error())
		return
	}

	fmt.Println("show data:")
	rows,err2:=db.Query("select * from `test`;")
	defer rows.Close()
	if err2!=nil{
		fmt.Println(err2.Error())
		return
	}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id,&name)
		fmt.Println(id,name)
	}
}


func main() {
	// unzip
	Untargz(srcfile,destfile)
	fmt.Println("un tar.gz ok")

	// add_user_chown_chmod
	group := "groupadd mysql"
	user := "sudo useradd -r -g mysql mysql"
	chown1 := "sudo chown -R mysql:mysql "+mysql_dir
	chown2 := "sudo chown -R mysql "+mysql_dir
	chmod := "sudo chmod -R 777 "+mysql_dir
	Cmd_root(group)
	Adduser(user,"mysql","user")
	Cmd(chown1)
	Cmd(chown2)
	Cmd(chmod)

	// install_libaio
	slibaio:="apt-cache search libaio"
	ilibaio := "apt-get install libaio1"
	Cmd_output(slibaio)
	Cmd_root(ilibaio)

	// my.cnf
	cp_cmd1 := fmt.Sprintf(`cp "%s" "%s"`,src1,dest1)
	Cmd_root(cp_cmd1)
	// mysql.server->mysql
	cp_cmd2 := fmt.Sprintf(`cp "%s" "%s"`,src2,dest2)
	Cmd_root(cp_cmd2)

	// add_path
	str_cmd1 := "sudo chmod -R 777 "+dest2
	str_cmd2 := "sudo touch "+src3
	str_cmd3 := "sudo chmod -R 777 "+src3
	Cmd(str_cmd1)
	Cmd(str_cmd2)
	Cmd(str_cmd3)
	Readline(dest2)
	mv_cmd := fmt.Sprintf(`sudo mv "%s" "%s"`,src3,dest2)
	Cmd(mv_cmd)

	// initialize
	init_cmd:="./mysqld --initialize --user=mysql --basedir="+mysql_dir+" --datadir="+mysql_dir+"data"
	Checkerr(init_cmd)
	chown := "sudo chown -R mysql:mysql "+mysql_dir+"data"
	Cmd(chown)
	data_cmd := "sudo chmod -R 777 "+mysql_dir+"data"
	Cmd(data_cmd)
	mysql_err_cmd:="less "+mysql_dir+"data/mysql.err"
	Cmd_output(mysql_err_cmd)

	str_cmd4:="./mysql_ssl_rsa_setup --datadir="+mysql_dir+"data"
	Checkerr(str_cmd4)

	str_cmd5:="./mysqld_safe --user=mysql &"
	Checkerr(str_cmd5)

	str_cmd6:="ps -ef|grep mysql"
	Checkerr(str_cmd6)

	// get temporary password
	str_pd:= Get_pd(mysql_dir+"data/mysql.err")
	fmt.Println("temporary password: ",str_pd)

	// reset password
	change_pd_cmd:="mysqladmin -uroot -p"+str_pd+" password mysql"
	Checkerr(change_pd_cmd)

	// connect database and create table
	Dbconnect()

}
