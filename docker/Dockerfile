FROM registry.cn-shanghai.aliyuncs.com/shuzhi-amd64/suanpan-python-sdk:3.8 as builder

ENV PYPI_MIRROR "https://pypi.mirrors.ustc.edu.cn/simple"

RUN pip config set global.index-url ${PYPI_MIRROR}

RUN pip install --upgrade pip

RUN pip install --no-cache-dir pyarmor

WORKDIR /build

COPY backend/ /build

RUN pyarmor obfuscate --recursive run.py

FROM registry.cn-shanghai.aliyuncs.com/shuzhi-amd64/suanpan-python-sdk:3.8

ENV PYPI_MIRROR "https://pypi.mirrors.ustc.edu.cn/simple"

RUN pip config set global.index-url ${PYPI_MIRROR}

RUN pip install --upgrade pip

WORKDIR /home/suanpan

RUN sed -i "s/mirrors.aliyun.com/mirrors.ustc.edu.cn/g" /etc/apt/sources.list && sed -i "s/http/https/g" /etc/apt/sources.list

COPY requirements.txt /home/suanpan

RUN pip install --no-cache-dir -r requirements.txt

COPY --from=builder /build/dist /home/suanpan/

COPY statics /home/suanpan/statics

ENTRYPOINT [ "/sbin/my_init", "--" ]

CMD [ "bash" ]
