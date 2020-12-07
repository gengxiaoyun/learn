package masterslave

import "testing"

const(
	destFile = "/home/gengxy/mysql01/"
	baseDir = destFile +"mysql-5.7.31-linux-glibc2.12-x86_64/"
	dir = ""
	dirInit = baseDir +"bin/"
	dirCp = baseDir +"mysqld_multi/"
	mysqlFile3306 = dirCp +"mysqld3306"
	mysqlFile3307 = dirCp +"mysqld3307"
)

func TestMyMulti(t *testing.T) {
	err := MyMulti(dirInit,dir,mysqlFile3306,mysqlFile3307)
	if err != nil{
		t.Fatal("failed")
	}
}
