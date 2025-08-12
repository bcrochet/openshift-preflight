package container

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/image"
)

func createTestImageRefWithLayerCount(layerCount int) image.ImageReference {
	config := image.ParsedConfig{
		Config: struct {
			Labels map[string]string `json:"Labels"`
			Cmd    []string          `json:"Cmd"`
			User   string            `json:"User"`
		}{},
	}

	configBytes, _ := json.Marshal(config)

	return image.ImageReference{
		ConfigBytes:    configBytes,
		MockLayerCount: layerCount,
	}
}

var _ = Describe("LessThanMaxLayers", func() {
	var (
		maxLayersCheck MaxLayersCheck
		imgRef         image.ImageReference
	)

	BeforeEach(func() {
		imgRef = createTestImageRefWithLayerCount(5) // Default to 5 layers (less than max)
	})

	Describe("Checking for less than max layers", func() {
		Context("When it has fewer layers than max", func() {
			It("should pass Validate", func() {
				ok, err := maxLayersCheck.Validate(context.TODO(), imgRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
		Context("When has more layers than max", func() {
			BeforeEach(func() {
				imgRef = createTestImageRefWithLayerCount(50) // More than max layers
			})
			It("should not succeed the check", func() {
				ok, err := maxLayersCheck.Validate(context.TODO(), imgRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})

	AssertMetaData(&maxLayersCheck)
})
