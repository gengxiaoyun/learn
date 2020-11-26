package main

import(
	"github.com/shirou/gopsutil/mem"
	"fmt"
	"strings"
	"os"
	"bufio"
	"io"
	"math"
	"strconv"
)


func GetMemPercent() (string,string,string) {
	memInfo,err:=mem.VirtualMemory()
	if err!=nil{
		fmt.Println(err)
	}

	var a,b int
	a = int(math.Floor(float64(memInfo.Total)/float64(1024*1024*1024)+0.5))
	fmt.Println("a: ",a)
	switch a{
	case 1:
		b=8
	case 2:
		b=16
	case 3:
		b=32
	default:
		b=64
	}

	fmt.Println("b: ",b)

	c:=strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Total) * 0.75 /float64(1024*1024)))))))

	d:=strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Available) * 0.3 / float64(1024*1024)))))))

	return strconv.Itoa(b),c,d

}

func Readline(filename,b,c,d string) {
	f,err:=os.Open(filename)
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer f.Close()

	out,err:=os.OpenFile(filename+"001", os.O_RDWR, 0777)
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer out.Close()
	buf:=bufio.NewReader(f)
	newline:=""
	for {
		line,_,err:=buf.ReadLine()
		if err==io.EOF{
			break
		}
		if err!=nil{
			fmt.Println(err.Error())
		}
		newline = string(line)
		if newline=="thread_cache_size=8"{
			newline = strings.Replace(newline,"8",b,1)
		}
		if newline=="innodb_buffer_pool_size=512M"{
			newline = strings.Replace(newline,"512",c,1)
		}
		if newline=="key_buffer_size=8M" {
			newline = strings.Replace(newline, "8", d, 1)

		}
		_,err1:=out.WriteString(newline+"\n")
		if err1!=nil{
			fmt.Println(err1.Error())
		}
	}

	err1:=os.Remove(filename)
	if err1!=nil{
		fmt.Println(err1.Error())
	}
	err2:=os.Rename(filename+"001",filename)
	if err2!=nil{
		fmt.Println(err2.Error())
	}
}

func main(){

	b,c,d:=GetMemPercent()
	fmt.Println("thread_cache_size: ",b)
	fmt.Println("innodb_buffer_pool_size: ",c)
	fmt.Println("key_buffer_size: ",d)
	Readline("my.cnf",b,c,d)

}