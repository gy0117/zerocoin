SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;


DROP TABLE IF EXISTS `exchange_order`;
CREATE TABLE `exchange_order`
(
    `id`             bigint(0) NOT NULL AUTO_INCREMENT,
    `order_id`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '订单id',
    `amount`         decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '买入或者卖出量',
    `base_symbol`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '结算单位',
    `canceled_time`  bigint(0) NOT NULL COMMENT '取消时间',
    `coin_symbol`    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '币单位',
    `completed_time` bigint(0) NOT NULL COMMENT '完成时间',
    `direction`      int(0) NOT NULL COMMENT '订单方向 0 买 1 卖',
    `user_id`        bigint(0) NOT NULL,
    `price`          decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '挂单价格',
    `status`         int(0) NOT NULL COMMENT '订单状态 0 交易中 1 完成 2 取消 3 超时',
    `symbol`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易对',
    `time`           bigint(0) NOT NULL COMMENT '挂单时间',
    `traded_amount`  decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '成交量',
    `turnover`       decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '成交额 ',
    `type`           int(0) NOT NULL COMMENT '挂单类型 0 市场价 1 最低价',
    `use_discount`   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '是否使用折扣 0 不使用 1使用',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `order_id`(`order_id`) USING BTREE,
    INDEX            `index_user_id_time`(`user_id`, `time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 84 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


SET FOREIGN_KEY_CHECKS = 1;