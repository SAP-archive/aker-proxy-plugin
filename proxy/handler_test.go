package proxy_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.infra.hana.ondemand.com/cloudfoundry/aker-proxy/proxy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Handler", func() {
	const headerKey = "X-Aker-Custom-Header"
	const headerValue = "SomeValue"

	var targetURL string
	var proxyPath string
	var preserveHeaders bool
	var handler http.Handler

	JustBeforeEach(func() {
		parsedTargetURL, err := url.Parse(targetURL)
		Ω(err).ShouldNot(HaveOccurred())
		handler = NewHandler(parsedTargetURL, proxyPath, preserveHeaders)
		Ω(handler).ShouldNot(BeNil())
	})

	Context("when configuration is valid", func() {
		const serverResponsePayload = "SomeContent"

		var fakeServer *ghttp.Server
		var response *httptest.ResponseRecorder
		var request *http.Request

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
