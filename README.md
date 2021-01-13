# learn

This project shows how to install mysql using go 
======

go run main.go
---

###### first, parse the configuration file, decompress the compressed package and establish ssh connection through the following function.
Flex()
UnTar()  
MySshConnect()

###### second, transfer configuration files, create new directory, add group and user and assign permissions to folders, get the mysql installation dependency package libaio.
CopyToRemote()
BasicWork()

###### third, initialize and start. After changing the temporary password, you can connect to the database to create a table and insert data. what's more, build a master-slave environment.
InitMysql()
MyMulti()

testfile
----
Some test files.

### Now, let's start the installation.

prepare work
------
###### you must download a mysql installation package from the official website. this is my folder.
srcfile = "/home/gengxy/mysql/mysql.tar.gz"

###### Enter the project working directory under GOPATH.
_$ cd $GOPATH/src/learn_
_<br>$ go run main.go_

###### Enter the URL in the browser.
_http://localhost:8080/?address=ip1:port1&address=ip2:port2&address=ip3:port3&user=xxx&password=xxx_
eg: http://localhost:8080/?address=192.168.186.137:3306&address=192.168.186.137:3307&address=192.168.186.138:3308&user=root&password=root
###### unit tests
_$ go test_