# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.15)
# Database: trackerdb
# Generation Time: 2014-02-01 23:29:25 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table account
# ------------------------------------------------------------

DROP TABLE IF EXISTS `account`;

CREATE TABLE `account` (
  `customerId` int(11) NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userName` varchar(30) DEFAULT NULL,
  `password` char(40) DEFAULT NULL COMMENT 'TODO: change from char to binary to save space',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table customer
# ------------------------------------------------------------

DROP TABLE IF EXISTS `customer`;

CREATE TABLE `customer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `phoneNumber` varchar(20) DEFAULT NULL,
  `address` varchar(30) DEFAULT NULL,
  `email` varchar(30) DEFAULT NULL,
  `firstName` varchar(20) DEFAULT NULL,
  `lastName` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table gpsDevice
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gpsDevice`;

CREATE TABLE `gpsDevice` (
  `id` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(50) DEFAULT NULL,
  `customer_id` varchar(50) DEFAULT NULL,
  `latitude` varchar(50) DEFAULT NULL,
  `longitude` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `gpsDevice` WRITE;
/*!40000 ALTER TABLE `gpsDevice` DISABLE KEYS */;

INSERT INTO `gpsDevice` (`id`, `name`, `customer_id`, `latitude`, `longitude`)
VALUES
  ('111111111111','device1','12','33.5522','14.2233'),
  ('222222222222','device2','13','66.1111','18.1111'),
  ('333333333333','device3','14','88.1111','19.3333');

/*!40000 ALTER TABLE `gpsDevice` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table ipAddress
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ipAddress`;

CREATE TABLE `ipAddress` (
  `listId` int(11) NOT NULL,
  `ipAddress` varchar(15) DEFAULT NULL,
  PRIMARY KEY (`listId`),
  CONSTRAINT `ipAddressToList` FOREIGN KEY (`listId`) REFERENCES `ipList` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table ipList
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ipList`;

CREATE TABLE `ipList` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `deviceId` int(11) DEFAULT NULL,
  `timestamp` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `ipListToDevice` (`deviceId`),
  CONSTRAINT `ipListToDevice` FOREIGN KEY (`deviceId`) REFERENCES `laptopDevice` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table keyLogs
# ------------------------------------------------------------

DROP TABLE IF EXISTS `keyLogs`;

CREATE TABLE `keyLogs` (
  `deviceId` int(11) NOT NULL,
  `timestamp` timestamp NULL DEFAULT NULL,
  `data` text,
  KEY `keyLogToDevice` (`deviceId`),
  CONSTRAINT `keyLogToDevice` FOREIGN KEY (`deviceId`) REFERENCES `laptopDevice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table laptopDevice
# ------------------------------------------------------------

DROP TABLE IF EXISTS `laptopDevice`;

CREATE TABLE `laptopDevice` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `deviceName` varchar(30) DEFAULT NULL,
  `customerId` int(11) DEFAULT NULL,
  `macAddress` varchar(12) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
