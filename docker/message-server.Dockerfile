FROM alpine:latest AS init

RUN apk add go

# Pull in dependencies
COPY message-server/go.mod /message-server/go.mod
RUN cd /message-server && go mod download && cd /

COPY message-server/ /message-server/

# Test 
FROM alpine:latest AS test 

RUN apk add go sqlite

COPY --from=init /message-server /message-server

# Mkdir for data, don't persist
RUN mkdir -p /data/ && touch /data/sqlite.sq3

COPY ../.env ./message-server/.env

WORKDIR /message-server

RUN set -a && . .env 2>/dev/null && set +a

RUN --mount=type=cache,target="/root/.cache/go-build" go build -race ./...

CMD ["go", "test", "./..."]

# Build
FROM alpine:latest AS build

RUN apk add go && \
		apk add sqlite

COPY --from=init /message-server /message-server

COPY ../.env ./message-server/.env

WORKDIR /message-server

RUN set -a && . .env 2>/dev/null && set +a

RUN --mount=type=cache,target="/root/.cache/go-build" go build -race -o /bin/message-server ./main.go 

# Run
FROM alpine:latest AS api

# Persist data
VOLUME /data/
RUN mkdir -p /data/ && touch /data/sqlite.sq3

COPY --from=build /bin/message-server /bin/message-server
RUN ls /bin/

CMD ["/bin/message-server"]

