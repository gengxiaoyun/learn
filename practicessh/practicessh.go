package practicessh

import (
	"fmt"
	"compress/gzip"
	"strings"
	"io"
	"os"
	"github.com/romberli/go-util/linux"
	_"github.com/go-sql-driver/mysql"
	"archive/tar"
	"bufio"
	"time"
)

const (
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

	setPdSql = "/usr/local/setpassword.sql"
	setMasterSql = "/usr/local/master.sql"
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
		return err
	}
	defer f.Close()

	out,err := os.OpenFile("./mysql.server", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil{
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
			return err
		}
	}
	err = os.Rename("./mysql.server","./mysql")
	if err != nil{
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

func CheckUser(sshConn *linux.MySSHConn,DbUser,cmdGroup,cmdUser string) error{
	checkUser := "cat /etc/passwd|grep " + DbUser
	_,output,err := sshConn.ExecuteCommand(checkUser)
	if err != nil{
		return err
	}
	if output == ""{
		_,_, err = sshConn.ExecuteCommand(cmdGroup)
		if err != nil{
			return err
		}
		_,_, err = sshConn.ExecuteCommand(cmdUser)
		if err != nil{
			return err
		}
		return nil
	}
	fmt.Println("User exists!")
	return nil
}

func CreateSomeDir(sshConn *linux.MySSHConn,dataDir,port,DbUser,
	DbGroup,cmdGroup,cmdUser string) error {
	// every time
	err = CheckUser(sshConn,DbUser,cmdGroup,cmdUser)
	if err != nil{
		return err
	}
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

	err = OwnAndMod(sshConn,DbUser,DbGroup,"/mysqldata")
	if err != nil{
		return err
	}
	err = OwnAndMod(sshConn,DbUser,DbGroup,"/mysqllog")
	if err != nil{
		return err
	}

	return nil
}

func CopyBinaryCommand(sshConn *linux.MySSHConn) error{
	cpMysqlCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMysql,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMysqlCmd)
	if err != nil{
		return err
	}
	cpMydCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyd,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMydCmd)
	if err != nil{
		return err
	}
	cpSafeCmd := fmt.Sprintf(`cp "%s" "%s"`,cpSafe,binFile)
	_,_,err = sshConn.ExecuteCommand(cpSafeCmd)
	if err != nil{
		return err
	}
	cpMyMultiCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyMulti,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyMultiCmd)
	if err != nil{
		return err
	}
	cpMyDumpCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyDump,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyDumpCmd)
	if err != nil{
		return err
	}
	cpMyBinlogCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyBinlog,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyBinlogCmd)
	if err != nil{
		return err
	}
	cpMyCnfEditCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyCnfEdit,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyCnfEditCmd)
	if err != nil{
		return err
	}
	cpMyPriDefCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyPriDef,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyPriDefCmd)
	if err != nil{
		return err
	}
	cpMyAdmCmd := fmt.Sprintf(`cp "%s" "%s"`,cpMyAdm,binFile)
	_,_,err = sshConn.ExecuteCommand(cpMyAdmCmd)
	if err != nil{
		return err
	}

	return nil
}

// copy mysql to /usr/local/, add group and user
func PrepareWork(sshConn *linux.MySSHConn,srcFile,destFile,exportStr,cmdSource,
	sLib,iLib string) error{
	err = sshConn.CopyToRemote(srcFile,destFile)
	if err != nil{
		return err
	}
	fmt.Println("CopyToRemote")

	err = CopyBinaryCommand(sshConn)
	if err != nil{
		return err
	}
	fmt.Println("CopyBinaryCommand")

	cmdChangeProfile := fmt.Sprintf(`sed -i "%s" /etc/profile`,exportStr)
	_,_,err = sshConn.ExecuteCommand(cmdChangeProfile)
	if err != nil{
		return err
	}
	fmt.Println("cmdChangeProfile")

	_,_, err = sshConn.ExecuteCommand(cmdSource)
	if err != nil{
		return err
	}
	fmt.Println("cmdSource")

	err = InstallTool(sshConn,sLib,iLib)
	if err != nil{
		return err
	}
	fmt.Println("InstallTool")

	return nil
}

func BasicWork(sshConn *linux.MySSHConn,srcFile,destFile,exportStr,cmdSource,cmdGroup,cmdUser,
	DbUser,DbGroup,file,sLib,iLib,dataDir,port string) error{

	err = CreateSomeDir(sshConn,dataDir,port,DbUser,DbGroup,cmdGroup,cmdUser)
	if err != nil{
		return err
	}

	// copy mysql to /usr/local/, add group and user
	err = PrepareWork(sshConn,srcFile,destFile,exportStr,
		cmdSource,sLib,iLib)
	if err != nil {
		return err
	}

	err = OwnAndMod(sshConn,DbUser,DbGroup,file)
	if err != nil{
		return err
	}

	return nil
}

func BasicWorkSlave(sshConn *linux.MySSHConn,srcFile,baseDir,DFileA,dataDir,DFileB,exportStr,cmdSource,
	cmdGroup,cmdUser,DbUser,DbGroup,sLib,iLib,port string) error{
	err = CreateSomeDir(sshConn,dataDir,port,DbUser,DbGroup,cmdGroup,cmdUser)
	if err != nil{
		return err
	}

	err = PrepareWork(sshConn,srcFile,baseDir,exportStr,
		cmdSource,sLib,iLib)
	if err != nil {
		return err
	}

	err = sshConn.CopyToRemote(DFileA,dataDir)
	if err != nil{
		return err
	}
	err = sshConn.CopyToRemote(DFileB,"/mysqldata/mysql"+port+"/log/")
	if err != nil{
		return err
	}

	err = OwnAndMod(sshConn,DbUser,DbGroup,"/mysqldata")
	if err != nil{
		return err
	}

	err = OwnAndMod(sshConn,DbUser,DbGroup,baseDir)
	if err != nil{
		return err
	}

	_, _, err = sshConn.ExecuteCommand(startCommand+port)
	if err != nil{
		return err
	}
	time.Sleep(time.Duration(20)*time.Second)

	return nil
}

func MyMulti(sshConn *linux.MySSHConn,dataDir,DFileA,logErrDir,DFileB,port string) error {
	setPdCommand := "mysql -uroot < " + setPdSql
	_,_,err = sshConn.ExecuteCommand(setPdCommand)
	if err != nil{
		return err
	}
	setMasterCommand := "mysql -uroot -pmysql < " + setMasterSql
	_,_,err = sshConn.ExecuteCommand(setMasterCommand)
	if err != nil{
		return err
	}

	_, _, err = sshConn.ExecuteCommand(stopCommand+port)
	if err != nil{
		return err
	}
	time.Sleep(time.Duration(10)*time.Second)

	err = sshConn.CopyFromRemote(dataDir,DFileA)
	if err != nil{
		return err
	}
	err = sshConn.CopyFromRemote(logErrDir,DFileB)
	if err != nil{
		return err
	}

	_, _, err = sshConn.ExecuteCommand(startCommand+port)
	if err != nil{
		return err
	}
	time.Sleep(time.Duration(15)*time.Second)

	return nil
}