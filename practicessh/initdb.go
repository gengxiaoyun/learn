package practicessh

import (
	"github.com/romberli/go-util/linux"
	"time"
)

const(
	initCommand = "/usr/local/mysql/bin/mysqld --initialize-insecure --user=mysql"
	setupCommand = "/usr/local/mysql/bin/mysql_ssl_rsa_setup"
	safeCommand = "mysqld_multi start "
)

// initialize
func InitMysql(sshConn *linux.MySSHConn,baseDir,user,group,dataDir,port string) error {
	var err error
	log_error := "/mysqldata/mysql" + port + "/log/mysqld.log"

	initCmd := initCommand + " --basedir=" + baseDir + " --datadir=" +
		dataDir + " --log-error=" + log_error
	_,_,err = sshConn.ExecuteCommand(initCmd)
	if err != nil{
		return err
	}

	err = OwnAndMod(sshConn,user,group,dataDir)
	if err != nil{
		return err
	}

	// setup
	setupCmd := setupCommand + " --datadir=" + dataDir
	_,_,err = sshConn.ExecuteCommand(setupCmd)
	if err != nil{
		return err
	}

	_,_,err = sshConn.ExecuteCommand(safeCommand+port)
	if err != nil{
		return err
	}

	time.Sleep(time.Duration(15)*time.Second)

	err = sshConn.CopyFromRemote(log_error, "./testfile")
	if err != nil{
		return err
	}

	return nil
}

