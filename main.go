package main

import (
	"fmt"
	//"github.com/gengxiaoyun/learn/mylinux"
	//"github.com/gengxiaoyun/learn/practicessh"
	//"github.com/gengxiaoyun/learn/dbsql"

	"log"
	"flag"
	"runtime"
	"os"
	"learn/practicessh"
	"learn/dbsql"
	"learn/mylinux"
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

	//cmdChangeProfile = "sed -i "/export\ GO111MODULE=on/a\export\ PATH\=\$PATH\:\/usr\/local\/mysql\/bin" /etc/profile"
	exportStr = "/'export GO111MODULE=on'/a\\export PATH=$PATH:/usr/local/mysql/bin"
	//cmdChangeProfile = "echo `export PATH=$PATH:/usr/local/mysql/bin` >> /etc/profile"
	cmdSource = "source /etc/profile"
	cmdGroup = "groupadd "+DbGroup
	cmdUser = `useradd -p "openssl passwd -1 -salt "some" user" -r -g mysql mysql`

	sLib = "apt-cache search libaio"
	iLib = "apt-get install libaio1"

	dir = ""

	//sqlFileResetPd = "dbsql/setpasswd.sql"
	sqlFileMaster = "dbsql/master.sql"
	sqlFileSlave = "dbsql/slave.sql"
	sqlFileTest = "dbsql/test.sql"

	// copy mysql data to local
	DFileB = "/usr/local/data/"

	rmCommand = "rm -rf "

	startCommand = "/usr/local/mysql/bin/mysqld_multi start "
	reportCommand = "mysqld_multi report"
)

var(
	address string
	user string
	pass string
	err error
	logFileName = flag.String("./log", "InstallMysql.log", "Log file name")
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
	//err = practicessh.UnTar(fileNameSource,fileNameDest)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//mvCmd := fmt.Sprintf(`sudo mv "%s" "%s"`, fileNameI, fileNameN)
	//err = mylinux.Cmd(mvCmd,dir)
	//if err != nil{
	//	fmt.Println(err.Error())
	//}
	//log.Println("unzip succeed")

	for i:=0;i<len(arr);i++ {
		dataDir := "/mysqldata/mysql" + arr[i][1] + "/data/"
		err = practicessh.ChangeMysqlServerFile(serverFile,dataDir)
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		log.Println("mysqlfile prepared")

		sshConn,err := practicessh.MySshConnect(arr[i][0],p,user,pass)
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		log.Println("ssh connect succeedd")

		if i==0 {
			err = sshConn.CopyToRemote(srcCnfFile,destCnfFile)
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			err = sshConn.CopyToRemote(srcMysqlFile, destMysqlFile)
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("my.cnf and mysql.server file copy succeedd")
			err = practicessh.BasicWork(sshConn,srcFile,destFile,exportStr,cmdSource,cmdGroup,cmdUser,
				DbUser,DbGroup,baseDir,sLib,iLib,dataDir,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("BasicWork succeedd")


			//psCommand := "ps -ef|grep mysql"
			//_,output,err := sshConn.ExecuteCommand(psCommand)
			//if err != nil{
			//	fmt.Println(err.Error())
			//	return err
			//}
			//fmt.Println("psCommand:\n")
			//fmt.Println(output)
			//log_error := "/mysqldata/mysql" + arr[i][1] + "/log/mysqld.log"
			//err = sshConn.CopyFromRemote(log_error, "./testfile")
			//if err != nil{
			//	fmt.Println(err.Error())
			//	return err
			//}
			//log.Println("InitMysql succeedd")


			err = practicessh.InitMysql(sshConn,baseDir,DbUser,DbGroup,dataDir,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("InitMysql succeedd")

			//pwd := practicessh.GetPassword("./testfile/mysqld.log")
			////err = practicessh.ResetPasswd("127.0.0.1",arr[i][1],pwd,sqlFileResetPd)
			//err = practicessh.ResetPasswd(sshConn,pwd)
			//if err != nil{
			//	fmt.Println(err.Error())
			//	return err
			//}
			//fmt.Println("===============ResetPasswd succeedd==============")


			err = practicessh.MyMulti(sshConn,sqlFileMaster,dataDir, destFile,"127.0.0.1",arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("MyMulti succeedd")
			rmCmd := rmCommand + destFile +"data/auto.cnf"
			err = mylinux.Cmd(rmCmd,dir)
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("remove auto.cnf succeedd")
			err = dbsql.ChangeSql(arr[i][0],arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("ChangeSql succeedd")
			fmt.Println("ChangeSql succeedd")

		} else {
			if arr[i][0] == arr[0][0] {
				err = practicessh.CreateSomeDir(sshConn,dataDir,arr[i][1],DbUser,DbGroup)
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				err = sshConn.CopyToRemote(DFileB,dataDir)
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				_, _, err = sshConn.ExecuteCommand(startCommand+arr[i][1])
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				fmt.Println("BasicWorkSlave.startCommand ok")
				err = practicessh.DbConnect(arr[i][0],arr[i][1],sqlFileSlave)
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				log.Println("slave DbConnect succeedd")
				fmt.Println("slave DbConnect succeedd")

			} else {
				err = sshConn.CopyToRemote(srcCnfFile,destCnfFile)
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				err = sshConn.CopyToRemote(srcMysqlFile, destMysqlFile)
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				log.Println("my.cnf and mysql.server file copy succeedd")
				fmt.Println("my.cnf and mysql.server file copy succeedd")


				err = practicessh.BasicWorkSlave(sshConn,srcFile,baseDir,DFileB,dataDir,exportStr,cmdSource,cmdGroup,cmdUser,
					DbUser,DbGroup,sLib,iLib,arr[i][1])
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				log.Println("BasicWorkSlave succeedd")
				fmt.Println("BasicWorkSlave succeedd")
				err = practicessh.DbConnect(arr[i][0],arr[i][1],sqlFileSlave)
				if err != nil{
					fmt.Println(err.Error())
					return err
				}
				log.Println("slave DbConnect succeedd")
				fmt.Println("slave DbConnect succeedd")
				}

			}

	}

	for i:=0;i<len(arr);i++ {
		sshConn,err := practicessh.MySshConnect(arr[i][0],p,user,pass)
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		log.Println("ssh connect succeedd")
		if i==0{
			_, _, err = sshConn.ExecuteCommand(reportCommand)
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("reportCommand succeedd")
			err = practicessh.DbConnect(arr[i][0],arr[i][1],sqlFileTest)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
			log.Println("master DbConnect succeedd")
		} else {
			err = practicessh.DbConnect(arr[i][0],arr[i][1],sqlFileSlave)
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("slave DbConnect succeedd")
		}
	}
	return nil
}


func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "InstallMysql start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//write log
	log.Printf("Log start! File:%v \n", "InstallMysql.log")

	err = StartMysql()
	if err != nil{
		log.Println(err.Error())
	}
}
