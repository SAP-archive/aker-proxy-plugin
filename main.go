package main

import (
	"github.wdf.sap.corp/I061150/aker-proxy/proxy"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func main() {
	plugin.ListenAndServe(proxy.NewHandler)
}
