stop slave;
change master to master_host='192.168.186.131', master_port=3306, master_user='replication', master_password='mysql', master_auto_position=1;
start slave;


