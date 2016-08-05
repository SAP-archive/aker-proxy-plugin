package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.infra.hana.ondemand.com/cloudfoundry/aker/plugin"
)

type handlerConfig struct {
	URL                     string `yaml:"url"`
	ProxyPath               string `yaml:"proxy_path"`
	PreserveInternalHeaders bool   `yaml:"preserve_internal_headers"`
}

func NewHandlerFromRawConfig(config []byte) (http.Handler, error) {
	cfg := handlerConfig{}
	if err := plugin.UnmarshalConfig(config, &cfg); err != nil {
		return nil, err
	}
	return NewHandlerFromConfig(cfg)
}

func NewHandlerFromConfig(cfg handlerConfig) (http.Handler, error) {
	targetURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}
	return NewHandler(targetURL, cfg.ProxyPath, cfg.PreserveInternalHeaders), nil
}

func NewHandler(targetURL *url.URL, proxyPath string, preserveHeaders bool) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = targetURL.Host
			req.URL.Scheme = targetURL.Scheme
			req.URL.Host = targetURL.Host
			originalPath := removeProxyPath(req.URL.Path, proxyPath)
			req.URL.Path = joinPaths(targetURL.Path, originalPath)
			if !preserveHeaders {
				removeInternalHeaders(req.Header)
			}
		},
	}
}

func removeInternalHeaders(headers http.Header) {
	for name := range headers {
		if strings.HasPrefix(strings.ToLower(name), "x-aker") {
			delete(headers, name)
		}
	}
}

func removeProxyPath(path, proxyPath string) string {
	return strings.TrimPrefix(path, proxyPath)
}

func joinPaths(first, second string) string {
	firstSlash := strings.HasSuffix(first, "/")
	secondSlash := strings.HasPrefix(second, "/")
	switch {
	case firstSlash && secondSlash:
		return first + second[1:]
	case !firstSlash && !secondSlash:
		return first + "/" + second
	default:
		return first + second
	}
}
