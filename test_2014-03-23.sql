# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.15)
# Database: trackerDb
# Generation Time: 2014-03-23 17:49:30 +0000
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

LOCK TABLES `account` WRITE;
/*!40000 ALTER TABLE `account` DISABLE KEYS */;

INSERT INTO `account` (`customerId`, `id`, `userName`, `password`)
VALUES
	(1,1,'test@email.com','ee946cfd0649268eae325634c974646f9547ee86');

/*!40000 ALTER TABLE `account` ENABLE KEYS */;
UNLOCK TABLES;


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

LOCK TABLES `customer` WRITE;
/*!40000 ALTER TABLE `customer` DISABLE KEYS */;

INSERT INTO `customer` (`id`, `phoneNumber`, `address`, `email`, `firstName`, `lastName`)
VALUES
	(1,'6661231234',NULL,'test@email.com','steven','whaley');

/*!40000 ALTER TABLE `customer` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table gpsDevice
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gpsDevice`;

CREATE TABLE `gpsDevice` (
  `id` varchar(50) NOT NULL DEFAULT '',
  `deviceName` varchar(50) DEFAULT NULL,
  `customerId` varchar(50) DEFAULT NULL,
  `latitude` varchar(50) DEFAULT NULL,
  `longitude` varchar(50) DEFAULT NULL,
  `isStolen` tinyint(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `gpsDevice` WRITE;
/*!40000 ALTER TABLE `gpsDevice` DISABLE KEYS */;

INSERT INTO `gpsDevice` (`id`, `deviceName`, `customerId`, `latitude`, `longitude`, `isStolen`)
VALUES
	('1','name','1',NULL,NULL,NULL),
	('111111111111','device1','12','33.5522','14.2233',0),
	('222222222222','device2','13','66.1111','18.1111',1),
	('333333333333','device3','14','88.1111','19.3333',0);

/*!40000 ALTER TABLE `gpsDevice` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table ipAddress
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ipAddress`;

CREATE TABLE `ipAddress` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ipAddress` varchar(15) DEFAULT NULL,
  `listId` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `ipAddress` WRITE;
/*!40000 ALTER TABLE `ipAddress` DISABLE KEYS */;

INSERT INTO `ipAddress` (`id`, `ipAddress`, `listId`)
VALUES
	(1,'127.0.01231.1',1),
	(2,'123.1.1.1',1),
	(3,'123.2.23.2',1),
	(4,'123.3.3.3',1),
	(5,'127.0.01231.1',2),
	(6,'123.1.1.1',2),
	(7,'123.2.23.2',2),
	(8,'123.3.3.3',2),
	(9,'127.0.01231.1',3),
	(10,'123.1.1.1',3),
	(11,'123.2.23.2',3),
	(12,'123.3.3.3',3);

/*!40000 ALTER TABLE `ipAddress` ENABLE KEYS */;
UNLOCK TABLES;


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

LOCK TABLES `ipList` WRITE;
/*!40000 ALTER TABLE `ipList` DISABLE KEYS */;

INSERT INTO `ipList` (`id`, `deviceId`, `timestamp`)
VALUES
	(1,1,NULL),
	(2,1,NULL),
	(3,1,NULL);

/*!40000 ALTER TABLE `ipList` ENABLE KEYS */;
UNLOCK TABLES;


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

LOCK TABLES `keyLogs` WRITE;
/*!40000 ALTER TABLE `keyLogs` DISABLE KEYS */;

INSERT INTO `keyLogs` (`deviceId`, `timestamp`, `data`)
VALUES
	(1,NULL,'it worked!keylog test');

/*!40000 ALTER TABLE `keyLogs` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table laptopDevice
# ------------------------------------------------------------

DROP TABLE IF EXISTS `laptopDevice`;

CREATE TABLE `laptopDevice` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `deviceName` varchar(30) DEFAULT NULL,
  `customerId` int(11) DEFAULT NULL,
  `macAddress` varchar(12) DEFAULT NULL,
  `isStolen` tinyint(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `laptopDevice` WRITE;
/*!40000 ALTER TABLE `laptopDevice` DISABLE KEYS */;

INSERT INTO `laptopDevice` (`id`, `deviceName`, `customerId`, `macAddress`, `isStolen`)
VALUES
	(1,NULL,NULL,NULL,1);

/*!40000 ALTER TABLE `laptopDevice` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
