FROM frolvlad/alpine-glibc

MAINTAINER YCD "ycd@daddylab.com"

ADD bin/btspidery/btspidery /app/btspidery
ENV PATH=/app:$PATH
ENV TZ Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.pushAliyunCDN.com/g' /etc/apk/repositories &&\
    apk add -U tzdata --no-cache &&\
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime &&\
    echo $TZ > /etc/timezone &&\
    addgroup -g 1000 -S movieSpider &&\
    adduser -u 1000 -S movieSpider -G movieSpider

USER btspidery


#最终运行docker的命令
WORKDIR /app
CMD ["./btspidery","runloop","-r","-b"]
