package practicessh

import (
	"fmt"
	"compress/gzip"
	"strings"
	"io"
	"os"

	//"github.com/pkg/sftp"
	"github.com/romberli/go-util/linux"
	//"golang.org/x/crypto/ssh"
	_"github.com/go-sql-driver/mysql"
	"database/sql"
	"io/ioutil"
	"archive/tar"
	"bufio"
)

const (
	uname="root"
	dbpd="mysql"
	//ip="127.0.0.1"
	dbname="mysql"

	stopCommand = "mysqld_multi stop "
	startCommand = "mysqld_multi start "
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
	modCmd := fmt.Sprintf(`sudo chmod -R g+rwx "%s"`,file)
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

func ChangeMysqlFile(filename,dataDir string) error {
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

func ChangeFile(sshConn *linux.MySSHConn,DbUser,DbGroup,file string) error {

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

func PrepareWork(sshConn *linux.MySSHConn,fileS,fileD,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
	dataDir,port,DbUser,DbGroup,sLib,iLib string) error{
	err = sshConn.CopyToRemote(fileS,fileD)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.CopyToRemote ok")
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

	//checkUser:="cat /etc/passwd|grep " + DbUser
	//_, output, err := sshConn.ExecuteCommand(checkUser)
	//if err != nil{
	//	fmt.Println(err.Error())
	//	return err
	//}
	//if output != ""{
	//	cmdUserDel := "userdel "+ DbUser
	//	_,_, err = sshConn.ExecuteCommand(cmdUserDel)
	//	if err != nil{
	//		fmt.Println(err.Error())
	//		return err
	//	}
	//}
	cmdUserDel := "userdel "+ DbUser
	_,_, err = sshConn.ExecuteCommand(cmdUserDel)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.checkUser ok")

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
	err = InstallTool(sshConn,sLib,iLib)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.InstallTool ok")
	err = CreateSomeDir(sshConn,dataDir,port,DbUser,DbGroup)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("PrepareWork.CreateSomeDir ok")
	return nil
}

func BasicWork(sshConn *linux.MySSHConn,fileNameSource,fileNameDest,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
	DbUser,DbGroup,file,sLib,iLib,dataDir,port string) error{

	err = PrepareWork(sshConn,fileNameSource,fileNameDest,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
		dataDir,port,DbUser,DbGroup,sLib,iLib)
	if err != nil {
		return err
	}
	fmt.Println("BasicWork.PrepareWork ok")
	err = ChangeFile(sshConn,DbUser,DbGroup,file)
	if err != nil{
		return err
	}
	fmt.Println("BasicWork.ChangeFile ok")
	return nil
}

func BasicWorkSlave(sshConn *linux.MySSHConn,DFileA,file,DFileB,dataDir,cmdChangeProfile,cmdSource,
	cmdGroup,cmdUser,DbUser,DbGroup,sLib,iLib,port string) error{

	err = PrepareWork(sshConn,DFileA, file,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
		dataDir,port,DbUser,DbGroup,sLib,iLib)
	if err != nil {
		return err
	}
	fmt.Println("BasicWorkSlave.PrepareWork ok")

	err = sshConn.CopyToRemote(DFileB, dataDir)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}

	err = MyOwn(sshConn,DbUser,DbGroup,dataDir)
	if err != nil{
		return err
	}
	err = MyMod(sshConn,dataDir)
	if err != nil{
		return err
	}

	err = ChangeFile(sshConn,DbUser,DbGroup,file)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("BasicWorkSlave.ChangeFile ok")
	_, _, err = sshConn.ExecuteCommand(startCommand+port)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("BasicWorkSlave.startCommand ok")
	return nil
}

func MysqlConnect(ip,port,pwd,file string) (db *sql.DB, err error){
	path := strings.Join([]string{uname,":",pwd,"@tcp(",ip,":",port,")/",dbname,"?charset=utf8&multiStatements=true"},"")
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

func ResetPasswd(ip,port,pwd,file string) error {
	db,err := MysqlConnect(ip,port,pwd,file)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()
	return nil
}

func DbConnect(ip,port,file string) error {
	db,err := MysqlConnect(ip,port,dbpd,file)
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

func MyMulti(sshConn *linux.MySSHConn,sqlFileMaster,SFileA,DFileA,SFileB,DFileB,ip,port string) error {
	err = DbConnect(ip,port,sqlFileMaster)
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
	fmt.Println("MyMulti.stopCommand ok")

	err = sshConn.CopyFromRemote(SFileA, DFileA)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("MyMulti.CopyFromRemote ok")

	err = sshConn.CopyFromRemote(SFileB, DFileB)
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

	//time.Sleep(time.Duration(3)*time.Second)
	//
	//_, _, err = conn.ExecuteCommand(reportCommand)
	//if err != nil{
	//	return err
	//}

	return nil
}