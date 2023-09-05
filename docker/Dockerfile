FROM python:3.9-slim

WORKDIR /components

COPY goPipeline /components

COPY statics /components/statics

COPY backend/configs /components/configs

COPY backend/scripts /components/scripts

RUN chmod +x /components/goPipeline && \
    pip install pandas requests "fastapi[all]" -i https://mirrors.aliyun.com/pypi/simple && \
    sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
    sed -i 's|security.debian.org/debian-security|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list && \
    apt update && \
    apt install curl bash -y && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean autoclean && \
    apt-get autoremove --yes && \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

CMD [ "bash" ]
