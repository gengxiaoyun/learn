package prepare

import (
	"log"
	"github.com/gengxiaoyun/learn/mylinux"
	"github.com/gengxiaoyun/learn/practicessh"
	"github.com/gengxiaoyun/learn/dbsql"
	"fmt"
	"time"
	"os"
)

const(
	fileNameSource = "/home/gengxy/mysql01/mysql.tar.gz"
	fileNameDest = "/home/gengxy/"
	fileNameI = fileNameDest + "mysql-5.7.31-linux-glibc2.12-x86_64/"
	fileNameN = fileNameDest + "mysql/"

	serverFile = fileNameDest + "mysql/support-files/mysql.server"
	srcMysqlFile = "./mysql"
	destMysqlFile = "/etc/init.d/mysql"

	srcCnfFile = "./my.cnf"
	destCnfFile = "/etc/my.cnf"

	srcFile = fileNameDest + "mysql/"
	destFile = "/usr/local/"
	baseDir  = destFile + "mysql/"

	DbUser = "mysql"
	DbGroup = "mysql"
	p = 22

	exportStr = "/export GO111MODULE=on/a\\export PATH=$PATH:/usr/local/mysql/bin"
	cmdSource = "source /etc/profile"
	cmdGroup = "groupadd "+DbGroup
	cmdUser = `useradd -p "openssl passwd -1 -salt "some" user" -r -g mysql mysql`

	sLib = "apt-cache search libaio"
	iLib = "apt-get install libaio1"

	dir = ""

	sqlFileSetPd = "./dbsql/setpassword.sql"
	sqlFileMaster = "./dbsql/master.sql"
	sqlFileSlave = "./dbsql/slave.sql"
	sqlFileTest = "./dbsql/test.sql"
	sqlFileSlaveTest = "./dbsql/slavetest.sql"

	// copy mysql data to local
	DFileA = "/usr/local/data/"
	DFileB = "/usr/local/log"
	rmCommand = "rm -rf "

	startCommand = "mysqld_multi start "
	reportCommand = "mysqld_multi report"
)

var(
	err error
)


func CheckPath(file string) error{
	_,err = os.Stat(file)
	if err == nil{
		return nil
	}

	err = mylinux.Cmd("mkdir " + file,dir)
	if err != nil{
		return err
	}
	return nil
}

