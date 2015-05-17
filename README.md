Etcd-with-discovery is a wrapper around [etcd](https://github.com/coreos/etcd), with some discovery-related additions that calls etcd with the right *advertising* related URLs.

## The problem this is trying to solve

As described in [this Google Groups thread](https://groups.google.com/d/msg/google-containers/rFIFD6Y0_Ew/GeDa8ZuPWd8J), when running etcd in the context of Kubernetes, there's no straighforward way to discover what IP should be used for the `-advertise-client-urls` etcd parameter.

Without a correct `-advertise-client-urls` etcd parameter, code that uses the etcd client API (such as etcdctl) will not work.  (curl, on the other hand, will still work)

## How to use

Instead of invoking etcd via:

```
$ etcd -listen-client-urls http://0.0.0.0:2379 -advertise-client-urls http://10.1.50.34:2379
```

You would instead run:

```
$ etcdisco -listen-client-urls http://0.0.0.0:2379 -advertise-client-urls http://{{ .LOCAL_IP }}:2379
```

This would have the effect of:

* etcdisco would attempt to discover the ip it should be advertising
* it will look over the command line args and replace `{{ .LOCAL_IP }}` with the discovered ip address
* it will invoke etcd with the modified command line args



