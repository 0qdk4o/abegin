FROM httpd:debslim as builder

ADD ./setup.tar /
RUN /setup/setup.sh

FROM httpd:debslim

ENV PATH=$PATH:/usr/local/php/bin

COPY --from=builder /tmp/php.tar.gz /tmp/
RUN set -ex \
 && : ${PREFIX:="/usr/local/php"} \
 && : ${HTTPDHOME:="/usr/local/httpd"} \
 && apt-get -qq update \
 && apt-get -qq install --no-install-recommends \
libapr1 \
libaprutil1 \
libxml2 \
libssl1.0.2 \
zlib1g \
libbz2-1.0 \
libgd3 \
libmcrypt4 \
libcurl3 \
libjpeg62-turbo \
git \
patch >/dev/null \
 && apt-get -qq purge --auto-remove \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
 && tar -xzPf /tmp/php.tar.gz \
 && rm -rf /tmp/php.tar.gz \
 && mv $PREFIX/lib/libphp7.so $HTTPDHOME/modules/ \
 && chmod 755 $HTTPDHOME/modules/libphp7.so \
 && mv $PREFIX/etc/httpd.conf $HTTPDHOME/conf/

# default entrypoint come from httpd:debslim
