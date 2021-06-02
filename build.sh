#! /bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 /cygdrive/d/opt/go/go1.16.4/bin/go build -o filetailer