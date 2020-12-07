package practicessh

import (
	"time"
	"bytes"
	"fmt"
	"net"
	"errors"
	"compress/gzip"
	"path/filepath"
	"strings"
	"io"
	"os"

	"github.com/pkg/sftp"
	//"github.com/romberli/go-util/linux"
	"golang.org/x/crypto/ssh"
	_"github.com/go-sql-driver/mysql"
	"bufio"
	"database/sql"
	"io/ioutil"
	"archive/tar"
)

const (
	MinUInt             uint        = 0
	EmptyString                     = ""
	DefaultExecFileMode os.FileMode = 0755
	LsCommand               = "ls"

	DefaultSplitStr           = "\n"
	DefaultStringZeroValue    = ""
	DefaultFailedReturnValue  = 1
	DefaultSSHTimeout         = 10 * time.Second
	DefaultSSHPortNum         = 22
	DefaultSSHUserName        = "root"
	DefaultSSHUserPass        = "root"
	DefaultByteBufferSize     = 1024 * 1024 // 1MB

	uname="root"
	pwd="mysql"
	ip="127.0.0.1"
	dbname="mysql"

	stopCommand = "mysqld_multi stop "
	startCommand = "mysqld_multi start "
)

//var conn *linux.MySSHConn

type MyConn struct {
	HostIp   string
	PortNum  int
	UserName string
	UserPass string
}

func NewMyConn(hostIP string, portNum int, userName string, userPass string) (conn *MyConn) {
	return &MyConn{
		hostIP,
		portNum,
		userName,
		userPass,
	}
}

func NewMyConnWithDefaultValue(hostIP string) (conn *MyConn) {
	return &MyConn{
		hostIP,
		DefaultSSHPortNum,
		DefaultSSHUserName,
		DefaultSSHUserPass,
	}
}

type MySSHConn struct {
	MyConn
	SSHClient *ssh.Client
	*sftp.Client
}


// create file
func (conn *MySSHConn) CreateFile(name string) (*os.File,error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name,"/")]),0755)
	if err!=nil{
		return nil,err
	}
	return os.Create(name)
}

