FROM golang:1.11.0

RUN mkdir -p /go/src/github.com/defraglabs/uptime

COPY . /go/src/github.com/defraglabs/uptime

WORKDIR /go/src/github.com/defraglabs/uptime
RUN go install
ENTRYPOINT /go/bin/uptime