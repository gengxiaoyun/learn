package basic

import (
	"fmt"
	"learn/mylinux"
	"archive/tar"
	"io"
	"os"
	"strings"
	"compress/gzip"
	"os/exec"
	"bytes"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type MyConn struct {
	HostIp   string
	PortNum  int
	UserName string
	UserPass string
}

type MySSHConn struct {
	MyConn
	SSHClient *ssh.Client
	*sftp.Client
}



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

// set password for user
func AddUser(command,username,password string) error {
	var err error
	cmd := exec.Command("/bin/bash","-c",command)
	err = cmd.Start()
	if err != nil{
		return err
	}
	err = cmd.Wait()
	if err != nil{
		return err
	}
	ps := exec.Command("echo",password)
	grep := exec.Command("passwd","-stdin",username)
	r,w := io.Pipe()
	defer r.Close()
	defer w.Close()
	ps.Stdout = w
	grep.Stdin = r
	var buffer bytes.Buffer
	grep.Stdout = &buffer
	err = ps.Start()
	if err != nil{
		return err
	}
	err = grep.Start()
	if err != nil{
		return err
	}
	err = ps.Wait()
	if err != nil{
		return err
	}
	w.Close()
	err = grep.Wait()
	if err != nil{
		return err
	}
	_,err = io.Copy(os.Stdout,&buffer)
	if err != nil{
		return err
	}
	//fmt.Println("Set password for user successfully")
	return nil
}

// Add user and group
func UserAndGroup(groupCommand,userCommand,user,group,
	file,userPassword,filePath,sLib,iLib,dir string) error {
	err := mylinux.CmdRoot(groupCommand)
	if err != nil{
		return err
	}
	err = AddUser(userCommand,user,userPassword)
	if err != nil{
		return err
	}

	mk_cmd := "sudo mkdir -p "+filePath
	err = mylinux.Cmd(mk_cmd,dir)
	if err != nil{
		return err
	}

	err = MyOwn(user,group,file,dir)
	if err != nil{
		return err
	}
	err = MyMod(file,dir)
	if err != nil{
		return err
	}

	err = InstallTool(sLib,iLib,dir)
	if err != nil{
		return err
	}
	return nil
}

func MyOwn(user,group,file,dir string) error {
	ownCmd := fmt.Sprintf(`sudo chown -R "%s"."%s" "%s"`,user,group,file)
	err := mylinux.Cmd(ownCmd,dir)
	if err != nil{
		return err
	}
	return nil
}

func MyMod(file,dir string) error {
	modCmd := fmt.Sprintf(`sudo chmod -R g+rwx "%s"`,file)
	err := mylinux.Cmd(modCmd,dir)
	if err != nil{
		return err
	}
	return nil
}

// install libaio
func InstallTool(sLib,iLib,dir string) error {
	err := mylinux.Cmd(sLib,dir)
	if err != nil{
		return err
	}
	err = mylinux.CmdRoot(iLib)
	if err != nil{
		return err
	}
	return nil
}