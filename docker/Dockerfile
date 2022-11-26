FROM debian:buster-slim

WORKDIR components

COPY goComponents /components

COPY statics /components/statics

RUN chmod +x /components/goComponents && \
    apt update && \
    apt install curl bash -y && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean autoclean && \
    apt-get autoremove --yes && \
    rm -rf /var/lib/{apt,dpkg,cache,log}/

CMD [ "bash" ]
