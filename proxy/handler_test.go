package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"

	. "github.com/SAP/aker-proxy-plugin/proxy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Handler", func() {
	Context("when creating instances", func() {
		var config []byte
		Context("from valid configuration", func() {
			var proxyHandler *httputil.ReverseProxy

			itShouldCreateValidProxyHandlerFromRawConfig := func() {
				handler, err := NewHandlerFromRawConfig(config)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(handler).ShouldNot(BeNil())
				proxyHandler = handler.(*httputil.ReverseProxy)
			}

			It("should be able to parse flush interval", func() {
				config = []byte("flush_interval: 300ms")
				itShouldCreateValidProxyHandlerFromRawConfig()
				Ω(proxyHandler.FlushInterval).Should(Equal(300 * time.Millisecond))
			})

			It("should default to 0 when flush interval not specified", func() {
				config = []byte("url: http://localhost:8080/")
				itShouldCreateValidProxyHandlerFromRawConfig()
				Ω(proxyHandler.FlushInterval).Should(Equal(time.Duration(0)))
			})
		})

		Context("from invalid configuration", func() {
			itShouldFailWhenConfigurationIsInvalid := func() {
				handler, err := NewHandlerFromRawConfig(config)
				Ω(err).Should(HaveOccurred())
				Ω(handler).Should(BeNil())
			}

			It("should fail", func() {
				config = []byte("invalid")
				itShouldFailWhenConfigurationIsInvalid()
			})

			It("should fail when it's invalid URL", func() {
				config = []byte("url: http://invalid URL")
				itShouldFailWhenConfigurationIsInvalid()
			})
		})
	})

	Context("when configuration is valid", func() {
		const headerKey = "X-Aker-Custom-Header"
		const headerValue = "SomeValue"
		const serverResponsePayload = "SomeContent"

		var fakeServer *ghttp.Server
		var response *httptest.ResponseRecorder
		var request *http.Request
		var targetURL string
		var proxyPath string
		var preserveHeaders bool
		var handler http.Handler
		var flushInterval time.Duration

		itShouldCallTheServer := func() {
			It("should call the server", func() {
				Ω(fakeServer.ReceivedRequests()).Should(HaveLen(1))
			})
		}

		itShouldReturnProperResponse := func() {
			It("should return proper response", func() {
				Ω(response.Code).Should(Equal(http.StatusOK))
				Ω(response.Body.String()).Should(Equal(serverResponsePayload))
			})
		}

		BeforeEach(func() {
			fakeServer = ghttp.NewServer()
			response = httptest.NewRecorder()
			var err error
			request, err = http.NewRequest("GET", "http://example.com/first/second?q=1", nil)
			request.Header.Add(headerKey, headerValue)
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			fakeServer.Close()
		})

		JustBeforeEach(func() {
			parsedTargetURL, err := url.Parse(targetURL)
			Ω(err).ShouldNot(HaveOccurred())
			handler = NewHandler(parsedTargetURL, proxyPath, preserveHeaders, flushInterval)
			Ω(handler).ShouldNot(BeNil())
			handler.ServeHTTP(response, request)
		})

		Context("when internal headers are not preserved", func() {
			BeforeEach(func() {
				preserveHeaders = false
				targetURL = fakeServer.URL()
				proxyPath = ""
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/first/second", "q=1"),
					ghttp.VerifyHeader(http.Header{
						headerKey: nil,
					}),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itShouldCallTheServer()
		})

		Context("when internal headers are preserved", func() {
			BeforeEach(func() {
				preserveHeaders = true
				targetURL = fakeServer.URL()
				proxyPath = ""
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/first/second", "q=1"),
					ghttp.VerifyHeader(http.Header{
						headerKey: []string{headerValue},
					}),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itShouldCallTheServer()
		})

		Context("when both target path and proxy path are empty", func() {
			BeforeEach(func() {
				targetURL = fakeServer.URL()
				proxyPath = ""
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/first/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when the target path is non-empty but proxy path is empty", func() {
			BeforeEach(func() {
				targetURL = fakeServer.URL() + "/zero"
				proxyPath = ""
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/zero/first/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when the proxy path is non-empty", func() {
			BeforeEach(func() {
				targetURL = fakeServer.URL()
				proxyPath = "/first"
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when proxy path does not match request path", func() {
			BeforeEach(func() {
				targetURL = fakeServer.URL()
				proxyPath = "/notFirst"
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/first/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when both target path and proxy path are non-empty", func() {
			Context("when target path does not end with slash", func() {
				BeforeEach(func() {
					targetURL = fakeServer.URL() + "/zero"
					proxyPath = "/first"
					fakeServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/zero/second", "q=1"),
						ghttp.RespondWith(http.StatusOK, serverResponsePayload),
					))
				})

				itShouldCallTheServer()

				itShouldReturnProperResponse()
			})

			Context("when target path ends with slash", func() {
				BeforeEach(func() {
					targetURL = fakeServer.URL() + "/zero/"
					proxyPath = "/first"
					fakeServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/zero/second", "q=1"),
						ghttp.RespondWith(http.StatusOK, serverResponsePayload),
					))
				})

				itShouldCallTheServer()

				itShouldReturnProperResponse()
			})

			Context("when proxy path is only slash", func() {
				BeforeEach(func() {
					targetURL = fakeServer.URL() + "/zero"
					proxyPath = "/"
					fakeServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/zero/first/second", "q=1"),
						ghttp.RespondWith(http.StatusOK, serverResponsePayload),
					))
				})

				itShouldCallTheServer()

				itShouldReturnProperResponse()
			})
		})
	})
})
