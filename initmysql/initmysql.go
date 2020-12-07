package initmysql

import (
	"learn/mylinux"
	"learn/basic"
	"regexp"
	"strings"
	"fmt"
	"bufio"
	"io"
	"os"
)

const(
	initCommand = "./mysqld --initialize --user=mysql"
	setupCommand = "./mysql_ssl_rsa_setup"
	safeCommand = "./mysqld_safe --user=mysql &"
	psCommand = "ps -ef|grep mysql"
	changePdCommand = "mysqladmin -uroot -p"
	mysqlPassword = "mysql"
)

// initialize
func InitMysql(baseDir,user,group,dataDir,dirInit,dir string) error {
	var err error
	initCmd := initCommand + " --basedir=" + baseDir + " --datadir=" +
		dataDir + " --log-error=" + dataDir +"mysql.err"
	err = mylinux.Cmd(initCmd, dirInit)
	if err != nil{
		return err
	}

	err = basic.MyOwn(user,group,dataDir,dir)
	if err != nil{
		return err
	}
	err = basic.MyMod(dataDir,dir)
	if err != nil{
		return err
	}

	errCmd := "less " + dataDir + "mysql.err"
	err = mylinux.Cmd(errCmd,dir)
	if err != nil{
		return err
	}

	// setup
	setupCmd := setupCommand + " --datadir=" + dataDir
	err = mylinux.Cmd(setupCmd, dirInit)
	if err != nil{
		return err
	}

	err = mylinux.Cmd(safeCommand, dirInit)
	if err != nil{
		return err
	}

	err = mylinux.Cmd(psCommand, dirInit)
	if err != nil{
		return err
	}

	// reset password
	strPd := GetPassword(dataDir +"mysql.err")
	fmt.Println("temporary password: ", strPd)
	pdCmd := changePdCommand + strPd +" password " + mysqlPassword
	err = mylinux.Cmd(pdCmd, dirInit)
	if err != nil{
		return err
	}

	return nil
}


// get temporary password
func GetPassword(filename string) string {
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



