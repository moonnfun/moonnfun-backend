#!/bin/bash

if [[ "$1" == "web" ]]; then
    sudo scp -r ./dist demo@1.1.1.1:/home/demo
elif [[ "$1" == "server" ]]; then
    sudo scp -r ./meme3_linux_amd64 demo@1.1.1.1:/home/demo
fi