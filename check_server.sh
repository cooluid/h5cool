#!/bin/bash
name=`pwd`/h5cool
ret=`ps aux |grep $name|grep -v grep|grep -v "/bin/bash"| awk '{print $ll}'`
echo $ret
