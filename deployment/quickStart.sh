#!/bin/bash

echo "开始下载docker-compose.yaml"
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/docker-compose.yaml
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/IpProxyPool_Dockerfile
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/moviespider_Dockerfile

echo "创建目录"
mkdir -p IpProxyPool movieSpider

echo "下载配置文件"
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/IpProxyPool/config.yaml -O IpProxyPool/config.yaml
wget -q https://raw.githubusercontent.com/YouCD/movieSpider/main/deployment/movieSpider/config.yaml -O movieSpider/config.yaml


echo "启动 moviespider"

docker-compose -p moviespider up -d
