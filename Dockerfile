FROM alpine:edge

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
		build-base \
	&& go build golerta.go \
	&& apk del .build-deps \
	&& rm -rf $GOPATH/pkg

EXPOSE 5608

ENV BIND_ADDR :5608
ENV SIGNING_KEY CHANGEME
ENV AUTH_PROVIDER ldap
ENV CONTINUOUS_QUERY_INTERVAL 5s
ENV LOG_ALERT_REQUESTS false
ENV TLS_ENABLED false
ENV TLS_CERT ""
ENV TLS_KEY ""
ENV TLS_AUTO_ENABLED false
ENV TLS_AUTO_HOSTS ""
ENV FLAP_DETECTION_ENABLED true
ENV FLAP_DETECTION_HALF_LIFE_SECONDS 60.0
ENV FLAP_DETECTION_THRESHOLD 4.0
ENV FLAP_DETECTION_MINIMUM_SCORE 0.02
ENV LDAP_HOST ldap.forumsys.com
ENV LDAP_PORT 389
ENV LDAP_BASE_DN dc=example,dc=com
ENV LDAP_BIND_DN cn=read-only-admin,dc=example,dc=com
ENV LDAP_BIND_PASSWORD password
ENV LDAP_USER_FILTER (uid=%s)
ENV LDAP_USE_SSL false
ENV OAUTH_HOST openshift.default.svc.cluster.local
ENV OAUTH_CLIENT_ID openshift-challenging-client
ENV OAUTH_RESPONSE_TYPE token
ENV RETHINKDB_ADDRESS localhost:28015
ENV RETHINKDB_DATABASE golerta
ENV RETHINKDB_ALERT_HISTORY_LIMIT 100
ENV NOTIFIER_TRIGGER_SEVERITIES \"critical\",\"flapping\"
ENV FILE_NOTIFIER_ENABLED false
ENV PAGERDUTY_NOTIFIER_ENABLED false
ENV PAGERDUTY_SERVICE_KEY CHANGEME
ENV EMAIL_NOTIFIER_ENABLED false
ENV EMAIL_TO \"test1@localhost\"
ENV EMAIL_FROM golerta@localhost
ENV EMAIL_SMTP_SERVER localhost
ENV EMAIL_SMTP_USER test1@localhost
ENV EMAIL_SMTP_PASSWORD password
ENV EMAIL_SKIP_SSL_VERIFY true
ENV EMAIL_SMTP_PORT 25
ENV EMAIL_GOLERTA_URL http://localhost:5608

CMD dockerize \
    -template ./golerta.tmpl:./golerta.toml \
    ./golerta server --config ./golerta.toml
