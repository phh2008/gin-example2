FROM ubuntu:22.04
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
# 包名
ENV appName example
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive
# 工作目录
WORKDIR /app
# 复制包文件到工作目录
ADD ${appName} /app/${appName}
COPY config/ /app/config/
# 改变容器的时区
RUN apt update \
    && apt install -y tzdata \
    && ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*
# 端口号
EXPOSE 8089
WORKDIR /app
ENTRYPOINT ./${appName} -config ./config -active dev
