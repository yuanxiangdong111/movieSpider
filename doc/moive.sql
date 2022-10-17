-- MySQL dump 10.13  Distrib 5.7.30, for el7 (x86_64)
--
-- Host: 127.0.0.1    Database: movie
-- ------------------------------------------------------
-- Server version	5.7.18

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
-- Current Database: `movie`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `movie` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;

USE `movie`;

--
-- Table structure for table `douban_video`
--

DROP TABLE IF EXISTS `douban_video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `douban_video` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `names` varchar(255) NOT NULL COMMENT '片名列表',
  `douban_id` varchar(255) NOT NULL COMMENT '豆瓣ID',
  `imdb_id` varchar(255) NOT NULL COMMENT 'imdbID',
  `row_data` longtext NOT NULL COMMENT '原始数据',
  `timestamp` bigint(11) NOT NULL COMMENT '修改创建时间',
  `type` varchar(255) NOT NULL COMMENT '类型',
  `playable` varchar(255) NOT NULL COMMENT '是否可以播放',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`names`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `feed_video`
--

DROP TABLE IF EXISTS `feed_video`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `feed_video` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '片名',
  `torrent_name` varchar(255) NOT NULL COMMENT '种子名',
  `torrent_url` varchar(255) NOT NULL COMMENT '种子引用地址',
  `magnet` longtext NOT NULL COMMENT '磁力链接',
  `year` varchar(255) NOT NULL COMMENT '年份',
  `type` varchar(255) NOT NULL COMMENT 'tv或movie',
  `row_data` longtext COMMENT '原始数据',
  `web` varchar(255) NOT NULL COMMENT '站点',
  `download` int(11) NOT NULL COMMENT '1:已经下载',
  `timestamp` bigint(11) NOT NULL COMMENT '修改创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`,`torrent_name`)
) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-10-16  1:05:13
