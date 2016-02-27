package proxy_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.wdf.sap.corp/I061150/aker-proxy/proxy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Handler", func() {
	var configuration []byte
	var handler http.Handler
	var handlerErr error

	buildConfiguration := func(targetURL, proxyPath string) []byte {
		buffer := bytes.Buffer{}
		buffer.WriteString("---\n")
		buffer.WriteString("url: ")
		buffer.WriteString(targetURL)
		buffer.WriteString("\n")
		buffer.WriteString("proxy_path: ")
		buffer.WriteString(proxyPath)
		buffer.WriteString("\n")
		return buffer.Bytes()
	}

	JustBeforeEach(func() {
		handler, handlerErr = NewHandler(configuration)
	})

	Context("when configuration is invalid YAML", func() {
		BeforeEach(func() {
			configuration = []byte("&asdINVALID_YAML:^HERE")
		})

		It("handler creation should fail", func() {
			Ω(handlerErr).Should(HaveOccurred())
		})
	})

	Context("when configuration contains invalid URL", func() {
		BeforeEach(func() {
			configuration = buildConfiguration("::INVALID::URL", "/")
		})

		It("handler creation should fail", func() {
			Ω(handlerErr).Should(HaveOccurred())
		})
	})

	Context("when configuration is valid", func() {
		const serverResponsePayload = "SomeContent"

		var fakeServer *ghttp.Server
		var response *httptest.ResponseRecorder
		var request *http.Request

		itHandlerCreationShouldNotFail := func() {
			It("handler creation should not fail", func() {
				Ω(handler).ShouldNot(BeNil())
				Ω(handlerErr).ShouldNot(HaveOccurred())
			})
		}

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
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			fakeServer.Close()
		})

		JustBeforeEach(func() {
			handler.ServeHTTP(response, request)
		})

		Context("when both target path and proxy path are empty", func() {
			BeforeEach(func() {
				configuration = buildConfiguration(fakeServer.URL(), "")
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/first/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itHandlerCreationShouldNotFail()

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when the target path is non-empty but proxy path is empty", func() {
			BeforeEach(func() {
				configuration = buildConfiguration(fakeServer.URL()+"/zero", "")
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/zero/first/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itHandlerCreationShouldNotFail()

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when the proxy path is non-empty", func() {
			BeforeEach(func() {
				configuration = buildConfiguration(fakeServer.URL(), "/first")
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itHandlerCreationShouldNotFail()

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when proxy path does not match request path", func() {
			BeforeEach(func() {
				configuration = buildConfiguration(fakeServer.URL(), "/notFirst")
				fakeServer.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/first/second", "q=1"),
					ghttp.RespondWith(http.StatusOK, serverResponsePayload),
				))
			})

			itHandlerCreationShouldNotFail()

			itShouldCallTheServer()

			itShouldReturnProperResponse()
		})

		Context("when both target path and proxy path are non-empty", func() {
			Context("when target path does not end with slash", func() {
				BeforeEach(func() {
					configuration = buildConfiguration(fakeServer.URL()+"/zero", "/first")
					fakeServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/zero/second", "q=1"),
						ghttp.RespondWith(http.StatusOK, serverResponsePayload),
					))
				})

				itHandlerCreationShouldNotFail()

				itShouldCallTheServer()

				itShouldReturnProperResponse()
			})

			Context("when target path ends with slash", func() {
				BeforeEach(func() {
					configuration = buildConfiguration(fakeServer.URL()+"/zero/", "/first")
					fakeServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/zero/second", "q=1"),
						ghttp.RespondWith(http.StatusOK, serverResponsePayload),
					))
				})

				itHandlerCreationShouldNotFail()

				itShouldCallTheServer()

				itShouldReturnProperResponse()
			})

			Context("when proxy path is only slash", func() {
				BeforeEach(func() {
					configuration = buildConfiguration(fakeServer.URL()+"/zero", "/")
					fakeServer.AppendHandlers(ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/zero/first/second", "q=1"),
						ghttp.RespondWith(http.StatusOK, serverResponsePayload),
					))
				})

				itHandlerCreationShouldNotFail()

				itShouldCallTheServer()

				itShouldReturnProperResponse()
			})
		})
	})
})
