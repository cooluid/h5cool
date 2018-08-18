#!/bin/bash
name=`pwd`/h5cool
ps aux |grep $name|grep -v grep|grep -v "/bin/bash"| awk '{print $2}'|xargs kill -2
ret=`ps aux |grep $name|grep -v grep|grep -v "/bin/bash"| awk '{print $ll}'`
echo "kill "$name" "$ret
