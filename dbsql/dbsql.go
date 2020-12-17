package dbsql

import (
	"strings"
	"fmt"
	"os"
	"bufio"
	"io"
)

// change master_host and master_port
func ChangeSql(report_ip,report_port string) error{
	f,err := os.Open("./dbsql/slave.sql")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer f.Close()

	out,err := os.OpenFile("./dbsql/slave.sql"+".mdf", os.O_RDWR|os.O_CREATE, 0777)
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
		if strings.Contains(newline,"192.168.186.132") {
			newline = strings.Replace(newline, "192.168.186.132", report_ip, 1)
		}
		if strings.Contains(newline,"3306") {
			newline = strings.Replace(newline, "3306", report_port, 1)
		}
		_,err = out.WriteString(newline+"\n")
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
	}

	err = os.Remove("./dbsql/slave.sql")
	if err != nil{
		return err
	}
	err = os.Rename("./dbsql/slave.sql"+".mdf","./dbsql/slave.sql")
	if err != nil{
		return err
	}

	return nil
}