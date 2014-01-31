# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.15)
# Database: trackerdb
# Generation Time: 2014-01-31 02:23:03 +0000
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
  `id` varchar(50) NOT NULL DEFAULT '',
  `username` varchar(30) DEFAULT NULL,
  `password` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `account` WRITE;
/*!40000 ALTER TABLE `account` DISABLE KEYS */;

INSERT INTO `account` (`id`, `username`, `password`)
VALUES
	('111111111111','catlvr666','cats'),
	('222222222222','pro_hacker','7331'),
	('333333333333','swagking420','weed');

/*!40000 ALTER TABLE `account` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table customer
# ------------------------------------------------------------

DROP TABLE IF EXISTS `customer`;

CREATE TABLE `customer` (
  `id` varchar(50) NOT NULL DEFAULT '',
  `phone_number` varchar(20) DEFAULT NULL,
  `address` varchar(30) DEFAULT NULL,
  `email` varchar(30) DEFAULT NULL,
  `first_name` varchar(11) DEFAULT NULL,
  `last_name` varchar(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `customer` WRITE;
/*!40000 ALTER TABLE `customer` DISABLE KEYS */;

INSERT INTO `customer` (`id`, `phone_number`, `address`, `email`, `first_name`, `last_name`)
VALUES
	('111111111111','3346661337','123 Fake St','steven@mensa.org','steven','whaley'),
	('222222222222',NULL,NULL,'jo@google.com',NULL,NULL),
	('333333333333',NULL,NULL,'al@microsoft.com',NULL,NULL);

/*!40000 ALTER TABLE `customer` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table gps_device
# ------------------------------------------------------------

DROP TABLE IF EXISTS `gps_device`;

CREATE TABLE `gps_device` (
  `id` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(50) DEFAULT NULL,
  `customer_id` varchar(50) DEFAULT NULL,
  `latitude` varchar(50) DEFAULT NULL,
  `longitude` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `gps_device` WRITE;
/*!40000 ALTER TABLE `gps_device` DISABLE KEYS */;

INSERT INTO `gps_device` (`id`, `name`, `customer_id`, `latitude`, `longitude`)
VALUES
	('111111111111','device1','12','33.5522','14.2233'),
	('222222222222','device2','13','66.1111','18.1111'),
	('333333333333','device3','14','88.1111','19.3333');

/*!40000 ALTER TABLE `gps_device` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table laptop_device
# ------------------------------------------------------------

DROP TABLE IF EXISTS `laptop_device`;

CREATE TABLE `laptop_device` (
  `id` varchar(100) DEFAULT '',
  `name` varchar(30) DEFAULT NULL,
  `customer_id` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `laptop_device` WRITE;
/*!40000 ALTER TABLE `laptop_device` DISABLE KEYS */;

INSERT INTO `laptop_device` (`id`, `name`, `customer_id`)
VALUES
	('111111111111','laptop','steven'),
	('222222222222','phone','leo'),
	('333333333333','tablet','nathan');

/*!40000 ALTER TABLE `laptop_device` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
