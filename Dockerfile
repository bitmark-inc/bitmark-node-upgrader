FROM golang:1.12.1-alpine

RUN apk update && apk upgrade && apk add git curl netcat-openbsd wget net-tools vim bash

ENV GOPATH /go
ENV PATH="/go/bin:${PATH}"
ENV GO111MODULE on

RUN cd /go/src && \
git clone https://github.com/bitmark-inc/bitmark-node-upgrader && \
cd /go/src/bitmark-node-upgrader && go mod download && \
go install && cd /go/bin

ADD dockerAssets/start.sh /
RUN cd / && chmod +x start.sh
CMD /start.sh

