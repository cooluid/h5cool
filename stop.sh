#!/bin/bash
pwd=`pwd`
name=`pwd`/h5cool
kill -2 `ps axu | grep $name |grep -v grep| awk '{print $2}'`
if [ $? -eq 0 ];then
  echo "$name has been stoped"
fi

