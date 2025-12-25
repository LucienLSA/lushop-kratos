/*
 Navicat Premium Data Transfer

 Source Server         : MySQL5.7
 Source Server Type    : MySQL
 Source Server Version : 50743 (5.7.43-log)
 Source Host           : localhost:3306
 Source Schema         : lushop_order_srv

 Target Server Type    : MySQL
 Target Server Version : 50743 (5.7.43-log)
 File Encoding         : 65001

 Date: 28/07/2025 21:51:40
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ordergoods
-- ----------------------------
DROP TABLE IF EXISTS `ordergoods`;
CREATE TABLE `ordergoods`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `order` int(11) NULL DEFAULT NULL COMMENT '订单ID',
  `goods` int(11) NULL DEFAULT NULL COMMENT '商品ID',
  `goods_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '商品名称',
  `goods_image` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '商品图片',
  `goods_price` float NULL DEFAULT NULL COMMENT '交易时的商品价格(不是最新的价格)',
  `nums` int(11) NULL DEFAULT NULL COMMENT '数量',
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_ordergoods_order`(`order`) USING BTREE,
  INDEX `idx_ordergoods_goods`(`goods`) USING BTREE,
  INDEX `idx_ordergoods_goods_name`(`goods_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of ordergoods
-- ----------------------------
INSERT INTO `ordergoods` VALUES (1, '2025-04-06 20:08:55.857', 0, '2025-04-06 20:08:55.857', 1, 421, '烟台红富士苹果12个 净重2.6kg以上 单果190-240g 新生鲜水果', 'https://py-go.oss-cn-beijing.aliyuncs.com/goods_images/df392d01993cdab9de740fe17798bda1', 44.9, 2, NULL);

-- ----------------------------
-- Table structure for orderinfo
-- ----------------------------
DROP TABLE IF EXISTS `orderinfo`;
CREATE TABLE `orderinfo`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL COMMENT '用户ID',
  `order_sn` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '订单号',
  `pay_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'alipay(支付宝), wechat(微信)',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'PAYING(待支付), TRADE_SUCCESS(成功), TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)',
  `trade_no` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '交易号',
  `order_mount` float NULL DEFAULT NULL COMMENT '总金额',
  `pay_time` datetime NULL DEFAULT NULL COMMENT '支付时间',
  `address` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '收货地址',
  `signer_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '收货名',
  `singer_mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '收获手机',
  `post` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '留言备注',
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_orderinfo_user`(`user`) USING BTREE,
  INDEX `idx_orderinfo_order_sn`(`order_sn`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of orderinfo
-- ----------------------------
INSERT INTO `orderinfo` VALUES (1, '2025-04-06 20:08:55.853', 0, '2025-04-06 20:08:55.853', 0, '202546208852959900182', '', '', '', 89.8, NULL, '北京市', 'bobby', '18787878787', '请尽快发货', NULL);

-- ----------------------------
-- Table structure for shoppingcart
-- ----------------------------
DROP TABLE IF EXISTS `shoppingcart`;
CREATE TABLE `shoppingcart`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `update_time` datetime(3) NULL DEFAULT NULL,
  `user` int(11) NULL DEFAULT NULL COMMENT '用户ID',
  `goods` int(11) NULL DEFAULT NULL COMMENT '商品ID',
  `nums` int(11) NULL DEFAULT NULL COMMENT '数量',
  `checked` tinyint(1) NULL DEFAULT NULL COMMENT '是否选中',
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_shoppingcart_goods`(`goods`) USING BTREE,
  INDEX `idx_shoppingcart_user`(`user`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of shoppingcart
-- ----------------------------
INSERT INTO `shoppingcart` VALUES (2, '2025-04-06 18:30:37.739', 0, '2025-04-06 18:30:37.739', 1, 422, 1, 0, NULL);
INSERT INTO `shoppingcart` VALUES (3, '2025-07-28 20:42:15.743', 0, '2025-07-28 20:42:15.743', 2, 422, 5, 0, NULL);
INSERT INTO `shoppingcart` VALUES (4, '2025-07-28 20:50:29.747', 0, '2025-07-28 20:50:29.747', 2, 423, 5, 0, NULL);

SET FOREIGN_KEY_CHECKS = 1;
