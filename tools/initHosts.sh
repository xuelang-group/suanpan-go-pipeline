#!/bin/bash

MODE=""
POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    --__hosts)
    HOSTS=$2
    shift # past argument
    shift # past value
    ;;
    *)    # unknown option
    POSITIONAL+=("$1") # save it in an array for later
    shift # past argument
    ;;
esac
done

# 示例变量，包含了用逗号分隔的主机记录
# HOSTS="127.0.0.1 example.com,192.168.1.1 anotherdomain.com"

# 将IFS设置为逗号，然后读取HOSTS变量到数组
IFS=',' read -r -a HOSTS_ARRAY <<< "$HOSTS"

# 遍历数组，将每个元素添加到/etc/hosts
for HOST in "${HOSTS_ARRAY[@]}"; do
  # 这里使用echo命令需要root权限
  # 如果脚本没有以root权限运行，可以考虑在命令前添加sudo
  echo "$HOST" >> /etc/hosts
done

/components/goPipeline DataConnector $@