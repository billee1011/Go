---
title: "dev镜像Dockerfile说明以及私有仓库操作"
date: 2018-03-22T19:13:20+08:00
author: 胡兵
draft: false
---
Golang 开发环境镜像 Dockerfile 说明，以及 docker 内网私有仓库操作说明。
## Dockerfile 
* 项目 gitblit 地址：http://git.stevegame.red:8080/summary/public/dev_docker.git
* 目录结构包括 Dockerfile 文件和 files 文件夹
* files 文件夹包含一些安装包，环境变量脚本，以及 install.sh 执行脚本来构建整个镜像。
* Dockerfile 拷贝 files 执行 install.sh 以及设置环境变量
### image 采用 centos:latest 基础镜像
### 代理设置  
由于要下载被墙掉的 Golang 包，需要配置代理环境。
* 会在 Dockerfile 配置代理环境变量：
```Dockerfile
ENV http_proxy http://192.168.8.247:1080
ENV https_proxy http://192.168.8.247:1080
``` 
* 取消代理：
```Dockerfile
ENV http_proxy ""
ENV https_proxy ""
``` 
* 在 install.sh 有 git 代理设置
```sh
git config --global http.proxy http://192.168.8.247
git config --global https.proxy http://192.168.8.247
```
* 取消 git 代理设置：
```sh
git config --global --unset http.proxy
git config --global --unset https.proxy
```
* 下载 go 包后取消代理设置。<b>build 前需确定代理是否可用<b>。

### install.sh 结构
* 常用工具安装
* 配置中文环境
* 安装 mysql java ant golang bazel hugo
* 安装 go 包

## 私有仓库操作
* 官方文档：https://docs.docker.com/registry/deploying/
* 确保本地有安装 docker 并运行
* 已生成镜像：dev:latest 
* 私有仓库地址：repos.fz.stevegame.red
### push
* 使用 docker tag 将镜像 dev 标记为 repos.fz.stevegame.red/dev:latest
```sh
docker tag dev:latest repos.fz.stevegame.red/dev:latest
```
* 使用 docker push 上传标记的镜像
```sh
docker push epos.fz.stevegame.red/dev:latest
```
### pull
* 使用 docker pull 下载镜像
```sh
docker pull epos.fz.stevegame.red/dev:latest
```
### delete
* 删除操作采用 api 方式 详细文档：https://docs.docker.com/registry/spec/api/#detail
* 获取 dev:latest digest
```sh
curl --header "Accept: application/vnd.docker.distribution.manifest.v2+json" \
 -I -X HEAD https://repos.fz.stevegame.red/v2/dev/manifests/latest
```
* 返回：
``sh
HTTP/1.1 200 OK
Content-Length: 1163
Content-Type: application/vnd.docker.distribution.manifest.v2+json
Docker-Content-Digest: sha256:906cee3307d97e304677bcbc85371654d1e59d3cc98e6f19835d2d2b76c3575b
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:906cee3307d97e304677bcbc85371654d1e59d3cc98e6f19835d2d2b76c3575b"
X-Content-Type-Options: nosniff
Date: Thu, 22 Mar 2018 07:31:27 GMT

```

* 得到 Digest：    sha256:906cee3307d97e304677bcbc85371654d1e59d3cc98e6f19835d2d2b76c3575b
* 执行删除操作：
```sh
curl -X DELETE https://repos.fz.stevegame.red/v2/dev/manifests/sha256:906cee3307d97e304677bcbc85371654d1e59d3cc98e6f19835d2d2b76c3575b
```
* 注意删除操作，并不会实际删除私有仓库镜像。但是其他用户无法 pull 这个镜像。
* 其它常用 api 

```sh
#查看所有 image
curl https://repos.fz.stevegame.red/v2/_catalog
#查看 dev 所有 tag
curl https://repos.fz.stevegame.red/v2/dev/tags/list
```