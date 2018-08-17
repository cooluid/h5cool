/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50161
Source Host           : localhost:3306
Source Database       : cool_account

Target Server Type    : MYSQL
Target Server Version : 50161
File Encoding         : 65001

Date: 2018-08-16 11:12:56
*/

SET FOREIGN_KEY_CHECKS=0;
-- ----------------------------
-- Table structure for `account`
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `openid` char(255) NOT NULL,
  `createtime` datetime NOT NULL,
  `gmlevel` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `index_account_openid` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of account
-- ----------------------------
