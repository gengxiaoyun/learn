drop table if exists `test`;
create table `test`(
    `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'primary key id',
    `name` varchar(30) NOT NULL COMMENT 'name',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='user table';