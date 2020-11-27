package dbsql

import (
	"strings"
	"database/sql"
	"fmt"
	"io/ioutil"
	_"github.com/go-sql-driver/mysql"
)

const(
	uname="root"
	pwd="mysql"
	ip="127.0.0.1"
	dbname="mysql"
)

func Dbconnect(port,file string) {
	path:=strings.Join([]string{uname,":",pwd,"@tcp(",ip,":",port,")/",dbname,"?charset=utf8&multiStatements=true"},"")
	db,_:=sql.Open("mysql",path)
	defer db.Close()
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err:=db.Ping()
	if err!=nil{
		fmt.Println("open database fail",err.Error())
		return
	}
	fmt.Println("connect success")
	sqlbytes,err := ioutil.ReadFile(file)
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