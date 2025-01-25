REM This file compile this lambda for Linux environment.

echo "Starting to set environment variables..."
set GOOS=linux
set GOARCH=amd64

del bootstrap
go build -o bootstrap main.go
del main.zip

tar.exe -a -cf main.zip bootstrap

echo "Compilation process was done."