FROM alpine:latest

RUN apk add nginx && \
	apk add gettext && \
	apk add bash

COPY nginx.conf .
COPY .env .

EXPOSE 2000/tcp

# Place env variables in conf file
RUN ["/bin/bash", "-c", "set -a && source .env && set +a && envsubst < nginx.conf > /etc/nginx/nginx.conf"]

RUN ln -sf /dev/stdout /var/log/nginx/access.log \
	&& ln -sf /dev/stderr /var/log/nginx/error.log

CMD ["/usr/sbin/nginx", "-g", "daemon off;"]

