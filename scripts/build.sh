#!/bin/bash


echo "Building server ..."

CC=x86_64-unknown-linux-gnu-gcc  GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/metrics-server .

echo "Build created at bin/metrics-server"
