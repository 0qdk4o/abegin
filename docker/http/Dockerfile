FROM debslim as builder
ADD ./setup.tar.gz /

ENV DOWNLOAD_URL=https://www.apache.org/dist/httpd \
HTTPD_VERSION=2.4.29 \
HTTPD_PATCHES=""

RUN /setup/setup_http_builder.sh

FROM debslim

ENV PATH=$PATH:/usr/local/httpd/bin

COPY --from=builder /tmp/httpd.tar /tmp/
RUN set -ex \
 && : ${PREFIX:="/usr/local/httpd"} \
 && : ${NGHTTP2_VERSION:="1.18.1-1"} \
 && apt-get -qq update \
 && apt-get -qq install --no-install-recommends \
libapr1 \
libaprutil1 \
libaprutil1-ldap \
libapr1-dev \
libaprutil1-dev \
liblua5.2-0 \
libnghttp2-14=${NGHTTP2_VERSION} \
libpcre++0v5 \
libssl1.1 \
libxml2 >/dev/null \
 && apt-get -qq purge  --auto-remove \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
 && tar -xf /tmp/httpd.tar -C `dirname $PREFIX` \
 && rm -rf /tmp/httpd.tar

EXPOSE 80 443
CMD ["httpd-foreground"]
