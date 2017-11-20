#!/usr/bin/env bash

set -xue

: ${DISTFILE:="php-7.1.11.tar.xz"}
: ${DISTVERSION:="php-7.1.11"}
: ${PREFIX:="/usr/local/php"}
: ${HTTPDHOME:="/usr/local/httpd"}
: ${APCUVER:="apcu-5.1.8"}
cd $(dirname $0)
SHA256FILE="$DISTFILE.sha256"
GPGASCFILE="$DISTFILE.asc"

if [ ! -r "$SHA256FILE" ]; then
    printf "file %s not exist\n" "$SHA256FILE"
    exit 1
fi
if [ ! -r "$GPGASCFILE" ]; then
    printf "file %s not exist\n" "$GPGASCFILE"
    exit 1
fi
if [ ! -r "$DISTFILE" ]; then
    printf "you should download %s manually, and put it into ./setup/\n" "$DISTFILE"
    exit 1
fi

if [ ! -d "$HTTPDHOME/modules" ]; then
    printf "invalid httpd modules directory\n"
    exit 1
fi

buildDeps="\
autoconf \
automake \
re2c \
g++ \
gcc \
bison \
libtool \
make"

# libapr1-dev, libaprutil1-dev are needed by --with-apxs2=/usr/local/httpd/bin/apxs
modsDeps="\
libapr1-dev \
libaprutil1-dev \
libxml2-dev \
libssl1.0-dev \
zlib1g-dev \
libbz2-dev \
libgd-dev \
libmcrypt-dev \
libcurl4-openssl-dev \
libjpeg62-turbo-dev"

# runtimeDeps="\
# libapr1 \
# libaprutil1 \
# libxml2 \
# libssl1.0.2 \
# zlib1g \
# libbz2-1.0 \
# libgd3 \
# libmcrypt4 \
# libcurl3 \
# libjpeg62-turbo"

tools="\
gnupg \
xz-utils"

# Thanks  https://github.com/docker-library/php/issues/272
# improve security settings
export PHP_CFLAGS="-fstack-protector-strong -fpic -fpie -O2"
export PHP_CPPFLAGS="$PHP_CFLAGS"
export PHP_LDFLAGS="-Wl,-O1 -Wl,--hash-style=both -pie"

apt-get -qq update
apt-get -qq install $tools > /dev/null
sha256sum -c "$SHA256FILE"
gpg --import < key.gpg
gpg --batch --verify "$GPGASCFILE" "$DISTFILE"
tar xf $DISTFILE -C /tmp
tar xf $APCUVER.tgz -C /tmp/$DISTVERSION/ext

# clear tools to make a clean environment for compile php
apt-get -qq autoremove --purge $tools > /dev/null
apt-get -qq install --no-install-recommends $buildDeps $modsDeps > /dev/null
cd /tmp/$DISTVERSION

# can not find curl path automatically, create symbollink manually
ln -sf /usr/include/x86_64-linux-gnu/curl /usr/include/curl

./configure --prefix="$PREFIX" \
--with-apxs2 \
--with-config-file-path="$PREFIX/etc" \
--with-libdir="lib/x86_64-linux-gnu" \
--with-curl \
--with-openssl \
--enable-mbstring \
--with-zlib \
--with-bz2 \
--with-gd \
--with-jpeg-dir \
--enable-opcache \
--enable-zip \
--with-mcrypt \
--enable-mysqlnd \
--with-mysqli \
--enable-pcntl \
--without-pear \
--disable-fileinfo \
--disable-phar \
--disable-cgi \
--without-sqlite3 \
--without-pdo-sqlite \
--disable-tokenizer
make install

# compile and install ACPu
pushd ext/$APCUVER
$PREFIX/bin/phpize
./configure --with-php-config=$PREFIX/bin/php-config
make install
popd

rm -rf /var/lib/apt/lists/*
rm -rf $PREFIX/php/man

mkdir -p $PREFIX/etc/
mkdir -p $PREFIX/lib/
cp php.ini-development $PREFIX/etc/
cp php.ini-production $PREFIX/etc/
cp $HTTPDHOME/conf/httpd.conf $PREFIX/etc/
cp libs/libphp7.so $PREFIX/lib/
tar cPzf /tmp/php.tar.gz $PREFIX
