package main

import (
	"github.infra.hana.ondemand.com/I061150/aker-proxy/proxy"
	"github.infra.hana.ondemand.com/I061150/aker/plugin"
)

func main() {
	plugin.ListenAndServe(proxy.NewHandlerFromRawConfig)
}
