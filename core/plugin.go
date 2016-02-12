package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func NewPlugin(log plugin.Logger) api.Plugin {
	return &plug{
		log: log,
	}
}

type plug struct {
	log       plugin.Logger
	target    *url.URL
	proxyPath string
	handler   http.Handler
}

type pluginConfig struct {
	URL       string `json:"url"`
	ProxyPath string `json:"proxy_path"`
}

func (p *plug) Config(data []byte) error {
	p.log.Info(fmt.Sprintf("Configuration: %s", string(data)))
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
	p.log.Info("Process...")

	if p.handler == nil {
		p.log.Error("Plugin has not been configured!")
		os.Exit(1)
	}

	req := &http.Request{
		URL:    context.Request.URL(),
		Method: context.Request.Method(),
		Body:   context.Request,
	}
	resp := &responseWrapper{
		delegate: context.Response,
	}
	p.handler.ServeHTTP(resp, req)
	resp.Close()

	p.log.Info("Done!")
	return true
}

type responseWrapper struct {
	delegate      api.Response
	headers       http.Header
	headerWritten bool
}

func (w *responseWrapper) Header() http.Header {
	return w.headers
}

func (w *responseWrapper) WriteHeader(status int) {
	for name, value := range w.headers {
		w.delegate.SetHeader(name, value)
	}
	w.delegate.WriteStatus(status)
	w.headerWritten = true
}

func (w *responseWrapper) Write(data []byte) (int, error) {
	return w.delegate.Write(data)
}

func (w *responseWrapper) Close() {
	if !w.headerWritten {
		w.WriteHeader(http.StatusOK)
	}
}
