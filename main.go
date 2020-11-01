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
)

func createfile(name string) (*os.File,error) {
	err:=os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name,"/")]),0755)
	if err!=nil{
		return nil,err
	}
	return os.Create(name)
}

func untargz(srcfile string, destfile string) error {
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
	// 读取文件
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
			file,err := createfile(filename)
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

func Cmd(command,password string) {
	cmd:=exec.Command("/bin/bash","-c",command)
	err:=cmd.Start()
	if err != nil{
		fmt.Println(err.Error())
	}
	cmd.Wait()
	ps:=exec.Command("echo",password)
	grep:=exec.Command("passwd","-stdin",password)
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

func Cmdfile(command,password string,f func(string)) {
	cmd:=exec.Command("/bin/bash","-c",command)
	err:=cmd.Start()
	if err != nil{
		fmt.Println(err.Error())
	}
	cmd.Wait()
	ps:=exec.Command("echo",password)
	grep:=exec.Command("passwd","-stdin",password)
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
	f("/etc/init.d/mysql")
}


func readline(filename string) {
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
	//err1:=os.Remove(filename)
	//if err1!=nil{
	//	fmt.Println(err1.Error())
	//}
	//err2:=os.Rename(filename+".mdf",filename)
	//if err2!=nil{
	//	fmt.Println(err2.Error())
	//}


}


//out,err:=cmd.Output()
//if err != nil{
//	fmt.Println(err.Error())
//}
//return string(out),err


func checkerr(command string) {
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


func main() {
	// file read
	//destfile:="/home/gengxy/mysql01/"
	//srcfile:="/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"
	//mysql_dir:=destfile+"mysql-5.7.31-linux-glibc2.12-x86_64/"
	//password:="mysql"

	//untargz(srcfile,destfile)
	//fmt.Println("un tar.gz ok")

	// add_user_chown_chmod
	//user:="sudo useradd mysql"
	//chown1:="sudo chown -R mysql:mysql "+mysql_dir
	//chown2:="sudo chown -R mysql "+mysql_dir
	//chmod:="sudo chmod -R 755 "+mysql_dir
	//Cmd(user,password)
	//Cmd(chown1,password)
	//Cmd(chown2,password)
	//Cmd(chmod,password)

	// install_libaio
	//slibaio:="sudo yum search libaio"
	//ilibaio:="sudo yum install libaio"
	//yum(slibaio)
	//yum(ilibaio)

	// copy_file
	//src := mysql_dir+"support-files/mysql.server"
	//dest := "/etc/init.d/mysql"
	//sercmd := fmt.Sprintf(`sudo cp "%s" "%s"`,src,dest)
	//Cmd(sercmd,password)

	// add_path
	//cmdstr1:="sudo chmod -R a=rwx /etc/init.d/mysql"
	//Cmd(cmdstr1,password)
	//cmdstr2:="sudo touch /etc/init.d/mysql.mdf"
	//Cmd(cmdstr2,password)
	//cmdstr3:="sudo chmod -R a=rwx /etc/init.d/mysql.mdf"
	//Cmd(cmdstr3,password)
	//readline("/etc/init.d/mysql")
	//src2 := "/etc/init.d/mysql.mdf"
	//dest2 := "/etc/init.d/mysql"
	//cmdremove := fmt.Sprintf(`sudo mv "%s" "%s"`,src2,dest2)
	//Cmd(cmdremove,password)

	// initialize
	//cmdstr4:=mysql_dir+"bin/mysqld --initialize --user=mysql --basedir="+mysql_dir+" --datadir="+mysql_dir+"data"
	//checkerr(cmdstr4)

	//b,err:=cmd.Output()
	//fmt.Println(string(b),err)

	//cmdstr5:=mysql_dir+"bin/mysql_ssl_rsa_setup --datadir="+mysql_dir+"data"
	//checkerr(cmdstr5)
	//cmdstr6:=mysql_dir+"bin/mysqld_safe --user=mysql &"
	//checkerr(cmdstr6)
	//cmdstr7:="ps -ef|grep mysql"
	//checkerr(cmdstr7)
	//
	//cmdstr8:="mysql -uroot -p"
	//checkerr(cmdstr8)
}


