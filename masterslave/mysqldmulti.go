package masterslave

import (
	"learn/mylinux"
	"fmt"
	"time"
)

const(
	stopCommand = "mysqld_multi stop 3306"
	rmCommand = "sudo rm -rf "
	startCommand = "mysqld_multi start 3306,3307"
	reportCommand = "mysqld_multi report"
)

func MyMulti(dirInit,dir,mysqlFile3306,mysqlFile3307 string) error {
	err := mylinux.Cmd(stopCommand, dirInit)
	if err != nil{
		return err
	}

	copyCmd := fmt.Sprintf(`cp -r "%s" "%s"`, mysqlFile3306, mysqlFile3307)
	err = mylinux.Cmd(copyCmd,dir)
	if err != nil{
		return err
	}

	rmCmd := rmCommand + mysqlFile3307 +"/data/auto.cnf.toml"
	err = mylinux.Cmd(rmCmd,dir)
	if err != nil{
		return err
	}

	err = mylinux.Cmd(startCommand, dirInit)
	if err != nil{
		return err
	}

	time.Sleep(time.Duration(3)*time.Second)

	err = mylinux.Cmd(reportCommand, dirInit)
	if err != nil{
		return err
	}

	return nil
}


