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
	"learn/mylinux"
)


func GetMemPercent() (string,string,string,error) {
	memInfo,err := mem.VirtualMemory()
	if err != nil{
		return nil,nil,nil,err
	}
	var a,b int
	a = int(math.Floor(float64(memInfo.Total)/float64(1024*1024*1024)+0.5))
	switch a {
	case 1:
		b = 8
	case 2:
		b = 16
	case 3:
		b = 32
	default:
		b = 64
	}
	c := strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Total) * 0.75 /float64(1024*1024)))))))
	d := strconv.Itoa(int(math.Exp2(float64(int(math.Log2(float64(memInfo.Available) * 0.3 / float64(1024*1024)))))))
	return strconv.Itoa(b),c,d,nil
}

func ChangeConfFile(srcCnf,pathTmp string) error {
	b,c,d,err := GetMemPercent()
	if err != nil{
		return err
	}
	f,err := os.Open(srcCnf)
	if err != nil{
		return err
	}
	defer f.Close()
	out,err := os.OpenFile(pathTmp, os.O_RDWR|os.O_CREATE, 0777)
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
		if newline == "thread_cache_size=16"{
			newline = strings.Replace(newline,"16",b,1)
		}
		if newline == "innodb_buffer_pool_size=512M"{
			newline = strings.Replace(newline,"512",c,1)
		}
		if newline == "key_buffer_size=16M" {
			newline = strings.Replace(newline, "16", d, 1)

		}
		_,err = out.WriteString(newline+"\n")
		if err != nil{
			return err
		}
	}
	err = os.Remove(srcCnf)
	if err != nil{
		return err
	}
	err = os.Rename(pathTmp,srcCnf)
	if err != nil{
		return err
	}

	return nil
}

// copy my.cnf
func CopyConfFile(srcCnf,destCnf string) error {
	cpCmd := fmt.Sprintf(`cp "%s" "%s"`, srcCnf, destCnf)
	err := mylinux.CmdRoot(cpCmd)
	if err != nil{
		return err
	}
	return nil
}
