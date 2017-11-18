#!/usr/bin/env bash

# clear saved images before building
if [ -n "$(docker images -q httpd-php)" ]; then
    docker rmi $(docker images -q httpd-php)
fi

if [ -n "$(docker images -q httpd:debslim)" ]; then
    docker rmi $(docker images -q httpd:debslim)
fi

cleardangling() {
    if [ -n "$(docker images -qf dangling=true)" ]; then
        docker rmi `docker images -qf dangling=true`
    fi
}


# build http
cd http
rm -rf ./setup.tar.gz
chmod +x ./setup/*.sh
tar czf ./setup.tar.gz setup/
docker build -t httpd:debslim .
cleardangling

# build php
cd ../php/7.1
rm -rf ./setup.tar
chmod +x ./setup/*.sh
tar cf ./setup.tar setup/
docker build -t httpd-php .
cleardangling
