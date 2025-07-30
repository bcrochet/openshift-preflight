package container

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/image"
)

func getTrademarkLabels(bad bool) map[string]string {
	labels := map[string]string{
		"name":       "name",
		"vendor":     "vendor",
		"maintainer": "maintainer",
	}

	if bad {
		labels["maintainer"] = "Red Hat"
	}

	return labels
}

func createTestImageRefWithTrademarkLabels(bad bool) image.ImageReference {
	config := image.ParsedConfig{
		Config: struct {
			Labels map[string]string `json:"Labels"`
			Cmd    []string          `json:"Cmd"`
			User   string            `json:"User"`
		}{
			Labels: getTrademarkLabels(bad),
		},
	}

	configBytes, _ := json.Marshal(config)

	return image.ImageReference{
		ConfigBytes: configBytes,
	}
}

var _ = Describe("HasNoProhibitedLabelsCheck", func() {
	var (
		hasProhibitedLabelsCheck HasNoProhibitedLabelsCheck
		imageRef                 image.ImageReference
	)

	BeforeEach(func() {
		imageRef = createTestImageRefWithTrademarkLabels(false)
	})

	Describe("Checking for prohibited labels", func() {
		Context("When it has no prohibited labels", func() {
			It("should pass Validate", func() {
				ok, err := hasProhibitedLabelsCheck.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
		Context("When it has prohibited labels", func() {
			BeforeEach(func() {
				imageRef = createTestImageRefWithTrademarkLabels(true)
			})
			It("should not pass Validate", func() {
				ok, err := hasProhibitedLabelsCheck.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})

	AssertMetaData(&hasProhibitedLabelsCheck)
})
