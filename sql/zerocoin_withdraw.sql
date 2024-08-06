SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `withdraw_record`;
CREATE TABLE `withdraw_record`
(
    `id`                 bigint(0) NOT NULL AUTO_INCREMENT,
    `user_id`            bigint(0) NOT NULL,
    `coin_id`            bigint(0) NOT NULL,
    `total_amount`       decimal(18, 8)                                                NOT NULL COMMENT '申请总数量',
    `fee`                decimal(18, 8)                                                NOT NULL COMMENT '手续费',
    `arrived_amount`     decimal(18, 8)                                                NOT NULL COMMENT '预计到账数量',
    `address`            varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '提现地址',
    `remark`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '备注',
    `transaction_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易编号',
    `can_auto_withdraw`  tinyint(0) NOT NULL COMMENT '是否能自动提币 0 false 1 true',
    `isAuto`             tinyint(0) NOT NULL COMMENT '是否是自动提现的 0 false 1 true',
    `status`             tinyint(0) NOT NULL COMMENT '状态 0 审核中 1 等待放币 2 失败 3 成功',
    `create_time`        bigint(0) NOT NULL COMMENT '创建时间',
    `deal_time`          bigint(0) NOT NULL COMMENT '处理时间',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;