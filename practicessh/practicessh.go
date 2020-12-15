package practicessh

import (
	"fmt"
	"compress/gzip"
	"strings"
	"io"
	"os"

	"github.com/romberli/go-util/linux"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"io/ioutil"
	"archive/tar"
	"bufio"
	"time"
)

const (
	uname="root"
	dbPasswd ="mysql"
	//ip="127.0.0.1"
	dbname="mysql"

	mysqlPassword = "mysql"
	changePdCommand = "mysqladmin -uroot -p"
	stopCommand = "mysqld_multi stop "
	startCommand = "mysqld_multi start "

	binFile = "/usr/bin"
	cpMysql = "/usr/local/mysql/bin/mysql"
	cpMyd = "/usr/local/mysql/bin/mysqld"
	cpSafe = "/usr/local/mysql/bin/mysqld_safe"
	cpMyMulti = "/usr/local/mysql/bin/mysqld_multi"
	cpMyDump = "/usr/local/mysql/bin/mysqldump"
	cpMyBinlog = "/usr/local/mysql/bin/mysqlbinlog"
	cpMyCnfEdit = "/usr/local/mysql/bin/mysql_config_editor"
	cpMyPriDef = "/usr/local/mysql/bin/my_print_defaults"
	cpMyAdm = "/usr/local/mysql/bin/mysqladmin"
)



// create file
func CreateFile(name string) (*os.File,error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name,"/")]),0755)
	if err!=nil{
		return nil,err
	}
	return os.Create(name)
}