func StartMysql(str []string,user,pass string) error{

	if err = practicessh.Init(); err != nil {
		log.Printf("conf.Init() err:%+v", err)
		return err
	}
	arr,err := practicessh.Flex(str)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("my.cnf prepared")

	err = practicessh.UnTar(fileNameSource,fileNameDest)
	if err != nil {
		return err
	}

	mvCmd := fmt.Sprintf(`sudo mv "%s" "%s"`, fileNameI, fileNameN)
	err = mylinux.Cmd(mvCmd,dir)
	if err != nil{
		return err
	}
	log.Println("unzip succeed")

	for i:=0;i<len(arr);i++ {
		dataDir := "/mysqldata/mysql" + arr[i][1] + "/data/"
		err = practicessh.ChangeMysqlServerFile(serverFile,dataDir)
		if err != nil{
			return err
		}
		log.Println("mysqlfile prepared")

		sshConn,err := practicessh.MySshConnect(arr[i][0],p,user,pass)
		if err != nil{
			return err
		}
		log.Println("ssh connect succeedd")

		if i==0 {
			err = sshConn.CopyToRemote(srcCnfFile,destCnfFile)
			if err != nil{
				return err
			}
			err = sshConn.CopyToRemote(srcMysqlFile, destMysqlFile)
			if err != nil{
				return err
			}
			log.Println("my.cnf and mysql.server file copy succeedd")
			err = practicessh.BasicWork(sshConn,srcFile,destFile,exportStr,cmdSource,cmdGroup,cmdUser,
				DbUser,DbGroup,baseDir,sLib,iLib,dataDir,arr[i][1])
			if err != nil{
				return err
			}
			log.Println("BasicWork succeedd")

			err = practicessh.InitMysql(sshConn,baseDir,DbUser,DbGroup,dataDir,arr[i][1])
			if err != nil{
				return err
			}
			log.Println("InitMysql succeedd")

			logErrDir := "/mysqldata/mysql" + arr[i][1] + "/log/"
			err = CheckPath(DFileA)
			if err != nil{
				return err
			}
			err = CheckPath(DFileB)
			if err != nil{
				return err
			}

			err = sshConn.CopyToRemote(sqlFileSetPd,destFile)
			if err != nil{
				return err
			}
			err = sshConn.CopyToRemote(sqlFileMaster,destFile)
			if err != nil{
				return err
			}
			err = sshConn.CopyToRemote(sqlFileTest,destFile)
			if err != nil{
				return err
			}

			err = practicessh.MyMulti(sshConn,dataDir,DFileA,logErrDir,DFileB,arr[i][1])
			if err != nil{
				return err
			}
			log.Println("MyMulti succeeded")

			rmCmd := rmCommand + destFile +"data/auto.cnf"
			err = mylinux.Cmd(rmCmd,dir)
			if err != nil{
				return err
			}
			log.Println("remove auto.cnf succeeded")
			err = dbsql.ChangeSql(arr[i][0],arr[i][1])
			if err != nil{
				return err
			}
			log.Println("ChangeSql succeeded")

		} else {
			setSlaveCommand := "mysql -uroot -pmysql -S /mysqldata/mysql"+arr[i][1]+"/mysql.sock < " + "/usr/local/slave.sql"

			if arr[i][0] == arr[0][0] {
				err = practicessh.CreateSomeDir(sshConn,dataDir,arr[i][1],
					DbUser,DbGroup,cmdGroup,cmdUser)
				if err != nil{
					return err
				}
				err = sshConn.CopyToRemote(DFileA,dataDir)
				if err != nil{
					return err
				}
				err = sshConn.CopyToRemote(DFileB,"/mysqldata/mysql"+arr[i][1]+"/log/")
				if err != nil{
					return err
				}
				err = practicessh.OwnAndMod(sshConn,DbUser,DbGroup,"/mysqldata")
				if err != nil{
					return err
				}
				_, _, err = sshConn.ExecuteCommand(startCommand+arr[i][1])
				if err != nil{
					return err
				}
				log.Println("Slave startCommand ok")
				time.Sleep(time.Duration(15)*time.Second)

				err = sshConn.CopyToRemote(sqlFileSlave,destFile)
				if err != nil{
					return err
				}
				err = sshConn.CopyToRemote(sqlFileSlaveTest,destFile)
				if err != nil{
					return err
				}

				_,_,err = sshConn.ExecuteCommand(setSlaveCommand)
				if err != nil{
					return err
				}
				log.Println("slave DbConnect succeedd")

			} else {
				err = sshConn.CopyToRemote(srcCnfFile,destCnfFile)
				if err != nil{
					return err
				}
				err = sshConn.CopyToRemote(srcMysqlFile, destMysqlFile)
				if err != nil{
					return err
				}
				log.Println("my.cnf and mysql.server file copy succeedd")

				err = practicessh.BasicWorkSlave(sshConn,srcFile,destFile,DFileA,dataDir,DFileB,exportStr,cmdSource,cmdGroup,cmdUser,
					DbUser,DbGroup,sLib,iLib,arr[i][1])
				if err != nil{
					return err
				}
				log.Println("BasicWorkSlave succeedd")

				err = sshConn.CopyToRemote(sqlFileSlave,destFile)
				if err != nil{
					return err
				}
				err = sshConn.CopyToRemote(sqlFileSlaveTest,destFile)
				if err != nil{
					return err
				}

				_,_,err = sshConn.ExecuteCommand(setSlaveCommand)
				if err != nil{
					return err
				}
				log.Println("slave DbConnect succeedd")

			}

		}

	}

	for i:=0;i<len(arr);i++ {
		sshConn,err := practicessh.MySshConnect(arr[i][0],p,user,pass)
		if err != nil{
			return err
		}
		log.Println("ssh connect succeedd")
		if i==0{
			_, output, err := sshConn.ExecuteCommand(reportCommand)
			if err != nil{
				return err
			}
			log.Println(output)
			log.Println("reportCommand succeedd")
			setMasterTestCommand := "mysql -uroot -pmysql < " + "/usr/local/test.sql"
			_,_,err = sshConn.ExecuteCommand(setMasterTestCommand)
			if err != nil{
				return err
			}
			log.Println("master test succeedd")

		} else {
			setSlaveTestCommand := "mysql -uroot -pmysql -S /mysqldata/mysql"+arr[i][1]+"/mysql.sock < " + "/usr/local/slavetest.sql"
			_,output,err := sshConn.ExecuteCommand(setSlaveTestCommand)
			if err != nil{
				return err
			}
			log.Println(output)
			log.Println("slave DbConnect succeedd")

		}
	}
	return nil
}