// unzip
func (conn *MySSHConn) UnTar(srcFile,destFile string) error {
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
			file,err := conn.CreateFile(filename)
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


func (conn *MySSHConn) MyOwn(DbUser,DbGroup,file string) error {
	ownCmd := fmt.Sprintf(`sudo chown -R "%s"."%s" "%s"`,DbUser,DbGroup,file)
	_,_,err := conn.ExecuteCommand(ownCmd)
	if err != nil{
		return err
	}
	return nil
}

func (conn *MySSHConn) MyMod(file string) error {
	modCmd := fmt.Sprintf(`sudo chmod -R g+rwx "%s"`,file)
	_,_,err := conn.ExecuteCommand(modCmd)
	if err != nil{
		return err
	}
	return nil
}

// install libaio
func (conn *MySSHConn) InstallTool(sLib,iLib string) error {
	_,_,err := conn.ExecuteCommand(sLib)
	if err != nil{
		return err
	}
	_,_,err = conn.ExecuteCommand(iLib)
	if err != nil{
		return err
	}
	return nil
}

func (conn *MySSHConn) ChangeMysqlFile(filename,dataDir string) error {
	f,err := os.Open(filename)
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	defer f.Close()

	out,err := os.OpenFile(filename+".mdf", os.O_RDWR, 0777)
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
	return nil
}

func (conn *MySSHConn) ChangeFile(DbUser,DbGroup,file,srcServer,destServer,strCommand,srcMdf,dataDir string) error {

	err = conn.MyOwn(DbUser,DbGroup,file)
	if err != nil{
		return err
	}
	err = conn.MyMod(file)
	if err != nil{
		return err
	}

	// mysql.server->mysql
	cpCmd := fmt.Sprintf(`cp "%s" "%s"`, srcServer, destServer)
	_,_,err := conn.ExecuteCommand(cpCmd)
	if err != nil{
		return err
	}
	// set basedir and datadir
	//err = MyMod(destServer,dir)
	//if err != nil{
	//	return err
	//}

	_,_,err = conn.ExecuteCommand(strCommand)
	if err != nil{
		return err
	}

	err = conn.MyMod(srcMdf)
	if err != nil{
		return err
	}

	err = conn.ChangeMysqlFile(destServer,dataDir)
	if err != nil{
		return err
	}

	mvCmd := fmt.Sprintf(`sudo mv "%s" "%s"`, srcMdf, destServer)
	_,_,err = conn.ExecuteCommand(mvCmd)
	if err != nil{
		return err
	}

	return nil
}

func (conn *MySSHConn) CreateSomeDir(dataDir,port,DbUser,DbGroup string) error {
	mkCmd := "sudo mkdir -p " + dataDir
	_,_,err = conn.ExecuteCommand(mkCmd)
	if err != nil{
		return err
	}
	mkLogCmd := "sudo mkdir " + "/mysqldata/mysql"+port+"/log/"
	_,_,err = conn.ExecuteCommand(mkLogCmd)
	if err != nil{
		return err
	}
	mkBinLogCmd := "sudo mkdir -p " + "/mysqllog/mysql"+port+"/binlog/"
	_,_,err = conn.ExecuteCommand(mkBinLogCmd)
	if err != nil{
		return err
	}
	mkRelayLogCmd := "sudo mkdir " + "/mysqllog/mysql"+port+"/relaylog/"
	_,_,err = conn.ExecuteCommand(mkRelayLogCmd)
	if err != nil{
		return err
	}
	err = conn.MyOwn(DbUser,DbGroup,"/mysqldata")
	if err != nil{
		return err
	}
	err = conn.MyMod("/mysqldata")
	if err != nil{
		return err
	}
	err = conn.MyOwn(DbUser,DbGroup,"/mysqllog")
	if err != nil{
		return err
	}
	err = conn.MyMod("/mysqllog")
	if err != nil{
		return err
	}
	return nil
}

func (conn *MySSHConn) PrepareWork(fileS,fileD,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
	dataDir,port,DbUser,DbGroup,sLib,iLib string) error{
	err = conn.CopyToRemote(fileS,fileD)
	if err != nil{
		return err
	}
	_, _, err = conn.ExecuteCommand(cmdChangeProfile)
	if err != nil{
		return err
	}
	_, _, err = conn.ExecuteCommand(cmdSource)
	if err != nil{
		return err
	}
	_, _, err = conn.ExecuteCommand(cmdGroup)
	if err != nil{
		return err
	}
	_, _, err = conn.ExecuteCommand(cmdUser)
	if err != nil{
		return err
	}

	err = conn.InstallTool(sLib,iLib)
	if err != nil{
		return err
	}

	err = conn.CreateSomeDir(dataDir,port,DbUser,DbGroup)
	if err != nil{
		return err
	}
	return nil
}

func (conn *MySSHConn) BasicWork(fileNameSource, fileNameDest,srcFile,destFile,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
	DbUser,DbGroup,file,sLib,iLib,srcServer,destServer,strCommand,srcMdf,dataDir,port string) error{

	err = conn.PrepareWork(fileNameSource, fileNameDest,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
		dataDir,port,DbUser,DbGroup,sLib,iLib)
	if err != nil {
		return err
	}

	err = conn.UnTar(srcFile,destFile)
	if err != nil {
		return err
	}

	err = conn.ChangeFile(DbUser,DbGroup,file,srcServer,destServer,
		strCommand,srcMdf,dataDir)
	if err != nil{
		return err
	}

	return nil
}

func (conn *MySSHConn) BasicWorkSlave(DFileA,DFileB,srcFile,destFile,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
DbUser,DbGroup,file,sLib,iLib,srcServer,destServer,strCommand,srcMdf,dataDir,port string) error{

	err = conn.PrepareWork(DFileA, file,cmdChangeProfile,cmdSource,cmdGroup,cmdUser,
		dataDir,port,DbUser,DbGroup,sLib,iLib)
	if err != nil {
		return err
	}

	err = conn.CopyToRemote(DFileB, dataDir)
	if err != nil{
		fmt.Println(err.Error())
	}

	err = conn.MyOwn(DbUser,DbGroup,dataDir)
	if err != nil{
		return err
	}
	err = conn.MyMod(dataDir)
	if err != nil{
		return err
	}

	err = conn.ChangeFile(DbUser,DbGroup,file,srcServer,destServer,
		strCommand,srcMdf,dataDir)
	if err != nil{
		return err
	}

	_, _, err = conn.ExecuteCommand(startCommand+port)
	if err != nil{
		return err
	}

	return nil
}

func (conn *MySSHConn) DbConnect(port,file string) error {
	path := strings.Join([]string{uname,":",pwd,"@tcp(",ip,":",port,")/",dbname,"?charset=utf8&multiStatements=true"},"")
	db,_ := sql.Open("mysql",path)
	defer db.Close()
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(10)
	err := db.Ping()
	if err != nil{
		fmt.Println("open database fail",err.Error())
		return err
	}
	fmt.Println("database connect success")
	sqlBytes,err := ioutil.ReadFile(file)
	if err != nil{
		return err
	}
	sqlTable := string(sqlBytes)
	fmt.Println(sqlTable)
	_,err = db.Exec(sqlTable)
	if err != nil{
		return err
	}

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

func (conn *MySSHConn) MyMulti(sqlFileMaster,sqlFileTest,SFileA,DFileA,SFileB,DFileB,port string) error {
	err = conn.DbConnect(port,sqlFileMaster)
	if err != nil{
		fmt.Println(err.Error())
	}

	_, _, err = conn.ExecuteCommand(stopCommand+port)
	if err != nil{
		return err
	}

	err = conn.CopyFromRemote(SFileA, DFileA)
	if err != nil{
		return err
	}
	err = conn.CopyFromRemote(SFileB, DFileB)
	if err != nil{
		return err
	}
	_, _, err = conn.ExecuteCommand(startCommand+port)
	if err != nil{
		return err
	}

	//time.Sleep(time.Duration(3)*time.Second)
	//
	//_, _, err = conn.ExecuteCommand(reportCommand)
	//if err != nil{
	//	return err
	//}

	return nil
}


func NewMySSHConn(hostIP string, portNum int, userName, userPass string) (*MySSHConn, error) {
	return NewMySSHConnWithOptionalArgs(hostIP, portNum, userName, userPass)
}

func NewMySSHConnWithOptionalArgs(hostIP string, in ...interface{}) (sshConn *MySSHConn, err error) {
	var (
		myConn       *MyConn
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
	)

	argLen := len(in)
	switch argLen {
	case 0:
		hostIP = strings.TrimSpace(hostIP)
		if hostIP == "" {
			return nil, errors.New("host ip could not be empty")
		}

		myConn = NewMyConnWithDefaultValue(hostIP)
	case 3:
		var (
			portNumValue  int
			userNameValue string
			userPassValue string
		)

		portNum := in[0]
		userName := in[1]
		userPass := in[2]

		switch portNum.(type) {
		case nil:
			portNumValue = DefaultSSHPortNum
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			portNumValue = portNum.(int)
		default:
			return nil, errors.New(
				fmt.Sprintf("port number must be integer type instead of %s",
					reflect.TypeOf(portNum).Name()))
		}

		switch userName.(type) {
		case nil:
			userNameValue = DefaultSSHUserName
		case string:
			userNameValue = strings.TrimSpace(userName.(string))
			if userNameValue == "" {
				userNameValue = DefaultSSHUserName
			}
		default:
			return nil, errors.New(
				fmt.Sprintf("user name must be string type instead of %s",
					reflect.TypeOf(portNum).Name()))
		}

		switch userPass.(type) {
		case nil:
			userPassValue = DefaultSSHUserPass
		case string:
			userPassValue = strings.TrimSpace(userPass.(string))
			if userPassValue == "" {
				userPassValue = DefaultSSHUserPass
			}
		default:
			return nil, errors.New(
				fmt.Sprintf("user password must be string type instead of %s",
					reflect.TypeOf(portNum).Name()))
		}

		myConn = NewMyConn(hostIP, portNumValue, userNameValue, userPassValue)
	default:
		return nil, errors.New(fmt.Sprintf("optional argument number must be 0 or 3 instead of %d", argLen))
	}

	// get auth method
	auth = append(auth, ssh.Password(myConn.UserPass))

	hostKeyCallBack := func(host string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            myConn.UserName,
		Auth:            auth,
		Timeout:         DefaultSSHTimeout,
		HostKeyCallback: hostKeyCallBack,
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", myConn.HostIp, myConn.PortNum)
	sshClient, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}

	// create sftp client
	sftpClient, err = sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	sshConn = &MySSHConn{
		*myConn,
		sshClient,
		sftpClient,
	}

	return sshConn, nil
}

// Close closes connections with the remote host
func (conn *MySSHConn) Close() (err error) {
	err = conn.Client.Close()
	if err != nil {
		return err
	}

	return conn.SSHClient.Close()
}

// ExecuteCommand executes shell command on the remote host
func (conn *MySSHConn) ExecuteCommand(cmd string) (result int, output string, err error) {
	var (
		stdOutBuffer bytes.Buffer
		stdErrBuffer bytes.Buffer
	)

	// create ssh session
	sshSession, err := conn.SSHClient.NewSession()
	if err != nil {
		return DefaultFailedReturnValue, DefaultStringZeroValue, err
	}
	defer func() { _ = sshSession.Close() }()

	sshSession.Stdout = &stdOutBuffer
	sshSession.Stderr = &stdErrBuffer

	// run command
	err = sshSession.Run(cmd)
	if err != nil {
		result = DefaultFailedReturnValue
		if stdErrBuffer.String() != EmptyString {
			err = errors.New(fmt.Sprintf("%s\n%s", err.Error(), stdErrBuffer.String()))
		}
	}

	output = stdOutBuffer.String() + stdErrBuffer.String()
	return result, output, err
}

// PathExists returns if given path exists
func (conn *MySSHConn) PathExists(path string) (bool, error) {
	_, err := conn.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// IsDir returns if given path on the remote host is a directory or not
func (conn *MySSHConn) IsDir(path string) (isDir bool, err error) {
	path = strings.TrimSpace(path)

	info, err := conn.Stat(path)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

// ListPath returns subdirectories and files of given path on the remote host, it returns a slice of sub paths
func (conn *MySSHConn) ListPath(path string) (subPathList []string, err error) {
	path = strings.TrimSpace(path)

	cmd := fmt.Sprintf("%s %s", LsCommand, path)
	_, subPathStr, err := conn.ExecuteCommand(cmd)
	if err != nil {
		return nil, err
	}

	subPathStr = strings.TrimSpace(subPathStr)
	if subPathStr != EmptyString {
		subPathList = strings.Split(subPathStr, DefaultSplitStr)
	}

	return subPathList, nil
}

// ReadDir returns subdirectories and files of given directory on the remote host, it returns a slice of os.FileInfo
func (conn *MySSHConn) ReadDir(dirName string) (fileInfoList []os.FileInfo, err error) {
	dirName = strings.TrimSpace(dirName)

	isDir, err := conn.IsDir(dirName)
	if err != nil {
		return nil, err
	}
	if !isDir {
		return nil, errors.New(fmt.Sprintf("it's NOT a directory. dir name: %s", dirName))
	}

	subPathList, err := conn.ListPath(dirName)
	if err != nil {
		return nil, err
	}
	for _, subPath := range subPathList {
		if subPath != EmptyString {
			fileNameAbs := filepath.Join(dirName, subPath)
			fileInfo, err := conn.Stat(fileNameAbs)
			if err != nil {
				return nil, err
			}

			fileInfoList = append(fileInfoList, fileInfo)
		}
	}

	return fileInfoList, err
}

// RemoveAll removes given path on the remote host, it will act like shell command "rm -rf $path",
// except that it will raise an error when something goes wrong.
func (conn *MySSHConn) RemoveAll(path string) (err error) {
	path = strings.TrimSpace(path)

	isDir, err := conn.IsDir(path)
	if err != nil {
		return err
	}

	if isDir {
		isEmpty, err := conn.IsEmptyDir(path)
		if err != nil {
			return err
		}

		if !isEmpty {
			subPathList, err := conn.ListPath(path)
			for _, subPath := range subPathList {
				subPathAbs := filepath.Join(path, subPath)
				err = conn.RemoveAll(subPathAbs)
				if err != nil {
					return err
				}
			}
		}

		err = conn.RemoveDirectory(path)
		if err != nil {
			return err
		}
	} else {
		err = conn.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}


// IsEmptyDir returns if  given directory is empty or not
func (conn *MySSHConn) IsEmptyDir(dirName string) (isEmpty bool, err error) {
	dirName = strings.TrimSpace(dirName)

	fileInfoList, err := conn.ReadDir(dirName)
	if err != nil {
		return false, err
	}

	if fileInfoList == nil {
		isEmpty = true
	}

	return isEmpty, nil
}

// CopyFile copy file content from source to destination, it doesn't care about which one local or remote
func (conn *MySSHConn) CopyFile(fileSource io.Reader, fileDest io.Writer, bufferSize int) (err error) {
	var n int

	if bufferSize <= int(MinUInt) {
		bufferSize = DefaultByteBufferSize
	}

	buf := make([]byte, bufferSize)

	for {
		n, err = fileSource.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		_, err = fileDest.Write(buf[int(MinUInt):n])
		if err != nil {
			return err
		}
	}

	return nil
}

// CopySingleFileFromRemote copies one single file from remote to local
func (conn *MySSHConn) CopySingleFileFromRemote(fileNameSource string, fileNameDest string) (err error) {
	var (
		fileDest   *os.File
		fileSource *sftp.File
	)

	fileNameSource = strings.TrimSpace(fileNameSource)
	fileNameDest = strings.TrimSpace(fileNameDest)

	if fileNameDest == EmptyString {
		fileNameDest = fileNameSource
	}

	// check if source path is a directory
	isDir, err := conn.IsDir(fileNameSource)
	if err != nil {
		return err
	}
	if isDir {
		return errors.New(fmt.Sprintf("it's NOT a file. file name: %s", fileNameSource))
	}

	// check if parent path of destination exists
	fileNameDestParent := filepath.Dir(fileNameDest)
	pathExists, err := conn.PathExists(fileNameDestParent)
	if err != nil {
		return nil
	}
	if !pathExists {
		return errors.New(fmt.Sprintf("parent path of destination does NOT exsists. path: %s", fileNameDest))
	}

	// check if destination path is a directory
	pathExists, err = conn.PathExists(fileNameDest)
	if err != nil {
		return err
	}
	if pathExists {
		isDir, err = conn.IsDir(fileNameDest)
		if err != nil {
			return err
		}
		if isDir {
			fileNameSourceBase := filepath.Base(fileNameSource)
			fileNameDest = filepath.Join(fileNameDest, fileNameSourceBase)
		}
	}

	fileSource, err = conn.Open(fileNameSource)
	if err != nil {
		return err
	}
	defer func() { _ = fileSource.Close() }()

	fileDest, err = os.Create(fileNameDest)
	if err != nil {
		return err
	}
	defer func() { _ = fileDest.Close() }()

	err = conn.CopyFile(fileSource, fileDest, DefaultByteBufferSize)
	if err != nil {
		return err
	}

	return nil
}


// CopySingleFileToRemote copies one single file from local to remote
func (conn *MySSHConn) CopySingleFileToRemote(fileNameSource string, fileNameDest string) (err error) {
	var (
		fileSource *os.File
		fileDest   *sftp.File
	)

	fileNameSource = strings.TrimSpace(fileNameSource)
	fileNameDest = strings.TrimSpace(fileNameDest)

	if fileNameDest == EmptyString {
		fileNameDest = fileNameSource
	}

	// check if source path is a directory
	isDir, err := conn.IsDir(fileNameSource)
	if err != nil {
		return err
	}
	if isDir {
		return errors.New(fmt.Sprintf("it's NOT a file. file name: %s", fileNameSource))
	}

	// check if parent path of destination exists
	fileNameDestParent := filepath.Dir(fileNameDest)
	pathExists, err := conn.PathExists(fileNameDestParent)
	if err != nil {
		return nil
	}
	if !pathExists {
		return errors.New(fmt.Sprintf("parent path of destination does NOT exsists. path: %s", fileNameDest))
	}

	// check if destination path is a directory
	pathExists, err = conn.PathExists(fileNameDest)
	if err != nil {
		return err
	}
	if pathExists {
		isDir, err = conn.IsDir(fileNameDest)
		if err != nil {
			return err
		}
		if isDir {
			fileNameSourceBase := filepath.Base(fileNameSource)
			fileNameDest = filepath.Join(fileNameDest, fileNameSourceBase)
		}
	}

	fileSource, err = os.Open(fileNameSource)
	if err != nil {
		return err
	}
	defer func() { _ = fileSource.Close() }()

	fileDest, err = conn.Create(fileNameDest)
	if err != nil {
		return err
	}
	defer func() { _ = fileDest.Close() }()

	err = conn.CopyFile(fileSource, fileDest, DefaultByteBufferSize)
	if err != nil {
		return err
	}

	return nil
}

// CopyFileListFromRemote copies given files from remote to local
//func (conn *MySSHConn) CopyFileListFromRemote(fileListSource []string, FileDirDest string) (err error) {
//	FileDirDest = strings.TrimSpace(FileDirDest)
//	if FileDirDest == DefaultStringZeroValue {
//		return errors.New("file destination directory should NOT an empty string")
//	}
//
//	pathExists, err := conn.PathExists(FileDirDest)
//	if err != nil {
//		return err
//	}
//
//	if !pathExists {
//		_, err = os.Create(FileDirDest)
//		if err != nil {
//			return err
//		}
//	}
//
//	for _, fileNameSource := range fileListSource {
//		fileNameSource = strings.TrimSpace(fileNameSource)
//		fileNameDest := path.Base(fileNameSource)
//
//		err = conn.CopySingleFileFromRemote(fileNameSource, path.Join(FileDirDest, fileNameDest))
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}


// CopyFileListFromRemoteWithNewName copies file from remote to local,
// it copies file contents and rename files to given file names
//func (conn *MySSHConn) CopyFileListFromRemoteWithNewName(fileListSource []string, FileListDest []string) (err error) {
//	if len(fileListSource) != len(FileListDest) {
//		return errors.New("the length of source and destination list MUST be exactly same")
//	}
//
//	for i, fileNameSource := range fileListSource {
//		fileNameSource = strings.TrimSpace(fileNameSource)
//
//		fileNameDest := FileListDest[i]
//		fileNameDest = strings.TrimSpace(fileNameDest)
//
//		if fileNameDest == DefaultStringZeroValue {
//			return errors.New("destination file name should not be an empty string")
//		}
//
//		err = conn.CopySingleFileFromRemote(fileNameSource, fileNameDest)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

// GetPathDirMapRemote reads all subdirectories and files of given directory on the remote host
// and calculate the relative path of rootPath,
// then map the absolute path of subdirectory names and file names as keys, relative paths as values to fileDirMap
func (conn *MySSHConn) GetPathDirMapRemote(fileDirMap map[string]string, dirName, rootPath string) (err error) {
	dirName = strings.TrimSpace(dirName)
	rootPath = strings.TrimSpace(rootPath)

	fileInfoList, err := conn.ReadDir(dirName)
	if err != nil {
		return err
	}

	if fileInfoList == nil {
		// it's an empty directory
		fileDirMap[dirName] = EmptyString
	}

	for _, fileInfo := range fileInfoList {
		fileName := fileInfo.Name()
		fileNameAbs := filepath.Join(dirName, fileName)

		if fileInfo.IsDir() {
			// call recursively
			err = conn.GetPathDirMapRemote(fileDirMap, fileNameAbs, rootPath)
			if err != nil {
				return err
			}
		} else {
			// get relative path with root path
			fileNameRel, err := filepath.Rel(rootPath, fileNameAbs)
			if err != nil {
				return err
			}

			fileDirMap[fileNameAbs] = fileNameRel
		}
	}

	return nil
}


// CopyDirFromRemote copies a directory with all subdirectories and files from remote to local
func (conn *MySSHConn) CopyDirFromRemote(dirNameSource, dirNameDest string) (err error) {
	dirNameSource = strings.TrimSpace(dirNameSource)
	dirNameDest = strings.TrimSpace(dirNameDest)

	// check if source path is a directory
	isDir, err := conn.IsDir(dirNameSource)
	if err != nil {
		return err
	}
	if !isDir {
		return errors.New(fmt.Sprintf("it's NOT a directory. dir name: %s", dirNameSource))
	}

	// check if parent path of destination exists
	dirNameDestParent := filepath.Dir(dirNameDest)
	pathExists, err := conn.PathExists(dirNameDestParent)
	if err != nil {
		return nil
	}
	if !pathExists {
		return errors.New(fmt.Sprintf("parent path of destination does NOT exsists. path: %s", dirNameDest))
	}

	pathSourceBase := filepath.Base(dirNameSource)
	pathDestBase := filepath.Base(dirNameDest)

	// get new destination path to act like shell command "scp -r"
	if pathSourceBase != pathDestBase {
		pathExists, err := conn.PathExists(dirNameDest)
		if err != nil {
			return err
		}
		if pathExists {
			isDir, err := conn.IsDir(dirNameDest)
			if err != nil {
				return err
			}
			if !isDir {
				return errors.New(fmt.Sprintf("it's NOT a directory. dir name: %s", dirNameDest))
			}

			dirNameDest = filepath.Join(dirNameDest, pathSourceBase)
		}
	}

	pathDirMap := make(map[string]string)
	// get map of source path and relative path with destination directory
	err = conn.GetPathDirMapRemote(pathDirMap, dirNameSource, dirNameSource)
	if err != nil {
		return err
	}

	for pathName, relDir := range pathDirMap {
		if relDir == EmptyString {
			// it's an empty directory, we just need to create it
			relDirSource, err := filepath.Rel(dirNameSource, pathName)
			if err != nil {
				return err
			}

			dirDestAbs := filepath.Join(dirNameDest, relDirSource)
			err = os.MkdirAll(dirDestAbs, DefaultExecFileMode)
			continue
		}

		relDir = filepath.Dir(relDir)
		DirDestAbs := filepath.Join(dirNameDest, relDir)
		err = os.MkdirAll(DirDestAbs, DefaultExecFileMode)
		if err != nil {
			return err
		}

		fileNameDest := GetFileNameDest(pathName, DirDestAbs)
		// copy file from remote
		err = conn.CopySingleFileFromRemote(pathName, fileNameDest)
		if err != nil {
			return err
		}
	}

	return nil
}


// CopyDirFromRemote copies a directory with all subdirectories and files from local to remote
func (conn *MySSHConn) CopyDirToRemote(dirNameSource, dirNameDest string) (err error) {
	dirNameSource = strings.TrimSpace(dirNameSource)
	dirNameDest = strings.TrimSpace(dirNameDest)

	// check if source path is a directory
	isDir, err := conn.IsDir(dirNameSource)
	if err != nil {
		return err
	}
	if !isDir {
		return errors.New(fmt.Sprintf("it's NOT a directory. dir name: %s", dirNameSource))
	}

	// check if parent path of destination exists
	dirNameDestParent := filepath.Dir(dirNameDest)
	pathExists, err := conn.PathExists(dirNameDestParent)
	if err != nil {
		return nil
	}
	if !pathExists {
		return errors.New(fmt.Sprintf("parent path of destination does NOT exsists. path: %s", dirNameDest))
	}

	pathSourceBase := filepath.Base(dirNameSource)
	pathDestBase := filepath.Base(dirNameDest)

	// get new destination path to act like shell command "scp -r"
	if pathSourceBase != pathDestBase {
		pathExists, err := conn.PathExists(dirNameDest)
		if err != nil {
			return err
		}
		if pathExists {
			isDir, err := conn.IsDir(dirNameDest)
			if err != nil {
				return err
			}
			if !isDir {
				return errors.New(fmt.Sprintf("it's NOT a directory. dir name: %s", dirNameDest))
			}

			dirNameDest = filepath.Join(dirNameDest, pathSourceBase)
		}
	}

	pathDirMap := make(map[string]string)
	// get map of source path and relative path with destination directory
	err = GetPathDirMapLocal(pathDirMap, dirNameSource, dirNameSource)
	if err != nil {
		return err
	}

	for pathName, relDir := range pathDirMap {
		if relDir == EmptyString {
			//
			relDirSource, err := filepath.Rel(dirNameSource, pathName)
			if err != nil {
				return err
			}

			dirDestAbs := filepath.Join(dirNameDest, relDirSource)
			err = conn.MkdirAll(dirDestAbs)
			continue
		}

		relDir = filepath.Dir(relDir)
		DirDestAbs := filepath.Join(dirNameDest, relDir)
		err = conn.MkdirAll(DirDestAbs)
		if err != nil {
			return err
		}

		fileNameDest := GetFileNameDest(pathName, DirDestAbs)
		err = conn.CopySingleFileToRemote(pathName, fileNameDest)
		if err != nil {
			return err
		}
	}

	return nil
}

// CopyFromRemote copies no matter a directory or a file from remote to local
func (conn *MySSHConn) CopyFromRemote(pathSource, pathDest string) (err error) {
	pathSource = strings.TrimSpace(pathSource)
	pathDest = strings.TrimSpace(pathDest)

	// check if source path is a directory
	isDir, err := conn.IsDir(pathSource)
	if err != nil {
		return err
	}
	if isDir {
		return conn.CopyDirFromRemote(pathSource, pathDest)
	}

	return conn.CopySingleFileFromRemote(pathSource, pathDest)
}

// CopyFromRemote copies no matter a directory or a file from local to remote
func (conn *MySSHConn) CopyToRemote(pathSource, pathDest string) (err error) {
	pathSource = strings.TrimSpace(pathSource)
	pathDest = strings.TrimSpace(pathDest)

	isDir, err := IsDir(pathSource)
	if err != nil {
		return err
	}
	if isDir {
		return conn.CopyDirToRemote(pathSource, pathDest)
	}

	return conn.CopySingleFileToRemote(pathSource, pathDest)
}



