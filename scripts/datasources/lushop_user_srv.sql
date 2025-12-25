/*
 Navicat Premium Data Transfer

 Source Server         : MySQL5.7
 Source Server Type    : MySQL
 Source Server Version : 50743 (5.7.43-log)
 Source Host           : localhost:3306
 Source Schema         : lushop_user_srv

 Target Server Type    : MySQL
 Target Server Version : 50743 (5.7.43-log)
 File Encoding         : 65001

 Date: 28/07/2025 21:51:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `add_time` datetime NULL DEFAULT NULL,
  `update_time` datetime NULL DEFAULT NULL,
  `deleted_at` datetime NULL DEFAULT NULL,
  `is_deleted` tinyint(1) NULL DEFAULT NULL,
  `mobile` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `nick_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `birthday` datetime NULL DEFAULT NULL,
  `gender` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'male' COMMENT '\'female表示女,male表示男\'',
  `role` bigint(20) NULL DEFAULT 1 COMMENT '\'1表示普通用户,2表示管理员\'',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uni_user_mobile`(`mobile`) USING BTREE,
  INDEX `idx_mobile`(`mobile`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES (2, '2025-07-21 15:58:26', '2025-07-21 21:48:30', NULL, 0, '15383660176', '$2a$12$.IMEmklx/M/Un9I4a4DMV.lOBYYhkDSweS3BREKY0Oz5j0lxY0aEm', 'lucien', '2001-04-12 00:00:00', 'male', 2);
INSERT INTO `user` VALUES (3, '2025-07-24 21:45:45', '2025-07-24 21:45:45', NULL, 0, '19821216806', '$2a$12$E0HdxIjgeTT/9ZLSwCXSs.hdXYPnJmJFmk7Ky4WgD113QkCzpd68q', '19821216806', NULL, 'male', 1);
INSERT INTO `user` VALUES (4, '2025-07-24 22:06:29', '2025-07-24 22:06:29', NULL, 0, '19221216806', '$2a$12$rpvfZQhi48h4qLH1J/8gKud/NxyyStGnt44u1N5jm27vMEhLXOMEa', '19221216806', NULL, 'male', 1);
INSERT INTO `user` VALUES (5, '2025-07-24 22:15:27', '2025-07-24 22:15:27', NULL, 0, '19121216806', '$2a$12$eJIgRPmo6aKo/DrvMNkSAu.Y3S/J2Yedl5Xk1nOHjGXlTr0z49XBq', '19121216806', NULL, 'male', 1);

SET FOREIGN_KEY_CHECKS = 1;
