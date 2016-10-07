package main

import (
	"github.com/SAP/aker-proxy-plugin/proxy"
	"github.com/SAP/aker/plugin"
	"github.com/SAP/gologger"
)

func main() {
	if err := plugin.ListenAndServeHTTP(proxy.NewHandlerFromRawConfig); err != nil {
		gologger.Fatalf("Error creating plugin: %v", err)
	}
}
