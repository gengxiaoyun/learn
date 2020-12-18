GRANT ALL PRIVILEGES ON *.* TO `root`@`%` IDENTIFIED BY 'mysql' WITH GRANT OPTION;
GRANT SHUTDOWN ON *.* TO `admin`@`localhost` identified by 'zh3p8ch2we';
grant replication slave, replication client on *.* to `replication`@`%` identified by 'mysql';
FLUSH PRIVILEGES;
create database `dbTest`;
use `dbTest`;
drop table if exists `test`;
create table `test`(
    `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'primary key id',
    `name` varchar(30) NOT NULL COMMENT 'name',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='user table';
insert into `test`(`id`,`name`) values(1,'lily'),(2,'cindy'),(3,'anna');