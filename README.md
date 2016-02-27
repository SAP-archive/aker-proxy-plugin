Aker Proxy Plugin
=================

The Aker proxy plugin forwards requests to a remote server. This plugin consumes the request and will be the last in the chain of plugins for a given path.

Following is an example of a valid configuration for this plugin.

```yaml
url: http://example.org
proxy_path: "/"
```

The `url` property specifies the remote endpoint to which requests will be forwarded.

The `proxy_path` property can be used to remove part of the original request path.

For example, with the following configuration in Aker,

```yaml
endpoints:
  - path: /two/segments
    plugins:
      - ...
        configuration:
          url: http://example.org/target
          proxy_path: "/two"
```

if you were to access Aker on `/two/segments/suffix`, the requests would be forwarded to `http://example.org/target/segments/suffix`.
