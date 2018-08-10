/*
配置数据库初始化脚本
*/

/*商品表*/
INSERT `t_common_config`(`key`, `subkey`, `value`) values(
    'charge', 'item_list', 
    '{
    "android": {
        "default": [
            {
                "item_id": 1,
                "name": "金豆 100",
                "tag": "热卖",
                "price": 600,
                "coin": 100,
                "present_coin": 0
            },
            {
                "item_id": 2,
                "name": "金豆 1000",
                "tag": "特惠",
                "price": 800,
                "coin": 1000,
                "present_coin": 0
            }
        ],
        "city400200": [
            {
                "item_id": 1,
                "name": "金豆 100",
                "tag": "热卖",
                "price": 600,
                "coin": 100,
                "present_coin": 0
            },
            {
                "item_id": 2,
                "name": "金豆 1000",
                "tag": "特惠",
                "price": 800,
                "coin": 1000,
                "present_coin": 0
            }
        ]
    },
    "iphone": {
        "default": [
            {
                "item_id": 1,
                "name": "金豆 100",
                "tag": "热卖",
                "price": 600,
                "coin": 100,
                "present_coin": 0
            },
            {
                "item_id": 2,
                "name": "金豆 1000",
                "tag": "特惠",
                "price": 800,
                "coin": 1000,
                "present_coin": 0
            }
        ],
        "city400200": [
            {
                "item_id": 1,
                "name": "金豆 100",
                "tag": "热卖",
                "price": 600,
                "coin": 100,
                "present_coin": 0
            },
            {
                "item_id": 2,
                "name": "金豆 1000",
                "tag": "特惠",
                "price": 800,
                "coin": 1000,
                "present_coin": 0
            }
        ]
    }
}'
);

/*每日最大充值数*/
INSERT `t_common_config`(`key`, `subkey`, `value`) values (
    'charge',
    'day_max',
    '{
    "max_charge": 200000
    }'
);

/*游戏配置*/
INSERT INTO `t_game_config` VALUES (1, 1, '血流麻将', 1, 4, 4, NULL, NULL, NULL, NULL, NULL, '2018-08-07 19:01:33', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (2, 2, '血战麻将', 1, 4, 4, NULL, NULL, NULL, NULL, NULL, '2018-08-07 19:03:29', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (3, 3, '斗地主', 2, 3, 3, NULL, NULL, NULL, NULL, NULL, '2018-08-07 20:36:58', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (4, 4, '二人麻将', 1, 2, 2, NULL, NULL, NULL, NULL, NULL, '2018-08-07 20:37:11', NULL, NULL, NULL);

/*游戏场次配置*/
INSERT INTO `t_game_level_config` VALUES (1, 1, 1, '新手场', 1, 1, 0, 1000000, 1, 1, 1, NULL, 1, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (2, 2, 1, '新手场', 1, 1, 0, 1000000, 1, 1, 1, NULL, 1, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (3, 3, 1, '新手场', 1, 1, 0, 1000000, 1, 1, 1, NULL, 1, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (4, 4, 1, '新手场', 1, 2, 0, 1000000, 1, 1, 1, NULL, 1, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (5, 1, 2, '中级场', 1, 5, 200, 1000000, 1, 100, 1, NULL, 1, NULL, '2018-08-10 10:40:50', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (6, 2, 2, '中级场', 1, 5, 200, 1000000, 1, 100, 1, NULL, 1, NULL, '2018-08-10 10:40:52', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (7, 3, 2, '中级场', 1, 5, 200, 1000000, 1, 100, 1, NULL, 1, NULL, '2018-08-10 10:41:35', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (8, 4, 2, '中级场', 1, 5, 200, 1000000, 1, 100, 1, NULL, 1, NULL, '2018-08-10 10:43:08', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (9, 1, 3, '大师场', 1, 10, 800, 1000000, 1, 200, 1, NULL, 1, NULL, '2018-08-10 10:43:11', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (10, 2, 3, '大师场', 1, 10, 800, 1000000, 1, 200, 1, NULL, 1, NULL, '2018-08-10 10:43:42', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (11, 3, 3, '大师场', 1, 10, 800, 1000000, 1, 200, 1, NULL, 1, NULL, '2018-08-10 10:45:00', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (12, 4, 3, '大师场', 1, 10, 800, 1000000, 1, 200, 1, NULL, 1, NULL, '2018-08-10 10:45:02', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (13, 1, 4, '土豪场', 1, 100, 100000, 10000000, 1, 500, 1, NULL, 1, NULL, '2018-08-10 10:47:47', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (14, 2, 4, '土豪场', 1, 100, 100000, 10000000, 1, 500, 1, NULL, 1, NULL, '2018-08-10 10:48:20', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (15, 3, 4, '土豪场', 1, 100, 100000, 10000000, 1, 500, 1, NULL, 1, NULL, '2018-08-10 10:48:50', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (16, 4, 4, '土豪场', 1, 100, 100000, 10000000, 1, 500, 1, NULL, 1, NULL, '2018-08-10 10:49:31', NULL, NULL, NULL);