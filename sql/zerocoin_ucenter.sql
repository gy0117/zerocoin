SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;


DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`                            bigint(0) NOT NULL AUTO_INCREMENT,
    `ali_no`                        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `qr_code_url`                   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `appeal_success_times`          int(0) NOT NULL,
    `appeal_times`                  int(0) NOT NULL,
    `application_time`              bigint(0) NOT NULL,
    `avatar`                        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `bank`                          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `branch`                        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `card_no`                       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `certified_business_apply_time` bigint(0) NOT NULL,
    `certified_business_check_time` bigint(0) NOT NULL,
    `certified_business_status`     int(0) NOT NULL,
    `channel_id`                    int(0) NOT NULL DEFAULT 0,
    `email`                         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `first_level`                   int(0) NOT NULL,
    `google_date`                   bigint(0) NOT NULL,
    `google_key`                    varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `google_state`                  int(0) NOT NULL DEFAULT 0,
    `id_number`                     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `inviter_id`                    bigint(0) NOT NULL,
    `is_channel`                    int(0) NOT NULL DEFAULT 0,
    `jy_password`                   varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `last_login_time`               bigint(0) NOT NULL,
    `city`                          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `country`                       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `district`                      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `province`                      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `login_count`                   int(0) NOT NULL,
    `login_lock`                    int(0) NOT NULL,
    `margin`                        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `user_level`                    int(0) NOT NULL,
    `mobile_phone`                  varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `password`                      varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `promotion_code`                varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `publish_advertise`             int(0) NOT NULL,
    `real_name`                     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `real_name_status`              int(0) NOT NULL,
    `registration_time`             bigint(0) NOT NULL,
    `salt`                          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `second_level`                  int(0) NOT NULL,
    `sign_in_ability`               tinyint(4) NOT NULL DEFAULT b'1',
    `status`                        int(0) NOT NULL,
    `third_level`                   int(0) NOT NULL,
    `token`                         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `token_expire_time`             bigint(0) NOT NULL,
    `transaction_status`            int(0) NOT NULL COMMENT '0：禁止交易',
    `transaction_time`              bigint(0) NOT NULL,
    `transactions`                  int(0) NOT NULL,
    `username`                      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `qr_we_code_url`                varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `wechat`                        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `local`                         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `integration`                   bigint(0) NOT NULL DEFAULT 0,
    `user_grade_id`                 bigint(0) NOT NULL DEFAULT 1 COMMENT '等级id',
    `kyc_status`                    int(0) NOT NULL DEFAULT 0 COMMENT 'kyc等级',
    `generalize_total`              bigint(0) NOT NULL DEFAULT 0 COMMENT '注册赠送积分',
    `inviter_parent_id`             bigint(0) NOT NULL DEFAULT 0,
    `super_partner`                 varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `kick_fee`                      decimal(19, 2)                                                NOT NULL,
    `power`                         decimal(8, 4)                                                 NOT NULL DEFAULT 0.0000 COMMENT '个人矿机算力(每日维护)',
    `team_level`                    int(0) NOT NULL DEFAULT 0 COMMENT '团队人数(每日维护)',
    `team_power`                    decimal(8, 4)                                                 NOT NULL DEFAULT 0.0000 COMMENT '团队矿机算力(每日维护)',
    `user_level_id`                 bigint(0) NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `UK_gc3jmn7c2abyo3wf6syln5t2i`(`username`) USING BTREE,
    UNIQUE INDEX `UK_10ixebfiyeqolglpuye0qb49u`(`mobile_phone`) USING BTREE,
    INDEX                           `FKbt72vgf5myy3uhygc90xna65j`(`local`) USING BTREE,
    INDEX                           `FK8jlqfg5xqj5epm9fpke6iotfw`(`user_level_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;



DROP TABLE IF EXISTS `user_address`;
CREATE TABLE `user_address`
(
    `id`          bigint(0) NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(0) NOT NULL,
    `coin_id`     bigint(0) NOT NULL,
    `address`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
    `remark`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '备注',
    `status`      tinyint(0) NOT NULL COMMENT '0 正常 1 非法',
    `create_time` bigint(0) NOT NULL,
    `delete_time` bigint(0) NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;



DROP TABLE IF EXISTS `user_transaction`;
CREATE TABLE `user_transaction`
(
    `id`           bigint(0) NOT NULL AUTO_INCREMENT COMMENT '编号',
    `address`      varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '充值或提现地址、或转账地址',
    `amount`       decimal(18, 8)                                                NOT NULL COMMENT '充币金额',
    `create_time`  bigint(0) NOT NULL COMMENT '创建时间',
    `fee`          decimal(19, 8)                                                NOT NULL COMMENT '交易手续费',
    `flag`         int(0) NOT NULL DEFAULT 0 COMMENT '标识位',
    `user_id`      bigint(0) NOT NULL COMMENT '会员ID',
    `symbol`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '币种名称 比如BTC',
    `type`         int(0) NOT NULL COMMENT '交易类型,  0 RECHARGE 充值类型 1 WITHDRAW 提现类型 ',
    `discount_fee` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '折扣手续费',
    `real_fee`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '实收手续费',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


DROP TABLE IF EXISTS `user_wallet`;
CREATE TABLE `user_wallet`
(
    `id`                  bigint(0) NOT NULL AUTO_INCREMENT,
    `address`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '充值地址',
    `balance`             decimal(18, 8)                                                NOT NULL COMMENT '可用余额',
    `frozen_balance`      decimal(18, 8)                                                NOT NULL COMMENT '冻结余额',
    `release_balance`     decimal(18, 8)                                                NOT NULL COMMENT '待释放余额',
    `is_lock`             int(0) NOT NULL DEFAULT 0 COMMENT '钱包不是锁定 0 否 1 是',
    `user_id`             bigint(0) NOT NULL COMMENT '用户id',
    `version`             int(0) NOT NULL COMMENT '版本',
    `coin_id`             bigint(0) NOT NULL COMMENT '货币id',
    `to_released`         decimal(18, 8)                                                NOT NULL COMMENT '待释放总量',
    `coin_name`           varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '货币名称',
    `address_private_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '私钥地址',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `UKm68bscpof0bpnxocxl4qdnvbe`(`user_id`, `coin_id`) USING BTREE,
    INDEX                 `FKf9tgbp9y9py8t9c5xj0lllcib`(`coin_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

SET FOREIGN_KEY_CHECKS = 1;