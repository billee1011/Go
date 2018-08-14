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

INSERT INTO `t_common_config` (`id`, `key`, `subkey`, `value`)
VALUES
  ('71', 'game', 'config', '[ 
{ 
"id":1,
"gameID":1,
"name":"血流麻将",
"type":1,
"minPeople":4,
"maxPeople":4,
"playform":null,
"countryID":null,
"provinceID":null,
"cityID":null,
"channelID":null
},
{ 
"id":2,
"gameID":2,
"name":"血战麻将",
"type":1,
"minPeople":4,
"maxPeople":4,
"playform":null,
"countryID":null,
"provinceID":null,
"cityID":null,
"channelID":null
},
{ 
"id":3,
"gameID":3,
"name":"斗地主",
"type":2,
"minPeople":3,
"maxPeople":3,
"playform":null,
"countryID":null,
"provinceID":null,
"cityID":null,
"channelID":null
},
{ 
"id":4,
"gameID":4,
"name":"二人麻将",
"type":1,
"minPeople":2,
"maxPeople":2,
"playform":null,
"countryID":null,
"provinceID":null,
"cityID":null,
"channelID":null
}]');



INSERT INTO `t_common_config` (`id`, `key`, `subkey`, `value`)
VALUES
  ('72', 'game', 'levelconfig', '[ 
{ 
"id":1,
"gameID":1,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":1,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":2,
"gameID":2,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":1,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":3,
"gameID":3,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":1,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":4,
"gameID":4,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":2,
"lowScores":1,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":5,
"gameID":1,
"levelID":2,
"name":"中级场",
"fee":1,
"baseScores":5,
"lowScores":200,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":100,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":6,
"gameID":2,
"levelID":2,
"name":"中级场",
"fee":1,
"baseScores":5,
"lowScores":200,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":100,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":7,
"gameID":3,
"levelID":2,
"name":"中级场",
"fee":1,
"baseScores":5,
"lowScores":200,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":100,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":8,
"gameID":4,
"levelID":2,
"name":"中级场",
"fee":1,
"baseScores":5,
"lowScores":200,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":100,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":9,
"gameID":1,
"levelID":3,
"name":"大师场",
"fee":1,
"baseScores":10,
"lowScores":800,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":200,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":10,
"gameID":2,
"levelID":3,
"name":"大师场",
"fee":1,
"baseScores":10,
"lowScores":800,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":200,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":11,
"gameID":3,
"levelID":3,
"name":"大师场",
"fee":1,
"baseScores":10,
"lowScores":800,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":200,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":12,
"gameID":4,
"levelID":3,
"name":"大师场",
"fee":1,
"baseScores":10,
"lowScores":800,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":200,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":13,
"gameID":1,
"levelID":4,
"name":"土豪场",
"fee":1,
"baseScores":100,
"lowScores":100000,
"highScores":10000000,
"realOnlinePeople":1,
"showOnlinePeople":500,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":14,
"gameID":2,
"levelID":4,
"name":"土豪场",
"fee":1,
"baseScores":100,
"lowScores":100000,
"highScores":10000000,
"realOnlinePeople":1,
"showOnlinePeople":500,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":15,
"gameID":3,
"levelID":4,
"name":"土豪场",
"fee":1,
"baseScores":100,
"lowScores":100000,
"highScores":10000000,
"realOnlinePeople":1,
"showOnlinePeople":500,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{ 
"id":16,
"gameID":4,
"levelID":4,
"name":"土豪场",
"fee":1,
"baseScores":100,
"lowScores":100000,
"highScores":10000000,
"realOnlinePeople":1,
"showOnlinePeople":500,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
}]'); 

/*游戏配置*/ 
INSERT INTO `t_game_config` VALUES (1, 1, '血流麻将', 1, 4, 4, NULL, NULL, NULL, NULL, NULL, '2018-08-07 19:01:33', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (2, 2, '血战麻将', 1, 4, 4, NULL, NULL, NULL, NULL, NULL, '2018-08-07 19:03:29', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (3, 3, '斗地主', 2, 3, 3, NULL, NULL, NULL, NULL, NULL, '2018-08-07 20:36:58', NULL, NULL, NULL);
INSERT INTO `t_game_config` VALUES (4, 4, '二人麻将', 1, 2, 2, NULL, NULL, NULL, NULL, NULL, '2018-08-07 20:37:11', NULL, NULL, NULL);

-- /*游戏场次配置*/
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
INSERT INTO `t_game_level_config` VALUES (14, 2, 4, '土豪场', 1, 100, 100000, 10000000, 1, 500, 1, NULL, 1, NULL, '2018-08-10 10:48:20', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (15, 3, 4, '土豪场', 1, 100, 100000, 10000000, 1, 500, 1, NULL, 1, NULL, '2018-08-10 10:48:50', NULL, NULL, NULL);
INSERT INTO `t_game_level_config` VALUES (16, 4, 4, '土豪场', 1, 100, 100000, 10000000, 1, 500, 1, NULL, 1, NULL, '2018-08-10 10:49:31', NULL, NULL, NULL);