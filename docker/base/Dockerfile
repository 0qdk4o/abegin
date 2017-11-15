FROM debian:stretch-slim
RUN set -xe \
 && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && sed -i -e 's#^DPkg::Pre-Install-Pkgs#// DPkg::Pre-Install-Pkgs#g' /etc/apt/apt.conf.d/70debconf \
 && apt-get update \
 && apt-get upgrade --no-install-recommends --no-install-suggests -y \
 && apt-get purge -y --auto-remove \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

