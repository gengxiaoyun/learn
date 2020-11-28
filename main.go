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
	"regexp"
	"github.com/gengxiaoyun/learn/linux"
	"github.com/gengxiaoyun/learn/prepare"
	"github.com/gengxiaoyun/learn/dbsql"
)

const(
	destfile = "/home/gengxy/mysql01/"
	srcfile = "/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"
	mysql_dir = destfile+"mysql-5.7.31-linux-glibc2.12-x86_64/"
	pathtmp = "prepare/my.cnf001"
	src1 = "prepare/my.cnf"
	dest1 = "/etc/my.cnf"
	src2 = mysql_dir+"support-files/mysql.server"
	dest2 = "/etc/init.d/mysql"
	src3 = "/etc/init.d/mysql.mdf"
	dir = ""
	dir_init = mysql_dir+"bin/"
	dir_cp = mysql_dir+"mysqld_multi/"

	port = "3306"
	port_2 = "3307"

	sql_file = "dbsql/test.sql"
	sql_file2 = "dbsql/test2.sql"
	sql_file3 = "dbsql/test3.sql"

	mysqlfile = dir_cp+"mysqld3306"
	mysqlfile2 = dir_cp+"mysqld3307"
)

// create file
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

// Add user and set password
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
	fmt.Println("Add user successfully")
}

// change file
// --/etc/init.d/mysql
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
			newline = strings.Replace(newline,"datadir=","datadir=/home/gengxy/mysql01/mysql-5.7.31-linux-glibc2.12-x86_64/mysqld_multi/mysqld3306/data/",1)
		}
		_,err1:=out.WriteString(newline+"\n")
		if err1!=nil{
			fmt.Println(err1.Error())
		}
	}
}

// get temporary password
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


func main() {

	Untargz(srcfile,destfile)
	fmt.Println("un tar.gz ok")

	group := "groupadd mysql"
	linux.Cmd_root(group)
	user := "sudo useradd -r -g mysql mysql"
	Adduser(user,"mysql","user")
	mk_cmd := "sudo mkdir -p "+mysql_dir+"mysqld_multi/mysqld3306/"
	linux.Cmd(mk_cmd,dir)
	chown1 := "sudo chown -R mysql:mysql "+mysql_dir
	linux.Cmd(chown1,dir)
	chown2 := "sudo chown -R mysql "+mysql_dir
	linux.Cmd(chown2,dir)
	chmod := "sudo chmod -R 777 "+mysql_dir
	linux.Cmd(chmod,dir)

	// install_libaio
	slibaio:="apt-cache search libaio"
	linux.Cmd(slibaio,dir)
	ilibaio := "apt-get install libaio1"
	linux.Cmd_root(ilibaio)

	// change my.cnf
	err:=prepare.Changefile(src1,pathtmp)
	if err!=nil{
		fmt.Println(err.Error())
	}

	// copy my.cnf
	cp_cmd1 := fmt.Sprintf(`cp "%s" "%s"`,src1,dest1)
	linux.Cmd_root(cp_cmd1)
	// mysql.server->mysql
	cp_cmd2 := fmt.Sprintf(`cp "%s" "%s"`,src2,dest2)
	linux.Cmd_root(cp_cmd2)

	// set basedir and datadir
	str_cmd1 := "sudo chmod -R 777 "+dest2
	linux.Cmd(str_cmd1,dir)
	str_cmd2 := "sudo touch "+src3
	linux.Cmd(str_cmd2,dir)
	str_cmd3 := "sudo chmod -R 777 "+src3
	linux.Cmd(str_cmd3,dir)
	Readline(dest2)
	mv_cmd := fmt.Sprintf(`sudo mv "%s" "%s"`,src3,dest2)
	linux.Cmd(mv_cmd,dir)

	// initialize
	init_cmd:="./mysqld --initialize --user=mysql --basedir="+mysql_dir+" --datadir="+mysql_dir+"mysqld_multi/mysqld3306/data/ --log-error="+mysql_dir+"mysqld_multi/mysqld3306/data/mysql.err"
	linux.Cmd(init_cmd,dir_init)
	chown := "sudo chown -R mysql:mysql "+mysql_dir+"mysqld_multi/mysqld3306/data"
	linux.Cmd(chown,dir)
	data_cmd := "sudo chmod -R 777 "+mysql_dir+"mysqld_multi/mysqld3306/data"
	linux.Cmd(data_cmd,dir)
	mysql_err_cmd:="less "+mysql_dir+"mysqld_multi/mysqld3306/data/mysql.err"
	linux.Cmd(mysql_err_cmd,dir)

	str_cmd4:="./mysql_ssl_rsa_setup --datadir="+mysql_dir+"mysqld_multi/mysqld3306/data"
	linux.Cmd(str_cmd4,dir_init)

	str_cmd5:="./mysqld_safe --user=mysql &"
	linux.Cmd(str_cmd5,dir_init)

	str_cmd6:="ps -ef|grep mysql"
	linux.Cmd(str_cmd6,dir_init)

	// get temporary password
	str_pd:= Get_pd(mysql_dir+"mysqld_multi/mysqld3306/data/mysql.err")
	fmt.Println("temporary password: ",str_pd)

	// reset password
	change_pd_cmd:="mysqladmin -uroot -p"+str_pd+" password mysql"
	linux.Cmd(change_pd_cmd,dir_init)

	// connect database and create table
	dbsql.Dbconnect(port,sql_file)

	stop_cmd := "mysqld_multi stop 3306"
	linux.Cmd(stop_cmd,dir_init)

	copy_cmd := fmt.Sprintf(`cp -r "%s" "%s"`,mysqlfile,mysqlfile2)
	linux.Cmd(copy_cmd,dir)

	rm_cmd:="sudo rm -rf "+mysqlfile2+"/data/auto.cnf"
	linux.Cmd(rm_cmd,dir)

	start_cmd:="mysqld_multi start 3306,3307"
	linux.Cmd(start_cmd,dir_init)

	time.Sleep(time.Duration(3)*time.Second)

	report_cmd:="mysqld_multi report"
	linux.Cmd(report_cmd,dir_init)

	dbsql.Dbconnect(port,sql_file3)
	dbsql.Dbconnect(port_2,sql_file2)
}
