package main

import (
	"github.wdf.sap.corp/I061150/aker-proxy/core"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func main() {
	p := core.NewPlugin()
	plugin.ListenAndServe(p)
}
