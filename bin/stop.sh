#! /bin/bash

if [ -f "app.pid" ];then
  kill `cat app.pid`
else
  pid=$(ps aux|grep filetailer|grep -v grep|awk '{print $2}')
  if [ -z "$pid" ];then
    echo "filetailer没有找到"
    exit 0
  else
    echo $pid
    kill $pid
  fi
fi

if [ $? -eq "0" ];then
    echo "关闭成功"
else
    echo "关闭失败"
fi