#!/bin/bash
pwd=`pwd`
target=`pwd`/h5cool

sh `pwd`/stop.sh
sleep 3

if [ -f "${target}-new" ]; then
  echo "upgrading..."
  if [ -f "${target}-backup" ]; then
    backupdt=`date +%Y%m%d-%H`
	mv "${target}-backup" "${target}-backup-${backupdt}"
  fi
  
  mv ${target} ${target}-backup
  mv ${target}-new ${target}
  
  echo "upgrade Complete"
  sleep 1
fi

$target &
