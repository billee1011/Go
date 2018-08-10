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

INSERT `t_common_config`(`key`, `subkey`, `value`) values ( 
    'prop', 
    'interactive',
    '[
            {
                "propID": 1,
                "name": "rose",
                "attrType": 1,
                "attrID":1,
                "attrValue":-100,
                "attrLimit":10000
            },
            {
                "propID": 2,
                "name": "beer",
                "attrType": 1,
                "attrID":1,
                "attrValue":-100,
                "attrLimit":10000
            },
            {
                "propID": 3,
                "name": "bomb",
                "attrType": 1,
                "attrID":1,
                "attrValue":-100,
                "attrLimit":10000
            },
            {
                "propID": 4,
                "name": "grabChicken",
                "attrType": 1,
                "attrID":1,
                "attrValue":-100,
                "attrLimit":10000
            },
            {
                "propID": 5,
                "name": "eggGun",
                "attrType": 1,
                "attrID":1,
                "attrValue":-10000,
                "attrLimit":500000
            }
    ]
    '
);

/*游戏配置*/ 
INSERT INTO `t_game_config` VALUES (1, 1, '血流麻将', 1, 4, 4, NULL, NULL, NULL, NULL, NULL, '2018-08-07 19:01:33', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (2, 2, '血战麻将', 1, 4, 4, NULL, NULL, NULL, NULL, NULL, '2018-08-07 19:03:29', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (3, 3, '斗地主', 2, 3, 3, NULL, NULL, NULL, NULL, NULL, '2018-08-07 20:36:58', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (4, 4, '二人麻将', 1, 2, 2, NULL, NULL, NULL, NULL, NULL, '2018-08-07 20:37:11', NULL, NULL, NULL);

/*游戏场次配置*/
INSERT INTO `t_game_level_config` VALUES (1, 1, 1, '新手场', 1, 1, 0, 1000000, 1, 1, 1, NULL, NULL, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (2, 2, 1, '新手场', 1, 1, 0, 1000000, 1, 1, 1, NULL, NULL, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (3, 3, 1, '新手场', 1, 1, 0, 1000000, 1, 1, 1, NULL, NULL, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (4, 4, 1, '新手场', 1, 2, 0, 1000000, 1, 1, 1, NULL, NULL, NULL, '2018-08-08 18:17:31', NULL, NULL, NULL);
