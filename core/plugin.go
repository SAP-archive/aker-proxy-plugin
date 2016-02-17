package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.wdf.sap.corp/I061150/aker/adapter"
	"github.wdf.sap.corp/I061150/aker/api"
)

func NewPlugin() api.Plugin {
	return &plug{}
}

type plug struct {
	target    *url.URL
	proxyPath string
	handler   http.Handler
}

type pluginConfig struct {
	URL       string `json:"url"`
	ProxyPath string `json:"proxy_path"`
}

func (p *plug) Config(data []byte) error {
	fmt.Printf("Configuration: %s\n", string(data))
	cfg := pluginConfig{}
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	targetURL, err := url.Parse(cfg.URL)
	if err != nil {
		return err
	}
	director := func(req *http.Request) {
		req.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = removeLeadingPath(req.URL.Path, cfg.ProxyPath)
		req.URL.Path = singleJoiningSlash(targetURL.Path, req.URL.Path)
	}
	p.handler = &httputil.ReverseProxy{
		Director: director,
	}
	return nil
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

func (p *plug) Process(context api.Context) bool {
	if p.handler == nil {
		log.Fatalln("Plugin has not been configured!")
	}

	req := adapter.NewRequest(context.Request)
	resp := adapter.NewResponseWriterAdapter(context.Response)
	p.handler.ServeHTTP(resp, req)
	resp.Flush()

	return true
}
