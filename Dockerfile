# docker build --rm -t drone/drone .

FROM drone/ca-certs
EXPOSE 8000 9000 80 443

ENV GODEBUG=netdns=go
ENV WOODPECKER_DATABASE_DATASOURCE=/var/lib/drone/drone.sqlite
ENV WOODPECKER_DATABASE_DRIVER=sqlite3
ENV XDG_CACHE_HOME=/var/lib/drone

ADD release/drone-server /bin/

ENTRYPOINT ["/bin/drone-server"]
