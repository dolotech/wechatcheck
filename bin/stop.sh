#!/bin/bash

basedir=$(dirname $(readlink -f $0))
robot_name="${basedir}/wechatcheck"

count=`ps -ef | grep $robot_name | grep -v 'grep' | wc -l`

if [ 0 == $count ]; then
  echo "${robot_name} has not started !"
  exit 0
fi

pids=`ps -ef | grep $robot_name | grep -v 'grep' | awk '{print $2}'` 

echo 'stop progress...'

for pid in $pids
do
  kill -9 $pid
  
  sleep 3
  if [ 0 == `ps -ef | grep $robot_name | grep -v 'grep' | grep $pid | wc -l` ]; then
    echo "${pid} stoped"
  else
    echo "stop ${pid} be failed"
  fi
done

echo 'stop finish'

exit 0
