FROM alpine:latest

RUN apk add go

COPY message-server/ message-server/

WORKDIR message-server/

CMD ["go", "run", "main.go"]

