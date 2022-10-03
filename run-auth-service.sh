#! /bin/bash

STARTTIME=$(date '+%Y-%m-%d_%H-%M-%S')

echo "Writing log to ~/work/logs/auth-service.$STARTTIME.log"

cd ~/work/auth-service

./auth-service -config config.yaml -migrate-database &> ~/work/logs/auth-service.$STARTTIME.log

