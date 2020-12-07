package basic

import (
	"os"
	"testing"
)

const(
	srcFile = "/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"
	destFile = "learn/testfile/unzipfile/"
	filePath = "learn/testfile/data"
	file = "learn/testfile/"
	dir = ""

	group = "testUser"
	user = "testUser"
	userPassword = "user"
	groupCommand = "groupadd "+group
	userCommand = "sudo useradd -r -g "+group+" "+user

	sLib = "apt-cache search libaio"
	iLib = "apt-get install libaio1"

)

var err error

func TestUnTar(t *testing.T) {
	err = os.Mkdir(destFile,os.ModePerm)
	if err != nil{
		t.Fatal("Mkdir failed")
	}
	err = UnTar(srcFile,destFile)
	if err != nil{
		t.Fatal("Unzip failed")
	}

	//t.Log("succeeded")
}

func TestUserAndGroup(t *testing.T) {
	err = UserAndGroup(groupCommand,userCommand,user,group,
		filePath,userPassword,filePath,sLib,iLib,dir)
	if err != nil{
		t.Fatal("Add User And Group failed")
	}
}
