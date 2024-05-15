SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `coin`;
CREATE TABLE `coin`
(
    `id`                  int(0) NOT NULL AUTO_INCREMENT,
    `name`                varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '货币',
    `can_auto_withdraw`   int(0) NOT NULL COMMENT '是否能自动提币',
    `can_recharge`        int(0) NOT NULL COMMENT '是否能充币, 0:不可以 1:可以',
    `can_transfer`        int(0) NOT NULL COMMENT '是否能转账',
    `can_withdraw`        int(0) NOT NULL COMMENT '是否能提币, 0:不可以 1:可以',
    `cny_rate`            double                                                        NOT NULL COMMENT '对人民币汇率',
    `enable_rpc`          int(0) NOT NULL COMMENT '是否支持rpc接口',
    `is_platform_coin`    int(0) NOT NULL COMMENT '是否是平台币',
    `max_tx_fee`          double                                                        NOT NULL COMMENT '最大提币手续费',
    `max_withdraw_amount` decimal(18, 8)                                                NOT NULL COMMENT '最大提币数量',
    `min_tx_fee`          double                                                        NOT NULL COMMENT '最小提币手续费',
    `min_withdraw_amount` decimal(18, 8)                                                NOT NULL COMMENT '最小提币数量',
    `name_cn`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '中文名称',
    `sort`                int(0) NOT NULL COMMENT '排序',
    `status`              tinyint(0) NOT NULL COMMENT '状态 0 正常 1非法',
    `unit`                varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '单位',
    `usd_rate`            double                                                        NOT NULL COMMENT '对美元汇率',
    `withdraw_threshold`  decimal(18, 8)                                                NOT NULL COMMENT '提现阈值',
    `has_legal`           tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是合法币种',
    `cold_wallet_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '冷钱包地址',
    `miner_fee`           decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '转账时付给矿工的手续费',
    `withdraw_scale`      int(0) NOT NULL DEFAULT 4 COMMENT '提币精度',
    `account_type`        int(0) NOT NULL DEFAULT 0 COMMENT '币种账户类型0：默认  1：EOS类型',
    `deposit_address`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '充值地址',
    `infolink`            varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '币种资料链接',
    `information`         varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '币种简介',
    `min_recharge_amount` decimal(18, 8)                                                NOT NULL COMMENT '最小充值数量',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (1, 'Bitcoin', 0, 0, 1, 0, 0, 0, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '比特币', 1, 0, 'BTC', 0, 0.10000000,
        0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (2, 'Bitcoincash', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '比特现金', 1, 0, 'BCH', 0,
        0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (3, 'DASH', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '达世币', 1, 0, 'DASH', 0, 0.10000000,
        0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (4, 'Ethereum', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '以太坊', 1, 0, 'ETH', 0,
        0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (5, 'GalaxyChain', 1, 1, 1, 1, 1, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '银河链', 1, 0, 'GCC', 0,
        0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (6, 'Litecoin', 1, 0, 1, 1, 1, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '莱特币', 1, 0, 'LTC', 0,
        0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (7, 'SGD', 1, 1, 1, 1, 0, 1, 0, 0.0002, 500.00000000, 1, 1.00000000, '新币', 4, 0, 'SGD', 0, 0.10000000, 1, '0',
        0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `zerocoin_coin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`,
                                   `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`,
                                   `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`,
                                   `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`,
                                   `withdraw_scale`,
                                   `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`)
VALUES (8, 'USDT', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '泰达币T', 1, 0, 'USDT', 0, 0.10000000,
        0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);



DROP TABLE IF EXISTS `exchange_coin`;
CREATE TABLE `exchange_coin`
(
    `id`                 bigint(0) NOT NULL AUTO_INCREMENT,
    `symbol`             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易币种名称，格式：BTC/USDT',
    `base_coin_scale`    int(0) NOT NULL COMMENT '基币小数精度',
    `base_symbol`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '结算币种符号，如USDT',
    `coin_scale`         int(0) NOT NULL COMMENT '交易币小数精度',
    `coin_symbol`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '交易币种符号',
    `enable`             int(0) NOT NULL COMMENT '状态，1：启用，2：禁止',
    `fee`                decimal(8, 4)                                                 NOT NULL COMMENT '交易手续费',
    `sort`               int(0) NOT NULL COMMENT '排序，从小到大',
    `enable_market_buy`  int(0) NOT NULL DEFAULT 1 COMMENT '是否启用市价买',
    `enable_market_sell` int(0) NOT NULL DEFAULT 1 COMMENT '是否启用市价卖',
    `min_sell_price`     decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '最低挂单卖价',
    `flag`               int(0) NOT NULL DEFAULT 0 COMMENT '标签位，用于推荐，排序等,默认为0，1表示推荐',
    `max_trading_order`  int(0) NOT NULL DEFAULT 0 COMMENT '最大允许同时交易的订单数，0表示不限制',
    `max_trading_time`   int(0) NOT NULL DEFAULT 0 COMMENT '委托超时自动下架时间，单位为秒，0表示不过期',
    `min_turnover`       decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '最小挂单成交额',
    `clear_time`         bigint(0) NOT NULL DEFAULT 0 COMMENT '清盘时间',
    `end_time`           bigint(0) NOT NULL DEFAULT 0 COMMENT '结束时间',
    `exchangeable`       int(0) NOT NULL DEFAULT 1 COMMENT ' 是否可交易，1：可交易',
    `max_buy_price`      decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '最高买单价',
    `max_volume`         decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '最大下单量',
    `min_volume`         decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT '最小下单量',
    `publish_amount`     decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT ' 活动发行数量',
    `publish_price`      decimal(18, 8)                                                NOT NULL DEFAULT 0.00000000 COMMENT ' 分摊发行价格',
    `publish_type`       int(0) NOT NULL DEFAULT 1 COMMENT '发行活动类型 1:无活动,2:抢购发行,3:分摊发行',
    `robot_type`         int(0) NOT NULL DEFAULT 0 COMMENT '机器人类型',
    `start_time`         bigint(0) NOT NULL DEFAULT 0 COMMENT '开始时间',
    `visible`            int(0) NOT NULL DEFAULT 1 COMMENT ' 前台可见状态',
    `zone`               int(0) NOT NULL DEFAULT 0 COMMENT '交易区域',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;


INSERT INTO `zerocoin_coin`.`exchange_coin`(`symbol`, `base_coin_scale`, `base_symbol`, `coin_scale`, `coin_symbol`,
                                            `enable`, `fee`, `sort`, `enable_market_buy`, `enable_market_sell`,
                                            `min_sell_price`, `flag`, `max_trading_order`, `max_trading_time`,
                                            `min_turnover`, `clear_time`, `end_time`, `exchangeable`, `max_buy_price`,
                                            `max_volume`, `min_volume`, `publish_amount`, `publish_price`,
                                            `publish_type`, `robot_type`, `start_time`, `visible`, `zone`)
VALUES ('BTC/USDT', 2, 'USDT', 2, 'BTC', 1, 0.0010, 1, 1, 1, 0.00000000, 1, 0, 0, 0.00000000, 1640998800000,
        1640998800000, 1, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 1, 0, 1640998800000, 1, 0);
INSERT INTO `zerocoin_coin`.`exchange_coin`(`symbol`, `base_coin_scale`, `base_symbol`, `coin_scale`, `coin_symbol`,
                                            `enable`, `fee`, `sort`, `enable_market_buy`, `enable_market_sell`,
                                            `min_sell_price`, `flag`, `max_trading_order`, `max_trading_time`,
                                            `min_turnover`, `clear_time`, `end_time`, `exchangeable`, `max_buy_price`,
                                            `max_volume`, `min_volume`, `publish_amount`, `publish_price`,
                                            `publish_type`, `robot_type`, `start_time`, `visible`, `zone`)
VALUES ('ETH/USDT', 2, 'USDT', 2, 'ETH', 1, 0.0010, 3, 1, 1, 0.00000000, 0, 0, 0, 0.00000000, 1640998800000,
        1640998800000, 1, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 1, 0, 1640998800000, 1, 0);

SET FOREIGN_KEY_CHECKS = 1;