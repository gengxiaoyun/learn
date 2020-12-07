package practicessh

import (
	"regexp"
	"strings"
	"fmt"
	"bufio"
	"io"
	"os"
)

const(
	initCommand = "/usr/local/mysql/bin/mysqld --initialize --user=mysql"
	setupCommand = "/usr/local/mysql/bin/mysqld/mysql_ssl_rsa_setup"
	safeCommand = "/usr/local/mysql/bin/mysqld/mysqld_safe --user=mysql &"
	psCommand = "ps -ef|grep mysql"
	changePdCommand = "mysqladmin -uroot -p"
	mysqlPassword = "mysql"
)

// initialize
func (conn *MySSHConn) InitMysql(baseDir,user,group,dataDir string) error {
	var err error
	initCmd := initCommand + " --basedir=" + baseDir + " --datadir=" +
		dataDir + " --log-error=" + dataDir +"mysql.err"
	_,_,err = conn.ExecuteCommand(initCmd)
	if err != nil{
		return err
	}

	err = conn.MyOwn(user,group,dataDir)
	if err != nil{
		return err
	}
	err = conn.MyMod(dataDir)
	if err != nil{
		return err
	}

	errCmd := "less " + dataDir + "mysql.err"
	_,output,err := conn.ExecuteCommand(errCmd)
	if err != nil{
		return err
	}
	fmt.Println(output)

	// setup
	setupCmd := setupCommand + " --datadir=" + dataDir
	_,_,err = conn.ExecuteCommand(setupCmd)
	if err != nil{
		return err
	}

	_,output,err = conn.ExecuteCommand(safeCommand)
	if err != nil{
		return err
	}
	fmt.Println(output)

	_,output,err = conn.ExecuteCommand(psCommand)
	if err != nil{
		return err
	}
	fmt.Println(output)

	// reset password
	strPd := conn.GetPassword(dataDir +"mysql.err")
	fmt.Println("temporary password: ", strPd)
	pdCmd := changePdCommand + strPd +" password " + mysqlPassword
	_,_,err = conn.ExecuteCommand(pdCmd)
	if err != nil{
		return err
	}

	return nil
}


// get temporary password
func (conn *MySSHConn) GetPassword(filename string) string {
	f,err := os.Open(filename)
	if err != nil{
		fmt.Println(err)
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

