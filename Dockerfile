FROM golang:wheezy

RUN go get -u -v -t github.com/tleyden/etcd-discovery/...

