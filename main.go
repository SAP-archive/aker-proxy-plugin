package main

import (
	"github.wdf.sap.corp/I061150/aker-proxy/core"
	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func main() {
	plugin.ListenAndServe("aker-proxy", func() (api.Plugin, error) {
		return core.NewPlugin(), nil
	})
}
