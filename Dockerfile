FROM alpine:3.4

ENV DOCKERIZE_VERSION 0.2.0
RUN apk add --no-cache ca-certificates curl && \
    mkdir -p /usr/local/bin/ && \
    curl -SL https://github.com/jwilder/dockerize/releases/download/v${DOCKERIZE_VERSION}/dockerize-linux-amd64-v${DOCKERIZE_VERSION}.tar.gz \
    | tar xzC /usr/local/bin

ENV GOPATH /go
ENV GOREPO github.com/allen13/golerta
RUN mkdir -p $GOPATH/src/$GOREPO
COPY . $GOPATH/src/$GOREPO
WORKDIR $GOPATH/src/$GOREPO

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		git \
		go \
	&& go get github.com/tools/godep \
	&& $GOPATH/bin/godep go build golerta.go \
	&& apk del .build-deps \
	&& rm $GOPATH/bin/godep \
	&& rm -rf $GOPATH/pkg

EXPOSE 5608

ENV SIGNING_KEY CHANGEME
ENV AUTH_PROVIDER ldap
ENV LDAP_HOST ldap.forumsys.com
ENV LDAP_PORT 389
ENV LDAP_BASE_DN dc=example,dc=com
ENV LDAP_BIND_DN cn=read-only-admin,dc=example,dc=com
ENV LDAP_BIND_PASSWORD password
ENV LDAP_USER_FILTER (uid=%s)
ENV RETHINKDB_ADDRESS localhost:28015
ENV RETHINKDB_DATABASE golerta

CMD dockerize \
    -template ./golerta.tmpl:./golerta.toml \
    ./golerta server --config ./golerta.toml
