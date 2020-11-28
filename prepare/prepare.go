package prepare

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
	c:=strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Total) * 0.75 /float64(1024*1024)))))))
	d:=strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Available) * 0.3 / float64(1024*1024)))))))
	return strconv.Itoa(b),c,d
}

func Changefile(filename,pathtmp string) error {
	b,c,d:=GetMemPercent()
	f,err:=os.Open(filename)
	if err!=nil{
		fmt.Println(err.Error())
		return err
	}
	defer f.Close()
	out,err:=os.OpenFile(pathtmp, os.O_RDWR|os.O_CREATE, 0777)
	if err!=nil{
		fmt.Println(err.Error())
		return err
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
		if newline=="thread_cache_size=16"{
			newline = strings.Replace(newline,"16",b,1)
		}
		if newline=="innodb_buffer_pool_size=512M"{
			newline = strings.Replace(newline,"512",c,1)
		}
		if newline=="key_buffer_size=16M" {
			newline = strings.Replace(newline, "16", d, 1)

		}
		_,err1:=out.WriteString(newline+"\n")
		if err1!=nil{
			fmt.Println(err1.Error())
			return err
		}
	}
	err1:=os.Remove(filename)
	if err1!=nil{
		fmt.Println(err1.Error())
	}
	err2:=os.Rename(pathtmp,filename)
	if err2!=nil{
		fmt.Println(err2.Error())
	}
	return nil
}
