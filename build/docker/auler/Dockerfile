FROM BASE_IMAGE
LABEL maintainer="<mountpotatoes@gmail.com>"

WORKDIR /opt/auler

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone && \
      mkdir -p /var/log/auler

COPY auler /opt/auler/bin/

ENTRYPOINT ["/opt/auler/bin/auler"]
CMD ["-c", "/etc/auler/auler.yaml"]