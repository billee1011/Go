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