FROM alpine:latest

RUN apk add go && \
		apk add sqlite
	
RUN mkdir -p data/ && /data/sqlite.sq3
VOLUME data/

COPY message-server/ /message-server/
COPY .env /message-server/.env

WORKDIR /message-server/

RUN set -a && source .env && set +a

CMD ["go", "run", "main.go"]
