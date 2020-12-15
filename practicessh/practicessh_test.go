package practicessh

import (
	"testing"
	"os"
)

var(
	srcFile string
	destFile string
	file string
	dataDir string
	hostIP string
	port string
)

const(
	portNum = 22
	user = "root"
	pass = "root"
	exportStr = "/export GO111MODULE=on/a\\export PATH=$PATH:/usr/local/mysql/bin"
	cmdSource = "source /etc/profile"
	DbUser = "mysql"
	DbGroup = "mysql"
	cmdGroup = "groupadd " + DbGroup
	cmdUser = `useradd -p "openssl passwd -1 -salt "some" user" -r -g mysql mysql`
	sLib = "apt-cache search libaio"
	iLib = "apt-get install libaio1"
	DFileA = "/usr/local/data/"
	DFileB = "/usr/local/log/"
)

func TestUnTar(t *testing.T) {
	srcFile = "/home/gengxy/mysql01/mysql.tar.gz"
	destFile = "./testfile/unzipfile/"

	err = os.Mkdir(destFile,os.ModePerm)
	if err != nil{
		t.Fatal("Mkdir failed")
	}
	err = UnTar(srcFile,destFile)
	if err != nil{
		t.Fatal("Unzip failed")
	}
}

func TestChangeMysqlServerFile(t *testing.T) {
	file = "./testfile/aaa"
	dataDir = "/mysqldata/mysql3306/data/"

	err = ChangeMysqlServerFile(file,dataDir)
	if err != nil{
		t.Fatal("failed")
	}
}

func TestBasicWork(t *testing.T) {
	hostIP = "192.168.186.132"
	srcFile = "/home/gengxy/mysql/"
	destFile = "/usr/local/"
	file = "/usr/local/mysql/"
	dataDir = "/mysqldata/mysql3306/data/"
	port = "3306"

	sshConn,err := MySshConnect(hostIP,portNum,user,pass)
	if err != nil{
		t.Fatal("failed")
	}
	err = BasicWork(sshConn,srcFile,destFile,exportStr,cmdSource,cmdGroup,cmdUser,
		DbUser,DbGroup,file,sLib,iLib,dataDir,port )
	if err != nil{
		t.Fatal("failed")
	}
}

func TestBasicWorkSlave(t *testing.T) {
	var baseDir string
	hostIP = "192.168.186.136"
	srcFile = "/home/gengxy/mysql/"
	baseDir = "/usr/local/mysql/"
	dataDir = "/mysqldata/mysql3308/data/"
	port = "3308"

	sshConn,err := MySshConnect(hostIP,portNum,user,pass)
	if err != nil{
		t.Fatal("failed")
	}
	err = BasicWorkSlave(sshConn,srcFile,baseDir,DFileA,dataDir,DFileB,exportStr,cmdSource,
		cmdGroup,cmdUser,DbUser,DbGroup,sLib,iLib,port)
	if err != nil{
		t.Fatal("failed")
	}
}

func TestMyMulti(t *testing.T) {
	var(
		connectStr string
		sqlFileMaster string
		logErrDir string
	)
	hostIP = "192.168.186.132"
	connectStr = "unix(/mysqldata/mysql3306/mysql.sock)"
	sqlFileMaster = "./dbsql/master.sql"
	dataDir = "/mysqldata/mysql3306/data/"
	logErrDir = "/mysqldata/mysql3306/log/"
	port = "3306"

	sshConn,err := MySshConnect(hostIP,portNum,user,pass)
	if err != nil{
		t.Fatal("failed")
	}
	err = MyMulti(sshConn,connectStr,sqlFileMaster,dataDir,DFileA,logErrDir,DFileB,port)
	if err != nil{
		t.Fatal("failed")
	}
}