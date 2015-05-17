FROM quay.io/coreos/etcd

RUN go get -u -v -t github.com/tleyden/etcd-discovery/...

