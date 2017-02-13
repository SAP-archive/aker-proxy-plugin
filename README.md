[![CLA assistant](https://cla-assistant.io/readme/badge/SAP/aker-proxy-plugin)](https://cla-assistant.io/SAP/aker-proxy-plugin)

Aker Proxy Plugin
=================

The Aker proxy plugin forwards requests to a remote server. This plugin consumes the request and will be the last in the chain of plugins for a given path.

# License
This project is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.

# User Guide

Following is an example of a valid configuration for this plugin.

```yaml
url: http://example.org
proxy_path: "/"
preserve_internal_headers: true
flush_interval: 300ms
```

The `url` property specifies the remote endpoint to which requests will be forwarded.

The `proxy_path` property can be used to remove part of the original request path.

The `preserve_internal_headers` property specifies whether `x-aker-*` headers will be forwarded to the remote target. If the remote resources is hosted by an untrusted provider, then it makes sense to keep this value `false`.

The `flush_interval` property can be used to specify the flush interval to the [ReverseProxy](https://github.com/golang/go/blob/master/src/net/http/httputil/reverseproxy.go). It shares the same format as the duration string in [ParseDuration](https://golang.org/pkg/time/#ParseDuration). Defaults to zero which means there is no periodic flushing.

For example, with the following configuration in Aker,

```yaml
endpoints:
  - path: /two/segments
    plugins:
      - ...
        configuration:
          url: http://example.org/target
          proxy_path: "/two"
          preserve_internal_headers: true
```

if you were to access Aker on `/two/segments/suffix`, the requests would be forwarded to `http://example.org/target/segments/suffix`.

## Tests

`aker-proxy-plugin` project contains unit tests, in order to execute them run the following command in project root directory.

```bash
ginkgo -r
```
