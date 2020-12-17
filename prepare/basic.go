package prepare

import (
	"log"
	"fmt"
	"flag"
	"github.com/gengxiaoyun/learn/mylinux"
	"github.com/gengxiaoyun/learn/practicessh"
	"github.com/gengxiaoyun/learn/dbsql"
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

	sqlFileMaster = "./dbsql/master.sql"
	sqlFileSlave = "./dbsql/slave.sql"
	sqlFileTest = "./dbsql/test.sql"

	// copy mysql data to local
	DFileA = "/usr/local/data/"
	DFileB = "/usr/local/log/"
	mkCommand = "mkdir " + DFileB
	rmCommand = "rm -rf "

	startCommand = "mysqld_multi start "
	reportCommand = "mysqld_multi report"
)

var(
	address string
	user string
	pass string
	err error
)

func init() {
	flag.StringVar(&address, "address", "192.168.186.132:3306", "set ip and port")
	flag.StringVar(&user, "user", "root", "set username")
	flag.StringVar(&pass, "pass", "root", "set password")
}

func StartMysql() error{

	if err = practicessh.Init(); err != nil {
		log.Printf("conf.Init() err:%+v", err)
		return err
	}
	arr,err := practicessh.Flex(address)
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
			connectStr := "unix(/mysqldata/mysql"+arr[i][1]+"/mysql.sock)"

			err = mylinux.Cmd(mkCommand,dir)
			if err != nil{
				return err
			}
			err = practicessh.MyMulti(sshConn,connectStr,sqlFileMaster,dataDir,DFileA,logErrDir,DFileB,arr[i][1])
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
			if arr[i][0] == arr[0][0] {
				err = practicessh.CreateSomeDir(sshConn,dataDir,arr[i][1],DbUser,DbGroup)
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

				//connectStr := "tcp("+arr[i][0]+":"+arr[i][1]+")"
				connectStr := "unix(/mysqldata/mysql"+arr[i][1]+"/mysql.sock)"
				err = practicessh.DbConnect(connectStr,sqlFileSlave)
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

				err = practicessh.BasicWorkSlave(sshConn,srcFile,baseDir,DFileA,dataDir,DFileB,exportStr,cmdSource,cmdGroup,cmdUser,
					DbUser,DbGroup,sLib,iLib,arr[i][1])
				if err != nil{
					return err
				}
				log.Println("BasicWorkSlave succeedd")

				//connectStr := "tcp("+arr[i][0]+":"+arr[i][1]+")"
				connectStr := "unix(/mysqldata/mysql"+arr[i][1]+"/mysql.sock)"
				err = practicessh.DbConnect(connectStr,sqlFileSlave)
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
		//connectStr := "tcp("+arr[i][0]+":"+arr[i][1]+")"
		connectStr := "unix(/mysqldata/mysql"+arr[i][1]+"/mysql.sock)"
		if i==0{
			_, _, err = sshConn.ExecuteCommand(reportCommand)
			if err != nil{
				return err
			}
			log.Println("reportCommand succeedd")
			err = practicessh.DbConnect(connectStr,sqlFileTest)
			if err != nil {
				return err
			}
			log.Println("master DbConnect succeedd")
		} else {
			err = practicessh.DbConnect(connectStr,sqlFileSlave)
			if err != nil{
				return err
			}
			log.Println("slave DbConnect succeedd")
		}
	}
	return nil
}