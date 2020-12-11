package initmysql

import "testing"

const(
	user = "mysql.server"
	group = "mysql.server"

	destFile = "/home/gengxy/mysql01/"
	baseDir = destFile +"mysql.server-5.7.31-linux-glibc2.12-x86_64/"
	filePath = baseDir +"mysqld_multi/mysqld3306/"
	dataDir = filePath+"data/"
	dir = ""
	dirInit = baseDir +"bin/"
)

func TestInitMysql(t *testing.T) {
	err := InitMysql(baseDir,user,group,dataDir,dirInit,dir)
	if err != nil{
		t.Fatal("failed")
	}
}
