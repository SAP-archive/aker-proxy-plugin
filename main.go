package main

import (
	"github.infra.hana.ondemand.com/cloudfoundry/aker-proxy-plugin/proxy"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/plugin"
	"github.infra.hana.ondemand.com/cloudfoundry/gologger"
)

func main() {
	if err := plugin.ListenAndServeHTTP(proxy.NewHandlerFromRawConfig); err != nil {
		gologger.Fatalf("Error creating plugin: %v", err)
	}
}
