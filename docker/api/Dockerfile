FROM golang:1.11.0

RUN mkdir -p /go/src/github.com/dineshs91/uptime

COPY . /go/src/github.com/dineshs91/uptime

WORKDIR /go/src/github.com/dineshs91/uptime
RUN go build
ENTRYPOINT /go/src/github.com/dineshs91/uptime/docker/api/entrypoint.sh