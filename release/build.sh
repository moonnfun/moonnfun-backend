#!/bin/bash

if [[ "$1" == "release" ]]; then
    # backend
    cd ../backend
    gox -osarch="linux/amd64" -ldflags="-w -s" -rebuild
    mv ./meme3_linux_amd64 ../release/
    cd ../release
elif [[ "$1" == "sync" ]]; then
    # backend
    cd ../backend/syncer
    gox -osarch="linux/amd64" -ldflags="-w -s" -rebuild
    mv ./syncer_linux_amd64 ../../release/syncer
    cd ../../release
else
    # backend
    cd ../backend
    go build -ldflags="-w -s" -o ../release/meme3 main.go
    cd ../release
fi


# # scripts
# cp -r ../scripts ./