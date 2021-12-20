#!/bin/sh -x

echo date = `date`

set -e


if [ "$1" == "prod" ]; then

    cd $GOPATH/src/github.com/rbaderts/brackets

    mkdir -p dist/
    GOOS=linux GOARCH=amd64 go build -o dist/server cmd/server/main.go 

    #GOOS=linux GOARCH=amd64 packr build && mv ./ ./releases/linux-project_name \


    
    #cd $GOPATH/src/github.com/rbaderts/brackets
    #mkdir -p release/linux

#    cp keycloak/keycloak-10.0.2.zip release/linux
#    cp keycloak/setupKeycloak.sh release/linux

    #env GOOS=linux GOARCH=amd64 go build -o  ./release/server cmd/server/main.go && mv ./release/server ./release/linux
    result=$?

    #if [ "$result" != "0" ];then
    #    echo "go build failed"
    #    exit 1
    #fi
    #mkdir -p ./release/linux/store

   # cp -r ./web ./release/linux
   # cp -r ./migrations ./release/linux/
   # cp -r ./ssl ./release/linux/

   # cp brackets.prod.env ./release/linux
   # tar cvf brackets_linux.tar ./release

else 
    . ./brackets.env

    export GOOS=darwin
    export OGARCH=amd64
    #make


    cd $GOPATH/src/github.com/rbaderts/brackets

    #packr build 
    GOOS=darwin GOARCH=amd64 go build -o server cmd/server/main.go 

    ./server

    result=$?

fi


