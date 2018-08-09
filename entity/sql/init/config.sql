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