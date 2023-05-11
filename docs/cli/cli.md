# CLI Client

[source](https://github.com/gokch/ipfs-mount/blob/main/cmd/cli/cli.go)



```
Usage:
  client [flags]

Flags:
  -c, --cids stringArray    download cid
  -e, --expire int          expire seconds (default 600)
  -h, --help                help for client
  -p, --paths stringArray   download path per cid
      --peers stringArray   connect peer id
  -r, --rootpath string     root path (default "./")
  -t, --timeout int         timeout seconds, 0 is no timeout
  -w, --worker int          worker size (default 1)
```