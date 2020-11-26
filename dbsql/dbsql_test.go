package dbsql

import (
	"strings"
	"database/sql"
	"fmt"
	"testing"
)

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
		t.Errorf("failed")
	}
	for rows.Next() {
		var id int
		var name string
		err:=rows.Scan(&id,&name)
		if err!=nil{
			t.Errorf("failed")
		}
	}
}