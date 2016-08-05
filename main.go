package main

import (
	"github.infra.hana.ondemand.com/cloudfoundry/aker-proxy/proxy"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/logging"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/plugin"
)

func main() {
	if err := plugin.ListenAndServeHTTP(proxy.NewHandlerFromRawConfig); err != nil {
		logging.Fatalf("Error creating plugin: %v", err)
	}
}
