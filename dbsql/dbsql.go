package dbsql

import (
	"strings"
	"database/sql"
	"fmt"
	"io/ioutil"
	_"github.com/go-sql-driver/mysql"
	"os"
	"bufio"
	"io"
)

const(
	uname="root"
	pwd="mysql.server"
	ip="127.0.0.1"
	dbname="mysql.server"
)


func DbConnect(port,file string) error {
	path := strings.Join([]string{uname,":",pwd,"@tcp(",ip,":",port,")/",dbname,"?charset=utf8&multiStatements=true"},"")
	db,_ := sql.Open("mysql",path)
	defer db.Close()
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err := db.Ping()
	if err != nil{
		fmt.Println("open database fail",err.Error())
		return err
	}
	fmt.Println("database connect success")
	sqlBytes,err := ioutil.ReadFile(file)
	if err != nil{
		return err
	}
	sqlTable := string(sqlBytes)
	fmt.Println(sqlTable)
	_,err = db.Exec(sqlTable)
	if err != nil{
		return err
	}

	rows,err := db.Query("select * from `test`;")
	defer rows.Close()
	if err != nil{
		return err
	}
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id,&name)
		if err != nil{
			return err
		}
		fmt.Println(id,name)
	}
	return nil
}

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
		if strings.Contains(newline,"192.168.186.131") {
			newline = strings.Replace(newline, "192.168.186.131", report_ip, 1)
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