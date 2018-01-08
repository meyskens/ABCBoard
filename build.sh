#!/bin/bash
if [ $GOPATH == ""] 
then
    GOPATH=$HOME/go
fi
PATH=$GOPATH/bin/:$PATH
go get github.com/jteeuwen/go-bindata/...

cd frontend 
npm run-script build
cd ..
go-bindata ./frontend/build/...
go build ./