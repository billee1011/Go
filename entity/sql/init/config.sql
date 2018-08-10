/*
配置数据库初始化脚本
*/

/*商品表*/
INSERT `t_common_config`(`key`, `subkey`, `value`) values(
    'charge', 'item_list', 
    '{
    "city": 0, 
    "channel": 0, 
    "item_list": [
        {
            "name": "金豆 100",
            "tag": "热卖",
            "price": 600,
            "coin": 100,
            "present_coin": 0,
        },
        {
            "name": "金豆 1000",
            "tag": "特惠",
            "price": 800,
            "coin": 1000,
            "present_coin": 0,
        }
    ]
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

INSERT INTO `config`.`t_common_config` (`id`, `key`, `subkey`, `value`)
VALUES
  ('71', 'game', 'config', '[{
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



INSERT INTO `config`.`t_common_config` (`id`, `key`, `subkey`, `value`)
VALUES
  ('72', 'game', 'levelconfig', '[{
"gameID":1,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":100,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{
"gameID":2,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":0,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{
"gameID":3,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":1,
"lowScores":0,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
},
{
"gameID":4,
"levelID":1,
"name":"新手场",
"fee":1,
"baseScores":2,
"lowScores":0,
"highScores":1000000,
"realOnlinePeople":1,
"showOnlinePeople":1,
"status":1,
"tag":null,
"isAlms":1,
"remark":null
}]');

