# learn

This project shows how to install mysql using go
======

main.go
---
###### first, decompress the compressed package through the following function.

Untargz()  

###### second, add group and user and assign permissions to folders by calling cmd commands.

Cmd_root(); Adduser(); Cmd()

###### third, you must get the mysql installation dependency package libaio by calling cmd commands.

###### fourth, modify the configuration file. 
eg:  /etc/my.cnf and /etc/init.d/mysql

###### finally, initialize. After changing the temporary password, you can connect to the database to create a table and insert data. what's more, build a master-slave environment.

Dbconnect()

main_test.go
---
This go file is mainly for unit testing the functions in the main file.

command.go
-----
cmd command functions.

command_test.go
----
This go file is mainly for unit testing the functions in the command file.

dbsql.go
----
connect to the database and create a table.

testfile
----
Some test files.

### Now, let's start the installation.

prepare
------
###### you must download a mysql installation package from the official website. this is my folder.
srcfile = "/home/gengxy/mysql/mysql-5.7.31-linux-glibc2.12-x86_64.tar.gz"

###### Configure environment variables. Add the following two lines to the profile file.

_$ vi /etc/profile_

export MYSQL_HOME=/usr/local/mysql
<br>export PATH=$MYSQL_HOME/bin:$PATH

_$ source /etc/profile_

###### Enter the project working directory under GOPATH.

_$ cd $GOPATH/src/learn_
_<br>$ go run main.go_

###### During the running of the program, you may be required to enter the root password.

###### unit tests

_$ go test_