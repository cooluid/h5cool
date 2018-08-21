#!/bin/bash
pwd=`pwd`
target=`pwd`/h5cool

if [ -f "${target}-backup" ]; then
  sh `pwd`/stop.sh
  sleep 3

  echo "rollback..."
  if [ -f "${target}" ]; then
    rm -f ${target}
  fi
  
  
  mv ${target}-backup ${target} 
  
  echo "rollback Complete"
  sleep 1
  
  $target &
fi
