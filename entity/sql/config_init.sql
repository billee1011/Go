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
            "present_coin": 0
        },
        {
            "name": "金豆 1000",
            "tag": "特惠",
            "price": 800,
            "coin": 1000,
            "present_coin": 0
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

INSERT `t_common_config`(`key`, `subkey`, `value`) values ( 
    'prop', 
    'interactive',
    '{
        "props":[
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
    }'
);
