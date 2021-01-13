package practicessh

import "testing"

func TestInitMysql(t *testing.T) {
	var(
		err error
		hostIP string
		portNum int
		user string
		pass string
		baseDir string
		myUser string
		myGroup string
		dataDir string
		port string
	)
	hostIP = "192.168.186.137"
	portNum = 22
	user = "root"
	pass = "Abc727364"
	baseDir = "/usr/local/mysql/"
	myUser = "mysql"
	myGroup = "mysql"
	dataDir = "/mysqldata/mysql3306/data/"
	port = "3306"

	sshConn,err := MySshConnect(hostIP,portNum,user,pass)
	if err != nil{
		t.Fatal("failed")
	}
	err = InitMysql(sshConn,baseDir,myUser,myGroup,dataDir,port)
	if err != nil{
		t.Fatal("failed")
	}
}
