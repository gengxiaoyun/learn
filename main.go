package main

import (
	"flag"
	"fmt"
	"github.com/gengxiaoyun/learn/prepare"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"runtime"
)

func main() {
	var (
		err         error
		logFileName = flag.String("./log", "InstallMysql.log", "Log file name")
	)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	//set logfile Stdout
	_, err = os.Stat(*logFileName)
	if err == nil {
		err = os.Remove(*logFileName)
		if err != nil {
			log.Println(err.Error())
		}
	}
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
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, c.QueryArray("address"))
		err = prepare.StartMysql(c.QueryArray("address"), c.Query("user"), c.Query("password"))
		if err != nil {
			log.Println(err.Error())
		}
		log.Println("AutoInStallMysql succeeded!")
	})
	r.Run(":8080")

}
