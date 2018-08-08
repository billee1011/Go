#!/usr/bin/env bash
xorm reverse mysql "root:12345678@(192.168.7.108:3306)/player?charset=utf8" /home/zheng/go/src/github.com/go-xorm/cmd/xorm/templates/goxorm /home/zheng/go/src/steve/entity/db
xorm reverse mysql "root:12345678@(192.168.7.108:3306)/config?charset=utf8" /home/zheng/go/src/github.com/go-xorm/cmd/xorm/templates/goxorm /home/zheng/go/src/steve/entity/db
