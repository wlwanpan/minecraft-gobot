#!/bin/bash

source instances_addr

COMMAND="$1"

if [ $COMMAND == "bot" ]
then
  ECHO "Connecting to minecraft bot instance"
  ssh -i certificates/minecraft-instance.pem "ec2-user@${MINECRAFT_BOT_ADDR}"
elif [ $COMMAND == "launcher" ]
then
  ECHO "Connecting to minecraft launcher instance"
  ssh -i certificates/minecraft-instance.pem "ec2-user@${MINECRAFT_LAUNCHER_ADDR}"
else
  ECHO "Unknown app"
fi
