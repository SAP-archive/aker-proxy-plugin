package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.wdf.sap.corp/I061150/aker-proxy/core"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func main() {
	log, err := plugin.NewLogger("aker-proxy")
	if err != nil {
		panic(err)
	}

	p := core.NewPlugin(log)
	plugin.ListenAndServe(p)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	for range ch {
		os.Exit(0)
	}
}
