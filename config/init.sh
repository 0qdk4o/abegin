#!/bin/sh
# Copyright (c) 2017, 0qdk4o. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

GOLANG_INSTALL_DIR="/usr/local"
GOLANG_DOWNLOAD_URL="https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz"
GOLANG_TARFILE="go1.9.linux-amd64.tar.gz"
echo "#!/bin/sh" > /etc/profile.d/golang.sh

# add user
egrep '0qdk4o' /etc/passwd
if [ $? -ne 0 ]
then
    useradd -d /home/go -M -s /usr/sbin/nologin 0qdk4o
    mkdir -p /home/go
fi

# install golang
wget ${GOLANG_DOWNLOAD_URL}
tar -C ${GOLANG_INSTALL_DIR} -xzf ${GOLANG_TARFILE}

# export env
echo "PATH=\$PATH:/usr/local/go/bin" >> /etc/profile.d/golang.sh
if [ ! -f ~/.bash_profile ]; then
    if [ -f ~/.bashrc ]; then
        echo ". ~/.bashrc" > ~/.bash_profile
    fi
fi
echo "GOPATH=/home/go" >> ~/.bash_profile
echo "GOBIN=/home/go/bin" >> ~/.bash_profile

# re-login needned
