package pyxis

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/go-digest"
)

var _ = Describe("Pyxis CheckRedHatLayers", func() {
	ctx := context.Background()
	var pyxisClient *pyxisClient
	mux := http.NewServeMux()
	mux.HandleFunc("/query/", pyxisGraphqlLayerHandler(ctx))

	Context("when some layers are provided", func() {
		BeforeEach(func() {
			pyxisClient = NewPyxisClient("my.pyxis.host/query/", "my-spiffy-api-token", "my-awesome-project-id", &http.Client{Transport: localRoundTripper{handler: mux}})
		})
		Context("and a layer is a known good layer", func() {
			It("should be a good layer", func() {
				// Create a test digest
				testDigest := digest.FromString("test-layer")
				certImages, err := pyxisClient.CertifiedImagesContainingLayers(ctx, []digest.Digest{testDigest})
				Expect(err).ToNot(HaveOccurred())
				Expect(certImages).ToNot(BeNil())
				Expect(certImages).ToNot(BeZero())
			})
		})
	})
})
