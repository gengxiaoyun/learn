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
	srcCnfFile = "./my.cnf"
	destCnfFile = "/etc/my.cnf"

	fileNameSource = "/home/gengxy/mysql01/mysql.tar.gz"
	fileNameDest = "/home/gengxy/"

	serverFile = fileNameDest + "mysql/support-files/mysql.server"
	srcMysqlFile = "./mysql"
	destMysqlFile = "/etc/init.d/mysql"

	fileNameI = fileNameDest + "mysql-5.7.31-linux-glibc2.12-x86_64/"
	fileNameN = fileNameDest + "mysql/"

	srcFile = fileNameDest + "mysql/"
	destFile = "/usr/local/"

	baseDir  = destFile + "mysql/"

	DbUser = "mysql"
	DbGroup = "mysql"
	p = 22

	cmdChangeProfile = "echo `export MYSQL_HOME=/usr/local/mysql` >> /etc/profile && echo `export PATH=$MYSQL_HOME/bin:$PATH` >> /etc/profile"
	cmdSource = "source /etc/profile"
	cmdGroup = "groupadd "+DbGroup
	cmdUser = `useradd -p "openssl passwd -1 -salt "some" user" -r -g mysql mysql`

	sLib = "apt-cache search libaio"
	iLib = "apt-get install libaio1"

	//srcServer = baseDir +"support-files/mysql.server"
	//destServer = "/etc/init.d/mysql"
	//srcMdf = "/etc/init.d/mysql.mdf"
	//strCommand = "sudo touch "+ srcMdf

	dir = ""

	sqlFileResetPd = "dbsql/setpasswd.sql"
	sqlFileMaster = "dbsql/master.sql"
	sqlFileSlave = "dbsql/slave.sql"
	sqlFileTest = "dbsql/test.sql"

	DFileA = "/usr/local/Installdata"  //copy install dir to local
	DFileB = "/usr/local/data"  // copy mysql.server data to local
	rmCommand = "sudo rm -rf "

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
		err = practicessh.ChangeMysqlFile(serverFile,dataDir)
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		log.Println("mysqlfile prepared")
		//p,_:=strconv.Atoi(arr[i][1])
		sshConn,err := practicessh.MySshConnect(arr[i][0],p,user,pass)
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		log.Println("ssh connect succeedd")
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
		if i==0 {
			err = practicessh.BasicWork(sshConn,srcFile,destFile,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
				DbUser,DbGroup,baseDir,sLib,iLib,dataDir,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("BasicWork succeedd")
			err = practicessh.InitMysql(sshConn,baseDir,DbUser,DbGroup,dataDir,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("InitMysql succeedd")
			pwd := practicessh.GetPassword("./testfile/mysqld.log")
			err = practicessh.ResetPasswd("127.0.0.1",arr[i][1],pwd,sqlFileResetPd)
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			err = practicessh.MyMulti(sshConn,sqlFileMaster,baseDir,DFileA,dataDir,DFileB,arr[i][0],arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("MyMulti succeedd")
			rmCmd := rmCommand + DFileB +"/auto.cnf"
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
		} else {
			err = practicessh.BasicWorkSlave(sshConn,DFileA,baseDir,DFileB,dataDir,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
				DbUser,DbGroup,sLib,iLib,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("BasicWorkSlave succeedd")
			err = practicessh.DbConnect(arr[i][0],arr[i][1],sqlFileSlave)
			if err != nil{
				fmt.Println(err.Error())
				return err
			}
			log.Println("slave DbConnect succeedd")
		}
	}

	for i:=0;i<len(arr);i++ {
		//p,_:=strconv.Atoi(arr[i][1])
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
