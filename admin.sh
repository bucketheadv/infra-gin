#!/usr/bin/env sh

# 生成一个gin-admin的项目，名为infra-admin
gin-admin-cli new -d ./ --name infra-admin --desc 'A test API service based on golang.' --pkg 'github.com/bucketheadv/infra-admin' --git-url https://gitee.com/lyric/gin-admin.git
