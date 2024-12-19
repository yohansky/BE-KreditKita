-- MySQL dump 10.13  Distrib 5.7.22, for Linux (x86_64)
--
-- Host: localhost    Database: kredit
-- ------------------------------------------------------
-- Server version	5.7.22

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `consumers`
--

DROP TABLE IF EXISTS `consumers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `consumers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `nik` bigint(20) unsigned DEFAULT NULL,
  `full_name` longtext,
  `legal_name` longtext,
  `place_of_birth` longtext,
  `date_of_birth` longtext,
  `salary` double DEFAULT NULL,
  `photo_ktp` longtext,
  `photo_selfie` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_consumers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `consumers`
--

LOCK TABLES `consumers` WRITE;
/*!40000 ALTER TABLE `consumers` DISABLE KEYS */;
INSERT INTO `consumers` VALUES (6,1234567,'Yohanes Hubert','Joey','Jakarta','27-04-2000',6000000,'https://res.cloudinary.com/dcwt4iksg/image/upload/v1734513165/kredit/foto_ktp/1734513164_miebabat.jpg','https://res.cloudinary.com/dcwt4iksg/image/upload/v1734513166/kredit/foto_selfie/1734513165_miebabat.jpg','2024-12-18 09:12:46.133','2024-12-18 09:12:46.133',NULL),(7,12,'Joey Aprilio','April','Jakarta','27-04-2000',7000000,'https://res.cloudinary.com/dcwt4iksg/image/upload/v1734525038/kredit/foto_ktp/1734525032_miebabat.jpg','https://res.cloudinary.com/dcwt4iksg/image/upload/v1734525044/kredit/foto_selfie/1734525038_miebabat.jpg','2024-12-18 12:30:45.626','2024-12-18 12:30:45.626',NULL);
/*!40000 ALTER TABLE `consumers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `limits`
--

DROP TABLE IF EXISTS `limits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `limits` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `tenor1` double DEFAULT NULL,
  `tenor2` double DEFAULT NULL,
  `tenor3` double DEFAULT NULL,
  `tenor4` double DEFAULT NULL,
  `remaining_tenor1` double DEFAULT NULL,
  `remaining_tenor2` double DEFAULT NULL,
  `remaining_tenor3` double DEFAULT NULL,
  `remaining_tenor4` double DEFAULT NULL,
  `consumer_id` bigint(20) unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_limits_consumer` (`consumer_id`),
  KEY `idx_limits_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_limits_consumer` FOREIGN KEY (`consumer_id`) REFERENCES `consumers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `limits`
--

LOCK TABLES `limits` WRITE;
/*!40000 ALTER TABLE `limits` DISABLE KEYS */;
INSERT INTO `limits` VALUES (6,618000,412000,206000,1113000,618000,412000,206000,1113000,6,'2024-12-18 09:12:46.146','2024-12-18 09:12:46.146',NULL),(7,2100000,2310000,2541000,2795100,2100000,2310000,141000,2899100,7,'2024-12-18 12:30:45.640','2024-12-19 01:28:37.310',NULL);
/*!40000 ALTER TABLE `limits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transactions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `nomor_kontrak` longtext,
  `otr` double DEFAULT NULL,
  `admin_fee` double DEFAULT NULL,
  `jumlah_cicilan` double DEFAULT NULL,
  `jumlah_bunga` double DEFAULT NULL,
  `nama_asset` longtext,
  `consumer_id` bigint(20) unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_transactions_consumer` (`consumer_id`),
  KEY `idx_transactions_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_transactions_consumer` FOREIGN KEY (`consumer_id`) REFERENCES `consumers` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
INSERT INTO `transactions` VALUES (3,'7-3',2400000,50000,816000,48000,'Motor Honda',7,'2024-12-18 12:38:13.392','2024-12-18 12:38:13.392',NULL),(4,'7-6',2600000,50000,442000,52000,'Mobil Kijang',7,'2024-12-18 12:41:39.130','2024-12-18 12:41:39.130',NULL);
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-12-19  4:31:41
