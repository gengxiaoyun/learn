package practicessh

import (
	"regexp"
	"strings"
	"fmt"
	"bufio"
	"io"
	"os"
	"github.com/romberli/go-util/linux"
)

const(
	initCommand = "/usr/local/mysql/bin/mysqld --initialize --user=mysql"
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

	errCmd := "less " + log_error
	_,output,err := sshConn.ExecuteCommand(errCmd)
	if err != nil{
		return err
	}
	fmt.Println(output)

	// setup
	setupCmd := setupCommand + " --datadir=" + dataDir
	_,output,err = sshConn.ExecuteCommand(setupCmd)
	if err != nil{
		return err
	}
	fmt.Println(output)

	_,output,err = sshConn.ExecuteCommand(safeCommand+port)
	if err != nil{
		return err
	}
	fmt.Println(output)

	err = sshConn.CopyFromRemote(log_error, "./testfile")
	if err != nil{
		return err
	}

	// reset password
	strPd,err := GetPassword("./testfile/mysqld.log")
	if err != nil{
		return err
	}
	pdCmd := changePdCommand + "'" + strPd + "'" +" password " + mysqlPassword
	_,_,err = sshConn.ExecuteCommand(pdCmd)
	if err != nil{
		return err
	}

	return nil
}

// get temporary password
func GetPassword(filename string) (string,error) {
	f,err := os.Open(filename)
	if err != nil{
		return "",err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var pd_str []string
	for {
		line,err := buf.ReadString('\n')
		if err == io.EOF{
			break
		}
		if err != nil{
			return "",err
		}
		if strings.Contains(line,"root@localhost:"){
			Regexp := regexp.MustCompile("(.*?)(root@localhost: )(.*?)\n$")
			pd_str = Regexp.FindStringSubmatch(line)
		}
	}
	return string(pd_str[3]),nil
}

