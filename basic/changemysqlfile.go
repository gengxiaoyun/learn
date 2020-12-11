package basic

import (
	"fmt"
	"learn/mylinux"
	"strings"
	"io"
	"os"
	"bufio"
)


// change file
// --/etc/init.d/mysql.server
func ChangeMysqlFile(filename string) error {
	f,err := os.Open(filename)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer f.Close()

	out,err := os.OpenFile(filename+".mdf", os.O_RDWR, 0777)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer out.Close()
	buf := bufio.NewReader(f)
	newline := ""
	for {
		line,_,err := buf.ReadLine()
		if err == io.EOF{
			break
		}
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		newline = string(line)
		if newline == "basedir=" {
			newline = strings.Replace(newline, "basedir=", "basedir=/home/gengxy/mysql01/mysql.server-5.7.31-linux-glibc2.12-x86_64/", 1)

		}
		if newline == "datadir="{
			newline = strings.Replace(newline,"datadir=","datadir=/home/gengxy/mysql01/mysql.server-5.7.31-linux-glibc2.12-x86_64/mysqld_multi/mysqld3306/data/",1)
		}
		_,err = out.WriteString(newline+"\n")
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}

func ChangeFile(srcServer,destServer,strCommand,srcMdf,dir string) error {
	// mysql.server.server->mysql.server
	cpCmd := fmt.Sprintf(`cp "%s" "%s"`, srcServer, destServer)
	err := mylinux.CmdRoot(cpCmd)
	if err != nil{
		return err
	}
	// set basedir and datadir
	err = MyMod(destServer,dir)
	if err != nil{
		return err
	}

	err = mylinux.Cmd(strCommand,dir)
	if err != nil{
		return err
	}

	err = MyMod(srcMdf,dir)
	if err != nil{
		return err
	}

	err = ChangeMysqlFile(destServer)
	if err != nil{
		return err
	}

	mvCmd := fmt.Sprintf(`sudo mv "%s" "%s"`, srcMdf, destServer)
	err = mylinux.Cmd(mvCmd,dir)
	if err != nil{
		return err
	}

	return nil
}

