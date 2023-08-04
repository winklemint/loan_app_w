-- MySQL dump 10.13  Distrib 8.0.33, for Win64 (x86_64)
--
-- Host: localhost    Database: loan_application
-- ------------------------------------------------------
-- Server version	8.0.33

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `loan_details_table`
--

DROP TABLE IF EXISTS `loan_details_table`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `loan_details_table` (
  `id` int NOT NULL AUTO_INCREMENT,
  `loan_type` varchar(40) DEFAULT NULL,
  `loan_amount` float DEFAULT NULL,
  `pincode` int DEFAULT NULL,
  `tenure` int DEFAULT NULL,
  `employment_type` varchar(90) DEFAULT NULL,
  `gross_monthly_income` float DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `last_modified` datetime DEFAULT NULL,
  `status` varchar(45) DEFAULT 'Pending',
  `is_delete` int DEFAULT '0',
  `user_id` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=172 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `loan_details_table`
--

LOCK TABLES `loan_details_table` WRITE;
/*!40000 ALTER TABLE `loan_details_table` DISABLE KEYS */;
INSERT INTO `loan_details_table` VALUES (164,'homeimprovement',676332,345234,0,'BussinessOwner',676332,'2023-07-10 18:28:41','2023-07-10 18:28:41','Pending',1,42),(165,'loanagainst',778134,123456,0,'Independent',778134,'2023-07-10 18:39:43','2023-07-10 18:39:43','Pending',0,41),(166,'loanagainst',413794,123456,0,'Independent',413794,'2023-07-10 19:29:37','2023-07-10 19:29:37','Pending',0,41),(167,'BuyHome',613272,123456,0,'SelfEmployed',613272,'2023-07-19 12:32:20','2023-07-19 12:32:20','Pending',1,41),(168,'loanagainst',223326,123456,0,'Independent',223326,'2023-07-19 13:58:03','2023-07-19 13:58:03','Pending',0,42),(169,'loanagainst',223326,123456,0,'Independent',223326,'2023-07-19 13:58:06','2023-07-19 13:58:06','Pending',0,42),(170,'buyplot',339151,123456,0,'SelfEmployed',339151,'2023-07-19 15:52:00','2023-07-19 15:52:00','Pending',0,41),(171,'loanagainst',447005,123456,0,'BussinessOwner',447005,'2023-07-21 15:00:00','2023-07-21 15:00:00','Pending',0,46);
/*!40000 ALTER TABLE `loan_details_table` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-07-28 11:16:45
