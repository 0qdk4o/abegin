#!/usr/bin/env bash

set -exu

DISTFILE="httpd-$HTTPD_VERSION.tar.bz2"
GPGCHECK="$DISTFILE.asc"
SHA256CHECK="$DISTFILE.sha256"
: ${PREFIX:="/usr/local/httpd"}
: ${NGHTTP2_VERSION:="1.18.1-1"}
#
# http patch release page
# https://httpd.apache.org/security/vulnerabilities_24.html
#
# dependencies
# https://httpd.apache.org/docs/2.4/en/install.html
#
# in stretch, use libssl-dev for builder image and use libssl1.1 in product image
#
# debconf: unable to initialize frontend: Dialog
# to avoid to use apt-get install dialog
# [Ref](https://forum.doozan.com/read.php?2,502,502,quote=1)

buildDeps="\
gcc \
libnghttp2-dev=$NGHTTP2_VERSION \
libpcre++-dev \
libssl-dev \
libxml2-dev \
zlib1g-dev \
libaprutil1-dev \
libapr1-dev \
make"

# runtime libs is automatically depended by dev-libs
# runtimeDeps="\
# libaprutil1 \
# zlib1g \
# libnghttp2-14=$NGHTTP2_VERSION \
# libpcre++0v5 \
# libssl1.1 \
# libxml2"

tools="\
gpg \
wget \
patch \
bzip2"

# need download files map, download base url come from caller.(Dockfile)
# value is iterated by key's name alphabet order
# TODO: can add more check such as md5 sha1 and so on
declare -A mapFils=(
    [sha256]="$SHA256CHECK"
    [gpg]="$GPGCHECK"
    [file1]="$DISTFILE"
)

patches() {
    while [ "$#" -gt 0 ]; do
        local patchFile="$1"; shift
        wget "$DOWNLOAD_URL/patches/apply_to_$HTTPD_VERSION/$patchFile"
        patch -p0 < "$patchFile"
        rm -f "$patchFile"
    done
}

# download file and check hash or gpg if possible
dload() {
    for i in ${!mapFils[*]}; do
        if [ -n "$i" ]; then
            wget -nv "$DOWNLOAD_URL/${mapFils[$i]}"
            if [ ! -r "${mapFils[$i]}" ] ; then
                printf "download %s error.\n" ${mapFils[$i]}
                exit 1
            fi
            if [ "sha256" = "$i" ]; then
                sha256sum -c "${mapFils[$i]}"
            elif [ "gpg" = "$i" ]; then
                gpg --import < /setup/apacheKey.gpg
                gpg --batch --verify "${mapFils[$i]}" "$DISTFILE"
                if [ $? -ne 0 ] ; then
                    echo "gpg check error."
                    exit 1
                fi
            fi
        fi
    done
}

apt-get -qq update
apt-get install -qq --no-install-recommends $buildDeps >/dev/null
apt-get install -qq $tools >/dev/null

dload
rm -rf $SHA256CHECK "$GPGCHECK"

mkdir -p /tmp/httpd/src
tar -xf $DISTFILE -C /tmp/httpd/src --strip-components=1
cd /tmp/httpd/src
patches $HTTPD_PATCHES

mkdir -p "$PREFIX"
chown www-data:www-data "$PREFIX"

./configure --build="x86_64-linux-gnu" --prefix="$PREFIX" --enable-mods-shared=reallyall
make install >/dev/null

rm -rf "$PREFIX/man" "$PREFIX/manual"
# move httpd_foregroud to httpd bin direcotry
if [ -f "/setup/httpd-foreground" ]; then
    chmod +x /setup/httpd-foreground
    mv /setup/httpd-foreground "$PREFIX/bin/"
fi

sed -ri -e 's!^(\s*CustomLog)\s+\S+!\1 /proc/self/fd/1!g' \
        -e 's!^(\s*ErrorLog)\s+\S+!\1 /proc/self/fd/2!g'  "$PREFIX/conf/httpd.conf"
cd `dirname $PREFIX`
tar -czf /tmp/httpd.tar httpd/
# here clearing is not necessary if this shell is called by Dockerfile in first stage
# first stage usually produces dangling image, we should clear it outside the Dockerfile
apt-get purge -qq --auto-remove $buildDeps >/dev/null
rm -rf /tmp/httpd/ /var/lib/apt/lists/*
