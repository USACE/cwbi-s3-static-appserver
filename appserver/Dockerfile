FROM golang:1.16-alpine AS builder

RUN addgroup -g 50000 -S data \
  && adduser -u 50000 -S data -G data

RUN apk update && apk add --no-cache git

# Copy In Source Code
WORKDIR /go/src/app
COPY main.go go.mod go.sum ./
# Copy In entrypoint.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod 700 /entrypoint.sh \
  && chown data:data /entrypoint.sh

# Build
RUN go get -d -v \
  && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -ldflags="-w -s" -o /go/bin/appserver

RUN mkdir /data && chown data:data /data && chmod 755 /data

USER data

VOLUME [ "/data" ]

ENTRYPOINT ["/entrypoint.sh"]
