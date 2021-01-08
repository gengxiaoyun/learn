package main

import (
	"fmt"
	"log"
	"flag"
	"runtime"
	"os"
	//"github.com/gengxiaoyun/learn/prepare"
	"learn/prepare"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		err error
		logFileName = flag.String("./log", "InstallMysql.log", "Log file name")
	)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "InstallMysql start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//write log
	log.Printf("Log start! File:%v \n", "InstallMysql.log")

	r := gin.Default()
	r.GET("/",func(c *gin.Context){
		c.JSON(200,c.QueryArray("address"))
		fmt.Println(c.QueryArray("address"))
		fmt.Println(c.Query("user"))
		fmt.Println(c.Query("password"))
		err = prepare.StartMysql(c.QueryArray("address"),c.Query("user"),c.Query("password"))

	})
	r.Run(":8080")

	if err != nil{
		log.Println(err.Error())
	}
	log.Println("AutoInStallMysql succeeded!")
}
