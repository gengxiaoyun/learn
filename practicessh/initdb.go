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
		fmt.Println(err.Error())
		return err
	}

	err = MyOwn(sshConn,user,group,dataDir)
	if err != nil{
		return err
	}
	err = MyMod(sshConn,dataDir)
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
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(output)

	_,output,err = sshConn.ExecuteCommand(safeCommand+port)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(output)

	err = sshConn.CopyFromRemote(log_error, "./testfile")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	// reset password
	strPd := GetPassword("./testfile/mysqld.log")
	fmt.Println("temporary password: ", strPd)
	pdCmd := changePdCommand + "'" + strPd + "'" +" password " + mysqlPassword
	_,_,err = sshConn.ExecuteCommand(pdCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	return nil
}


// get temporary password
func GetPassword(filename string) string {
	f,err := os.Open(filename)
	if err != nil{
		fmt.Println(err.Error())
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
			fmt.Println(err.Error())
		}
		if strings.Contains(line,"root@localhost:"){
			Regexp := regexp.MustCompile("(.*?)(root@localhost: )(.*?)\n$")
			pd_str = Regexp.FindStringSubmatch(line)
		}
	}
	return string(pd_str[3])
}

