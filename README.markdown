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
know where they come from.
