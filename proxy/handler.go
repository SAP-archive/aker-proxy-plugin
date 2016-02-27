package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.wdf.sap.corp/I061150/aker/logging"

	"github.com/cloudfoundry-incubator/candiedyaml"
)

type handlerConfig struct {
	URL       string `yaml:"url"`
	ProxyPath string `yaml:"proxy_path"`
}

func NewHandler(config []byte) (http.Handler, error) {
	logging.Infof("Configuration: %s", string(config))
	cfg := handlerConfig{}
	if err := candiedyaml.Unmarshal(config, &cfg); err != nil {
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
			originalPath := removeProxyPath(req.URL.Path, cfg.ProxyPath)
			req.URL.Path = joinPaths(targetURL.Path, originalPath)
		},
	}, nil
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
