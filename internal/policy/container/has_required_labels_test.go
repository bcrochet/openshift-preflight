package container

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/image"
)

func getLabels(bad bool) map[string]string {
	labels := map[string]string{
		"name":        "name",
		"maintainer":  "maintainer",
		"vendor":      "vendor",
		"version":     "version",
		"release":     "release",
		"summary":     "summary",
		"description": "description",
	}

	if bad {
		delete(labels, "description")
	}

	return labels
}

func createTestImageRef(bad bool) image.ImageReference {
	config := image.ParsedConfig{
		Config: struct {
			Labels map[string]string `json:"Labels"`
			Cmd    []string          `json:"Cmd"`
			User   string            `json:"User"`
		}{
			Labels: getLabels(bad),
		},
	}

	configBytes, _ := json.Marshal(config)

	return image.ImageReference{
		ConfigBytes: configBytes,
	}
}

var _ = Describe("HasRequiredLabels", func() {
	var (
		hasRequiredLabelsCheck HasRequiredLabelsCheck
		imageRef               image.ImageReference
	)

	BeforeEach(func() {
		imageRef = createTestImageRef(false)
	})

	Describe("Checking for required labels", func() {
		Context("When it has required labels", func() {
			It("should pass Validate", func() {
				ok, err := hasRequiredLabelsCheck.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
		Context("When it does not have required labels", func() {
			BeforeEach(func() {
				imageRef = createTestImageRef(true)
			})
			It("should not succeed the check", func() {
				ok, err := hasRequiredLabelsCheck.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})

	AssertMetaData(&hasRequiredLabelsCheck)
})
