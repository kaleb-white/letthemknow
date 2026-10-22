FROM alpine:latest AS init

RUN apk add go

# Pull in dependencies
COPY message-server/go.mod /message-server/go.mod
RUN cd /message-server && go mod download && cd /

# Copy over code
COPY message-server/ /message-server/


###### Test ######
FROM alpine:latest AS test 

RUN apk add go sqlite

COPY --from=init /message-server /message-server

# Mkdir for data, don't persist
RUN mkdir -p /data/ && touch /data/sqlite.sq3

WORKDIR /message-server

RUN --mount=type=cache,target="/root/.cache/go-build" go build -race ./...

CMD ["go", "test", "./..."]

###### Build ######
FROM alpine:latest AS build

RUN apk add go && \
		apk add sqlite

COPY --from=init /message-server /message-server

WORKDIR /message-server

RUN --mount=type=cache,target="/root/.cache/go-build" go build -race -o /bin/message-server ./*.go 

###### Run ######
FROM alpine:latest AS api

# Persist data
VOLUME /data/
RUN mkdir -p /data/ && touch /data/sqlite.sq3

COPY --from=build /bin/message-server /bin/message-server

CMD ["/bin/message-server"]