// unzip
func UnTar(srcFile,destFile string) error {
	fr, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer fr.Close()
	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()
	// tar read
	tr := tar.NewReader(gr)
	// read file
	for {
		h, err := tr.Next()
		if err!=nil{
			if err == io.EOF {
				break
			}else {
				return err
			}
		}
		filename := destFile + h.Name
		if h.Typeflag == tar.TypeDir {
			if err:=os.MkdirAll(filename,os.FileMode(h.Mode));err!=nil{
				return err
			}
		}else{
			file,err := CreateFile(filename)
			if err != nil{
				return err
			}
			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func MySshConnect(hostIP string, portNum int, userName, userPass string) (sshConn *linux.MySSHConn,err error){
	sshConn,err = linux.NewMySSHConn(hostIP, portNum, userName, userPass)
	if err != nil{
		return nil,err
	}
	return sshConn,nil
}

func MyOwn(sshConn *linux.MySSHConn,DbUser,DbGroup,file string) error {
	ownCmd := fmt.Sprintf(`sudo chown -R "%s"."%s" "%s"`,DbUser,DbGroup,file)
	_,_,err := sshConn.ExecuteCommand(ownCmd)
	if err != nil{
		return err
	}
	return nil
}

func MyMod(sshConn *linux.MySSHConn,file string) error {
	modCmd := fmt.Sprintf(`sudo chmod -R 700 "%s"`,file)
	_,_,err := sshConn.ExecuteCommand(modCmd)
	if err != nil{
		return err
	}
	return nil
}

// install libaio
func InstallTool(sshConn *linux.MySSHConn,sLib,iLib string) error {
	_,_,err := sshConn.ExecuteCommand(sLib)
	if err != nil{
		return err
	}
	_,_,err = sshConn.ExecuteCommand(iLib)
	if err != nil{
		return err
	}
	return nil
}

func ChangeMysqlServerFile(filename,dataDir string) error {
	f,err := os.Open(filename)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer f.Close()

	out,err := os.OpenFile("./mysql.server", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer out.Close()
	buf := bufio.NewReader(f)
	newline := ""
	for {
		line,_,err := buf.ReadLine()
		if err == io.EOF{
			break
		}
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
		newline = string(line)
		if newline == "basedir=" {
			newline = strings.Replace(newline, "basedir=", "basedir=/usr/local/mysql/", 1)

		}
		if newline == "datadir="{
			newline = strings.Replace(newline,"datadir=","datadir=" + dataDir,1)
		}
		_,err = out.WriteString(newline+"\n")
		if err != nil{
			fmt.Println(err.Error())
			return err
		}
	}
	err = os.Rename("./mysql.server","./mysql")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func OwnAndMod(sshConn *linux.MySSHConn,DbUser,DbGroup,file string) error {

	err = MyOwn(sshConn,DbUser,DbGroup,file)
	if err != nil{
		return err
	}
	err = MyMod(sshConn,file)
	if err != nil{
		return err
	}

	return nil
}

func CheckPath(sshConn *linux.MySSHConn,path,mkCmd string) error {
	pathExists, err := sshConn.PathExists(path)
	if err != nil {
		return err
	}
	if pathExists {
		err = sshConn.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	_,_,err = sshConn.ExecuteCommand(mkCmd)
	if err != nil{
		return err
	}

	return nil
}

func CreateSomeDir(sshConn *linux.MySSHConn,dataDir,port,DbUser,DbGroup string) error {
	// every time
	mkCmd := "sudo mkdir -p " + dataDir
	err := CheckPath(sshConn,dataDir,mkCmd)
	if err != nil{
		return err
	}

	mkLogCmd := "sudo mkdir -p " + "/mysqldata/mysql"+port+"/log/"
	err = CheckPath(sshConn,"/mysqldata/mysql"+port+"/log/",mkLogCmd)
	if err != nil{
		return err
	}
	mkBinLogCmd := "sudo mkdir -p " + "/mysqllog/mysql"+port+"/binlog/"
	err = CheckPath(sshConn,"/mysqllog/mysql"+port+"/binlog/",mkBinLogCmd)
	if err != nil{
		return err
	}
	mkRelayLogCmd := "sudo mkdir -p " + "/mysqllog/mysql"+port+"/relaylog/"
	err = CheckPath(sshConn,"/mysqllog/mysql"+port+"/relaylog/",mkRelayLogCmd)
	if err != nil{
		return err
	}
	err = MyOwn(sshConn,DbUser,DbGroup,"/mysqldata")
	if err != nil{
		return err
	}
	err = MyMod(sshConn,"/mysqldata")
	if err != nil{
		return err
	}
	err = MyOwn(sshConn,DbUser,DbGroup,"/mysqllog")
	if err != nil{
		return err
	}
	err = MyMod(sshConn,"/mysqllog")
	if err != nil{
		return err
	}
	return nil
}

func CopyBinaryCommand(sshConn *linux.MySSHConn) error{
	cpMysqlCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMysql,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMysqlCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpMydCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyd,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMydCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpSafeCmd := fmt.Sprintf(`cp "%s" "%s"`,cpSafe,binFile)
	_,_,err = sshConn.ExecuteCommand(cpSafeCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpMyMultiCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyMulti,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyMultiCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpMyDumpCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyDump,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyDumpCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpMyBinlogCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyBinlog,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyBinlogCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpMyCnfEditCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyCnfEdit,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyCnfEditCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpMyPriDefCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyPriDef,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyPriDefCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	cpMyAdmCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyAdm,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyAdmCmd)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func PrepareWork(sshConn *linux.MySSHConn,srcFile,destFile,exportStr,cmdSource,cmdGroup,cmdUser,
	DbUser,sLib,iLib string) error{
	err = sshConn.CopyToRemote(srcFile,destFile)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.CopyToRemote ok")

	err = CopyBinaryCommand(sshConn)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	cmdChangeProfile := fmt.Sprintf(`sed -i "%s" /etc/profile`,exportStr)
	_,_,err = sshConn.ExecuteCommand(cmdChangeProfile)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.cmdChangeProfile ok")

	_,_, err = sshConn.ExecuteCommand(cmdSource)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.cmdSource ok")

	err = InstallTool(sshConn,sLib,iLib)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.InstallTool ok")


	checkUser := "cat /etc/passwd|grep mysql"
	_,output,err := sshConn.ExecuteCommand(checkUser)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	if output != ""{
		fmt.Println("User exists!")
		return nil
	}
	_,_, err = sshConn.ExecuteCommand(cmdGroup)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.cmdGroup ok")
	_,_, err = sshConn.ExecuteCommand(cmdUser)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.cmdUser ok")

	return nil
}

func BasicWork(sshConn *linux.MySSHConn,srcFile,destFile,exportStr,cmdSource,cmdGroup,cmdUser,
	DbUser,DbGroup,file,sLib,iLib,dataDir,port string) error{

	err = CreateSomeDir(sshConn,dataDir,port,DbUser,DbGroup)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.CreateSomeDir ok")

	// copy mysql to /usr/local/
	err = PrepareWork(sshConn,srcFile,destFile,exportStr,
		cmdSource,cmdGroup,cmdUser,DbUser,sLib,iLib)
	if err != nil {
		return err
	}
	fmt.Println("BasicWork.PrepareWork ok")

	err = OwnAndMod(sshConn,DbUser,DbGroup,file)
	if err != nil{
		return err
	}
	fmt.Println("BasicWork.ChangeFile ok")
	return nil
}

func BasicWorkSlave(sshConn *linux.MySSHConn,srcFile,baseDir,DFileA,dataDir,DFileB,exportStr,cmdSource,
	cmdGroup,cmdUser,DbUser,DbGroup,sLib,iLib,port string) error{
	err = CreateSomeDir(sshConn,dataDir,port,DbUser,DbGroup)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.CreateSomeDir ok")

	err = PrepareWork(sshConn,srcFile,baseDir,exportStr,
		cmdSource,cmdGroup,cmdUser,DbUser,sLib,iLib)
	if err != nil {
		return err
	}
	fmt.Println("BasicWorkSlave.PrepareWork ok")

	err = sshConn.CopyToRemote(DFileA,dataDir)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	err = sshConn.CopyToRemote(DFileB,"/mysqldata/mysql"+port+"/log/")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	err = OwnAndMod(sshConn,DbUser,DbGroup,"/mysqldata")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	err = OwnAndMod(sshConn,DbUser,DbGroup,baseDir)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("BasicWorkSlave.CopyToRemote ok")

	_, _, err = sshConn.ExecuteCommand(startCommand+port)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("BasicWorkSlave.startCommand ok")
	return nil
}

func MysqlConnect(pwd,file,connectStr string) (db *sql.DB, err error){
	path := strings.Join([]string{uname,":",pwd,"@"+connectStr+"/",dbname,"?charset=utf8&multiStatements=true"},"")
	db,_ = sql.Open("mysql",path)

	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil{
		fmt.Println("open database fail",err.Error())
		return nil,err
	}
	fmt.Println("database connect success")

	sqlBytes,err := ioutil.ReadFile(file)
	if err != nil{
		fmt.Println(err.Error())
		return nil,err
	}
	sqlTable := string(sqlBytes)
	fmt.Println(sqlTable)
	_,err = db.Exec(sqlTable)
	if err != nil{
		return nil,err
	}
	return db,nil
}

func DbConnect(connectStr,file string) error {
	db,err := MysqlConnect(dbPasswd,file,connectStr)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()

	rows,err := db.Query("select * from `test`;")
	defer rows.Close()
	if err != nil{
		return err
	}
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id,&name)
		if err != nil{
			return err
		}
		fmt.Println(id,name)
	}
	return nil
}

func MyMulti(sshConn *linux.MySSHConn,connectStr,sqlFileMaster,dataDir,DFileA,logErrDir,DFileB,port string) error {
	err = DbConnect(connectStr,sqlFileMaster)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("MyMulti.DbConnect ok")
	_, _, err = sshConn.ExecuteCommand(stopCommand+port)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	time.Sleep(time.Duration(5)*time.Second)
	fmt.Println("MyMulti.stopCommand ok")

	err = sshConn.CopyFromRemote(dataDir,DFileA)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("MyMulti.CopyFromRemote ok")
	err = sshConn.CopyFromRemote(logErrDir,DFileB)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("MyMulti.CopyFromRemote ok")

	_, _, err = sshConn.ExecuteCommand(startCommand+port)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("MyMulti.startCommand ok")
	time.Sleep(time.Duration(5)*time.Second)

	return nil
}