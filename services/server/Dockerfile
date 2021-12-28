FROM golang:1.17-alpine

ADD . /go/src/goproject
WORKDIR /go/src/goproject

RUN go get goproject
RUN go install

ENTRYPOINT ["/go/bin/goproject -port=8080 -addr=0.0.0.0"]