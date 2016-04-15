package plugin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"

	"github.infra.hana.ondemand.com/I061150/aker/socket"
)

type Plugin struct {
	http.Handler
	socketPath string
}

func (p *Plugin) SocketPath() string {
	if p == nil {
		return ""
	}
	return p.socketPath
}

func Open(name string, config []byte, next *Plugin) (*Plugin, error) {
	socketPath := socket.GetUniqueSocketPath("aker-plugin")

	setup, err := json.Marshal(&pluginSetup{
		SocketPath:        socketPath,
		ForwardSocketPath: next.SocketPath(),
		Configuration:     config,
	})
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(name)
	cmd.Stdin = bytes.NewReader(setup)
	cmd.Stdout = newPluginLogWriter(name, os.Stdout)
	cmd.Stderr = newPluginLogWriter(name, os.Stderr)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Plugin{
		socketPath: socketPath,
		Handler:    socket.Proxy(socketPath),
	}, nil
}
