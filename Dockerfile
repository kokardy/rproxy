from golang:1.15-buster

ENV GOPATH=/go
RUN mkdir -p /go

ENV SCHEME=http RHOST="www.google.com" ADDR=":80" ORI="" DEST=""

RUN go get github.com/kokardy/rproxy

WORKDIR /go/src/github.com/kokardy/rproxy
RUN go build .

CMD ["./server.sh"]

