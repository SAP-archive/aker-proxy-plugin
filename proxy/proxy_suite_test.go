package proxy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.infra.hana.ondemand.com/I061150/aker/logging"

	"testing"
)

func TestProxy(t *testing.T) {
	logging.DefaultLogger = new(logging.MutedLogger)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Proxy Suite")
}
