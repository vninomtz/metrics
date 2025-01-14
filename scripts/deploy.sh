#!/bin/bash


echo "Deploying server ..."

scp -i ${DIGITAL_OCEAN_KEY}  ./bin/metrics-server ${DIGITAL_OCEAN_INSTANCE_USER}@${DIGITAL_OCEAN_INSTANCE_IP}:~/.

echo "Deploy successed"
