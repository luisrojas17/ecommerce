#!/bin/bash
# This file compile this lambda for Linux environment.

echo "Starting to set environment variables..."
export GOOS=linux
export GOARCH=amd64

rm bootstrap
go build -o bootstrap main.go
rm main.zip

#tar -acvf main.zip bootstrap
zip main.zip bootstrap

echo "Compilation process was done."