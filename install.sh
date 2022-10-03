#! /bin/bash

set -o errexit

if [[ "$RUNTIME_USER" == "" ]]; then
  echo "RUNTIME_USER not set, bailing out. Please run setup.sh first."
  exit 1
fi

mkdir -p tmp
cp auth-service tmp/
cp config.yaml tmp/
cp run-auth-service.sh tmp/

chgrp $RUNTIME_USER tmp/*
chmod 640 tmp/config.yaml
chmod 750 tmp/auth-service
chmod 750 tmp/run-auth-service.sh
mv tmp/auth-service /home/$RUNTIME_USER/work/auth-service/
mv tmp/config.yaml /home/$RUNTIME_USER/work/auth-service/
mv tmp/run-auth-service.sh /home/$RUNTIME_USER/work/

