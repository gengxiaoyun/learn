package main

import (
	"archive/tar"
	"os"
	"fmt"
	"compress/gzip"
	"io"
	"strings"
)

func createfile(name string) (*os.File,error) {
	err:=os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name,"/")]),0755)
	if err!=nil{
		return nil,err
	}
	return os.Create(name)
}

func untargz(srcfile string, destfile string) error {
	fr, err := os.Open(srcfile)
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
	// 读取文件
	for {
		h, err := tr.Next()
		if err!=nil{
			if err == io.EOF {
				break
			}else {
				return err
			}
		}
		filename := destfile + h.Name
		if h.Typeflag == tar.TypeDir {
			if err:=os.MkdirAll(filename,os.FileMode(h.Mode));err!=nil{
				return err
			}
		}else{
			file,err := createfile(filename)
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



func main() {
	// file read
	destfile:="/home/gengxy/mysql01/"
	srcfile:="/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"
	untargz(srcfile,destfile)

	//fr, err := os.Open("/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//defer fr.Close()
	//// gzip read
	//gr, err := gzip.NewReader(fr)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//defer gr.Close()
	//// tar read
	//tr := tar.NewReader(gr)
	//// 读取文件
	//for {
	//	h, err := tr.Next()
	//	if err!=nil{
	//		if err == io.EOF {
	//			break
	//		}else {
	//			fmt.Println(err.Error())
	//		}
	//	}
	//	filename := destfile + h.Name
	//	if h.Typeflag == tar.TypeDir {
	//		if err:=os.MkdirAll(filename,os.FileMode(h.Mode));err!=nil{
	//			fmt.Println(err.Error())
	//		}
	//	}else{
	//		file,err := createfile(filename)
	//		if err != nil{
	//			fmt.Println(err.Error())
	//		}
	//		_, err = io.Copy(file, tr)
	//		if err != nil {
	//			fmt.Println(err.Error())
	//		}
	//	}
	//}
	fmt.Println("un tar.gz ok")
}


//func main() {
//	zipFile:="/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"
//	dest:="/home/gengxy/mysql-5.7.31"
//	tarGzUnzip(zipFile,dest)
//}

