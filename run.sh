#!/bin/bash
path=`pwd`
./stop.sh
sleep 3

$path/h5cool &
