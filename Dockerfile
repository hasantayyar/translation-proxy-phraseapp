FROM golang:1.10-alpine

RUN apk update && apk upgrade
RUN apk add git

RUN mkdir -p /go/src/github.com/thesoenke/translation-proxy
COPY . /go/src/github.com/thesoenke/translation-proxy

WORKDIR /go/src/github.com/thesoenke/translation-proxy
RUN go get
RUN go build -o /usr/local/bin/translation-proxy

CMD /usr/local/bin/translation-proxy
