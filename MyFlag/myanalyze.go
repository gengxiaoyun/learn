package MyFlag

import (
	"flag"
	"fmt"
	"strings"
)

var (
	address string
	user string
	pass string
)

func init(){
	flag.StringVar(&address,"address","192.168.186.132:3306","set ip and port")
	flag.StringVar(&user,"user","root","set username")
	flag.StringVar(&pass,"pass","root","set password")
}

func FlagCommand() [][]string{

	flag.Parse()
	str := strings.Split(address,",")
	a := len(str)
	arr := make([][]string,a)
	for i:=0;i<a;i++{
		arr[i] = make([]string,2)
	}
	for i:=0;i<a;i++ {
		fmt.Println(str[i])
		newStr := strings.Split(str[i], ":")
		ip := newStr[0]
		port := newStr[1]
		arr[i][0] = ip
		arr[i][1] = port
	}
	fmt.Println(arr,len(arr))
	return arr
}

