# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.15)
# Database: trackerdb
# Generation Time: 2014-04-01 03:07:38 +0000
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
  `customerId` int(11) DEFAULT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userName` varchar(30) DEFAULT NULL,
  `password` char(40) DEFAULT NULL COMMENT 'TODO: change from char to binary to save space',
  `admin` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `accountToCustomer` (`customerId`),
  CONSTRAINT `accountToCustomer` FOREIGN KEY (`customerId`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table coordinates
# ------------------------------------------------------------

DROP TABLE IF EXISTS `coordinates`;

CREATE TABLE `coordinates` (
  `deviceId` int(11) NOT NULL,
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  KEY `coordsToGpsDevice` (`deviceId`),
  CONSTRAINT `coordsToGpsDevice` FOREIGN KEY (`deviceId`) REFERENCES `gpsDevice` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
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
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `deviceName` varchar(30) DEFAULT NULL,
  `customerId` int(11) DEFAULT NULL,
  `deviceId` varchar(10) DEFAULT NULL COMMENT 'deviceId for gpsDevices are phone numbers',
  `isStolen` int(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `gpsDeviceToCustomer` (`customerId`),
  CONSTRAINT `gpsDeviceToCustomer` FOREIGN KEY (`customerId`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table ipAddress
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ipAddress`;

CREATE TABLE `ipAddress` (
  `listId` int(11) NOT NULL,
  `ipAddress` varchar(20) DEFAULT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  KEY `ipAddressToList` (`listId`),
  CONSTRAINT `ipAddressToList` FOREIGN KEY (`listId`) REFERENCES `ipList` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table ipList
# ------------------------------------------------------------

DROP TABLE IF EXISTS `ipList`;

CREATE TABLE `ipList` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `deviceId` int(11) DEFAULT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `ipListToDevice` (`deviceId`),
  CONSTRAINT `ipListToDevice` FOREIGN KEY (`deviceId`) REFERENCES `laptopDevice` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=latin1;



# Dump of table keyLogs
# ------------------------------------------------------------

DROP TABLE IF EXISTS `keyLogs`;

CREATE TABLE `keyLogs` (
  `deviceId` int(11) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `data` text,
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
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
  `deviceId` varchar(12) DEFAULT NULL COMMENT 'deviceId for laptopDevices are macAddresses',
  `isStolen` int(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `laptopToCustomer` (`customerId`),
  CONSTRAINT `laptopToCustomer` FOREIGN KEY (`customerId`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
