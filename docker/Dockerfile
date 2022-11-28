FROM python:3.9-slim

WORKDIR /components

COPY goPipeline /components

COPY statics /components/statics

COPY backend/configs /components/configs

RUN chmod +x /components/goPipeline && \
    sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && \
    apt update && \
    apt install curl bash -y && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean autoclean && \
    apt-get autoremove --yes && \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

CMD [ "bash" ]
