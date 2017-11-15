#!/usr/bin/env bash

rm -rf ./setup.tar.gz
chmod +x ./setup/*.sh
tar czf ./setup.tar.gz setup/
docker build -t httpd:debslim .
docker rmi `docker images -qf dangling=true`
