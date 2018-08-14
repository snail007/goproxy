create database IF NOT EXISTS taskdb;
use taskdb;

DROP TABLE IF EXISTS `pattern_table`;
CREATE TABLE `pattern_table` (
  `pattern` varchar(1024) NOT NULL,
  `proxyName` varchar(1024),
  `priority` int NOT NULL,
  PRIMARY KEY (`pattern`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `proxy`;
CREATE TABLE `proxy` (
  `source` varchar(1024) NOT NULL,
  `endpoint` varchar(1024),
  `port` int,
  `proxyType` varchar(1024) NOT NULL,
  `user` varchar(1024),
  `password` varchar(1024),
  `apiEndpoint` varchar(1024)
  PRIMARY KEY (`source`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
