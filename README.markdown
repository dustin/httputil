# HTTP Tracking

HTTP client usage tracking is useful when you have HTTP client
activity and want to know more about them.  Example:

```
In-flight HTTP requests:
  servicing GET "http://host1:8484/.cbfs/blob/[oid]" for 56.705274ms
  servicing GET "http://host1:8484/.cbfs/blob/[oid]" for 73.17592ms
  servicing GET "http://host2:8484/.cbfs/blob/[oid]" for 1.528276ms
  servicing GET "http://host3:8484/.cbfs/blob/[oid]" for 20.101718ms
```

You can also enable stack tracking on each one of those so you can
know where they come from.  Example:

```
In-flight HTTP requests:
  servicing GET "http://host1:8484/.cbfs/blob/[oid]" for 23.659896ms
    - net/http.send() - $GOROOT/src/pkg/net/http/client.go:139
    - net/http.(*Client).send() - $GOROOT/src/pkg/net/http/client.go:94
    - net/http.(*Client).doFollowingRedirects() - $GOROOT/src/pkg/net/http/client.go:251
    - net/http.(*Client).Get() - $GOROOT/src/pkg/net/http/client.go:243
    - net/http.Get() - $GOROOT/src/pkg/net/http/client.go:224
    - github.com/couchbaselabs/cbfs/client.fetchWorker.Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/fetch.go:64
    - github.com/couchbaselabs/cbfs/client.(*fetchWorker).Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/client.go:1
    - github.com/dustin/go-saturate.(*Saturator).destWorker() - $GOPATH/src/github.com/dustin/go-saturate/saturate.go:74
  servicing GET "http://host2:8484/.cbfs/blob/[oid]" for 414.055us
    - net/http.send() - $GOROOT/src/pkg/net/http/client.go:139
    - net/http.(*Client).send() - $GOROOT/src/pkg/net/http/client.go:94
    - net/http.(*Client).doFollowingRedirects() - $GOROOT/src/pkg/net/http/client.go:251
    - net/http.(*Client).Get() - $GOROOT/src/pkg/net/http/client.go:243
    - net/http.Get() - $GOROOT/src/pkg/net/http/client.go:224
    - github.com/couchbaselabs/cbfs/client.fetchWorker.Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/fetch.go:64
    - github.com/couchbaselabs/cbfs/client.(*fetchWorker).Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/client.go:1
    - github.com/dustin/go-saturate.(*Saturator).destWorker() - $GOPATH/src/github.com/dustin/go-saturate/saturate.go:74
  servicing GET "http://host1:8484/.cbfs/blob/[oid]" for 617.711758ms
    - net/http.send() - $GOROOT/src/pkg/net/http/client.go:139
    - net/http.(*Client).send() - $GOROOT/src/pkg/net/http/client.go:94
    - net/http.(*Client).doFollowingRedirects() - $GOROOT/src/pkg/net/http/client.go:251
    - net/http.(*Client).Get() - $GOROOT/src/pkg/net/http/client.go:243
    - net/http.Get() - $GOROOT/src/pkg/net/http/client.go:224
    - github.com/couchbaselabs/cbfs/client.fetchWorker.Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/fetch.go:64
    - github.com/couchbaselabs/cbfs/client.(*fetchWorker).Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/client.go:1
    - github.com/dustin/go-saturate.(*Saturator).destWorker() - $GOPATH/src/github.com/dustin/go-saturate/saturate.go:74
  servicing GET "http://host3:8484/.cbfs/blob/[oid]" for 19.561697ms
    - net/http.send() - $GOROOT/src/pkg/net/http/client.go:139
    - net/http.(*Client).send() - $GOROOT/src/pkg/net/http/client.go:94
    - net/http.(*Client).doFollowingRedirects() - $GOROOT/src/pkg/net/http/client.go:251
    - net/http.(*Client).Get() - $GOROOT/src/pkg/net/http/client.go:243
    - net/http.Get() - $GOROOT/src/pkg/net/http/client.go:224
    - github.com/couchbaselabs/cbfs/client.fetchWorker.Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/fetch.go:64
    - github.com/couchbaselabs/cbfs/client.(*fetchWorker).Work() - $GOPATH/src/github.com/couchbaselabs/cbfs/client/client.go:1
    - github.com/dustin/go-saturate.(*Saturator).destWorker() - $GOPATH/src/github.com/dustin/go-saturate/saturate.go:74
```
