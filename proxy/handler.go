package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type handlerConfig struct {
	URL       string `json:"url"`
	ProxyPath string `json:"proxy_path"`
}

func NewHandler(config []byte) (http.Handler, error) {
	fmt.Printf("Configuration: %s\n", string(config))
	cfg := handlerConfig{}
	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	targetURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			req.URL.Path = removeLeadingPath(req.URL.Path, cfg.ProxyPath)
			req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)
		},
	}, nil
}

func removeLeadingPath(path, leading string) string {
	leading = strings.TrimRight(leading, "/")
	return path[len(leading):]
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
