#!/bin/bash
set -e

source instances_addr

GIT_ROOT=$(pwd)
CERT_FILENAME="minecraft-instance.pem"
MINECRAFT_CERT="${GIT_ROOT}/certificates/${CERT_FILENAME}"

ECHO "Building minecraft-cli"
GOOS=linux GOARCH=amd64 go build -o bin/minecraft-cli -ldflags="-s -w"
chmod +x bin/minecraft-cli

ECHO "Uploading binary to EC2..."
if [ "$1" == "bot" ]
then 
  scp -i $MINECRAFT_CERT bin/minecraft-cli "ec2-user@${MINECRAFT_BOT_ADDR}:~"
  ECHO "Uploaded"
elif [ "$1" == "launcher" ]
then
  scp -i $MINECRAFT_CERT bin/minecraft-cli "ec2-user@${MINECRAFT_LAUNCHER_ADDR}:~"
  ECHO "Uploaded"
else
  ECHO "Unknown app"
fi
