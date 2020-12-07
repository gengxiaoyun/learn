package main

import (
	"fmt"
	//"github.com/gengxiaoyun/learn/mylinux"
	//"github.com/gengxiaoyun/learn/prepare"
	//"github.com/gengxiaoyun/learn/dbsql"
	//"learn/prepare"
	//"learn/dbsql"
	//"learn/basic"
	//"learn/initmysql"
	//"learn/masterslave"

	"log"
	"flag"
	"runtime"
	"os"
	"learn/practicessh"
	"learn/mylinux"
	"learn/dbsql"
)

const(
	//srcFile = "/home/gengxy/mysql/mysql.tar.gz"
	//destFile = "/home/gengxy/mysql01/"
	//baseDir = destFile +"mysql/"

	srcFile = "practicessh/my.cnf"
	destFile = "/etc/my.cnf"

	fileNameSource = "/home/gengxy/mysql/mysql.tar.gz"
	fileNameDest = "/usr/local/mysql01"
	fileDest = "/usr/local/"
	baseDir  = fileDest + "mysql/"

	//cmdUnzip = "tar -xvf mysql.tar.gz -C /usr/local"

	DbUser = "mysql"
	DbGroup = "mysql"
	//userPassword = "user"

	cmdChangeProfile = "echo `export MYSQL_HOME=/usr/local/mysql` >> /etc/profile && echo `export PATH=$MYSQL_HOME/bin:$PATH` >> /etc/profile"
	cmdSource = "source /etc/profile"
	cmdGroup = "groupadd "+DbGroup
	cmdUser = `useradd -p "openssl passwd -1 -salt "some" user" -r -g mysql mysql`

	sLib = "apt-cache search libaio"
	iLib = "apt-get install libaio1"

	srcServer = baseDir +"support-files/mysql.server"
	destServer = "/etc/init.d/mysql"
	srcMdf = "/etc/init.d/mysql.mdf"
	strCommand = "sudo touch "+ srcMdf

	dir = ""

	sqlFileMaster = "dbsql/master.sql"
	sqlFileSlave = "dbsql/slave.sql"
	sqlFileTest = "dbsql/test.sql"

	DFileA = "/usr/local/Installdata"
	DFileB = "/usr/local/data"
	rmCommand = "sudo rm -rf "

	reportCommand = "mysqld_multi report"
)

var(
	address string
	user string
	pass string
	err error
	logFileName = flag.String("./log", "golangMysql.log", "Log file name")
)

func init() {
	flag.StringVar(&address, "address", "192.168.186.132:3306", "set ip and port")
	flag.StringVar(&user, "user", "root", "set username")
	flag.StringVar(&pass, "pass", "root", "set password")
}


func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "golangMysql start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//write log
	log.Printf("Log start! File:%v \n", "golangMysql.log")


	//err = basic.UnTar(srcFile,destFile)
	//if err != nil{
	//	log.Printf("Un tar.gz failed. message:%s", err.Error())
	//}
	//log.Println("un tar.gz succeed")
	//
	//err = basic.UserAndGroup(groupCommand,userCommand,user,group,
	//	baseDir,userPassword,filePath,sLib,iLib,dir)
	//if err != nil{
	//	log.Printf("Prepare work failed. message:%s", err.Error())
	//}
	//log.Println("Prepare work succeed")
	//
	//// create my.cnf.toml
	//err = prepare.ChangeConfFile(srcCnf,pathTmp)
	//if err != nil{
	//	log.Printf("Change my.cnf.toml failed. message:%s", err.Error())
	//}
	//err = prepare.CopyConfFile(srcCnf,destCnf)
	//if err != nil{
	//	log.Printf("Copy my.cnf.toml failed. message:%s", err.Error())
	//}
	//log.Println("Create my.cnf.toml succeed")
	//
	//err = basic.ChangeFile(srcServer,destServer,strCommand,srcMdf,dir)
	//if err != nil{
	//	log.Printf("Move mysql.server failed. message:%s", err.Error())
	//}
	//
	//err = initmysql.InitMysql(baseDir,user,group,dataDir,dirInit,dir)
	//if err != nil{
	//	log.Printf("Initialize mysql failed. message:%s", err.Error())
	//}
	//log.Println("Initialize mysql succeed")
	//
	//// connect database and create table
	//err = dbsql.DbConnect(port3306, sqlFileMaster)
	//if err != nil{
	//	log.Printf("Connect databases failed. message:%s", err.Error())
	//}
	//log.Println("Connect databases succeed")
	//
	//err = masterslave.MyMulti(dirInit,dir,mysqlFile3306,mysqlFile3307)
	//if err != nil{
	//	log.Printf("Copy master-slave failed. message:%s", err.Error())
	//}
	//
	//err = dbsql.DbConnect(port3306, sqlFileTest)
	//if err != nil{
	//	log.Printf("Insert data on master failed. message:%s", err.Error())
	//}
	//
	//err = dbsql.DbConnect(port3307, sqlFileSlave)
	//if err != nil{
	//	log.Printf("Build master-slave failed. message:%s", err.Error())
	//}
	//log.Println("Build master-slave succeed")

	//flag.Parse()



	if err = practicessh.Init(); err != nil {
		log.Printf("conf.Init() err:%+v", err)
		//fmt.Println("failed")
	}
	arr,err := practicessh.Flex(address)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("======================================")

	for i:=0;i<len(arr);i++ {
		dataDir := "/mysqldata/mysql" + arr[i][1] + "/data/"
		sshConn,err := practicessh.NewMySSHConn(arr[i][0],arr[i][1],user,pass)
		if err != nil{
			fmt.Println(err.Error())
		}
		err = sshConn.CopyToRemote(srcFile,destFile)
		if err != nil{
			fmt.Println(err.Error())
		}
		if i==0 {
			err = sshConn.BasicWork(fileNameSource,fileNameDest,fileNameDest, fileDest,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
				DbUser,DbGroup,baseDir,sLib,iLib,srcServer,destServer,strCommand,srcMdf,dataDir,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
			}

			err = sshConn.InitMysql(baseDir,DbUser,DbGroup,dataDir)
			if err != nil{
				fmt.Println(err.Error())
			}

			err = sshConn.MyMulti(sqlFileMaster,sqlFileTest,baseDir,DFileA,dataDir,DFileB,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
			}

			rmCmd := rmCommand + DFileB +"/auto.cnf"
			err = mylinux.Cmd(rmCmd,dir)
			if err != nil{
				fmt.Println(err.Error())
			}

			err = dbsql.ChangeSql(arr[i][0],arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
			}

		} else if{
			err = sshConn.BasicWorkSlave(DFileA,DFileB,fileNameDest, fileDest,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
				DbUser,DbGroup,baseDir,sLib,iLib,srcServer,destServer,strCommand,srcMdf,dataDir,arr[i][1])
			if err != nil{
				fmt.Println(err.Error())
			}

			err = sshConn.DbConnect(arr[i][1],sqlFileSlave)
			if err != nil{
				fmt.Println(err.Error())
			}
		}
	}


	for i:=0;i<len(arr);i++ {
		sshConn,err := practicessh.NewMySSHConn(arr[i][0],arr[i][1],user,pass)
		if err != nil{
			fmt.Println(err.Error())
		}
		if i==0{
			_, _, err = sshConn.ExecuteCommand(reportCommand)
			if err != nil{
				fmt.Println(err.Error())
			}

			err = sshConn.DbConnect(arr[i][1], sqlFileTest)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if{
			err = sshConn.DbConnect(arr[i][1],sqlFileSlave)
			if err != nil{
				fmt.Println(err.Error())
			}
		}
	}
}
