FROM golang:1.17-alpine

ADD . /go/src/streamserver
WORKDIR /go/src/streamserver

RUN go get streamserver
RUN go install

ENTRYPOINT ["/go/bin/streamserver -port=8080 -addr=0.0.0.0"]