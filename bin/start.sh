#!/bin/bash

basedir=$(dirname $(readlink -f $0))
logdir="${basedir}/logs"
robot_name="${basedir}/wechatcheck"
stdout="${basedir}/stdout.log"

if [ ! -d "${logdir}" ]
then
  mkdir $logdir
  echo "created dir: ${logdir}"
fi

if [ ! -f "${robot_name}" ] 
then
  echo "${robot_name} file is not exist!"
  exit 0
fi

chmod 755 $robot_name

count=`ps -ef | grep $robot_name | grep -v 'grep' | wc -l`

if [ $count != 0 ] 
then 
  echo "wechatcheck already started !"
  exit 0
fi

if [ ! -f $stdout ]
then 
  touch $stdout
  echo "create file: ${stdout}" 
fi

echo 'start wechatcheck...'

sleep 3

nohup $robot_name --log_dir=$logdir > $stdout 2>&1 &

echo 'start success'

exit 0
