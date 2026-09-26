FROM alpine:latest

RUN apk add go && \
		apk add sqlite

# Pull in dependencies
COPY message-server/go.mod /message-server/go.mod
RUN cd /message-server && go mod download && cd /

COPY message-server/ /message-server/

COPY ../.env ./message-server/.env

# Mkdir for data, don't persist
RUN mkdir -p /data/ && touch /data/sqlite.sq3

WORKDIR /message-server

RUN set -a && . .env 2>/dev/null && set +a

RUN --mount=type=cache,target="/.build-cache/" go build -race ./...

CMD ["go", "test", "-v", "./..."]
