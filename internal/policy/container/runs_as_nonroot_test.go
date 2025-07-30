package container

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/image"
)

func createTestImageRefWithUser(user string) image.ImageReference {
	config := image.ParsedConfig{
		Config: struct {
			Labels map[string]string `json:"Labels"`
			Cmd    []string          `json:"Cmd"`
			User   string            `json:"User"`
		}{
			User: user,
		},
	}

	configBytes, _ := json.Marshal(config)

	return image.ImageReference{
		ConfigBytes: configBytes,
	}
}

var _ = Describe("RunAsNonRoot", func() {
	var (
		runAsNonRoot RunAsNonRootCheck
		imageRef     image.ImageReference
	)

	BeforeEach(func() {
		imageRef = createTestImageRefWithUser("1000")
	})

	Describe("Checking manifest user is not root", func() {
		Context("When manifest user is not root", func() {
			It("should pass Validate", func() {
				ok, err := runAsNonRoot.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
	})
	Describe("Checking manifest user is root", func() {
		Context("When manifest user is empty", func() {
			BeforeEach(func() {
				imageRef = createTestImageRefWithUser("")
			})
			It("should not pass Validate", func() {
				ok, err := runAsNonRoot.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
		Context("When manifest user is string root", func() {
			BeforeEach(func() {
				imageRef = createTestImageRefWithUser("root")
			})
			It("should not pass Validate", func() {
				ok, err := runAsNonRoot.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
		Context("When manifest user is UID 0", func() {
			BeforeEach(func() {
				imageRef = createTestImageRefWithUser("0")
			})
			It("should not pass Validate", func() {
				ok, err := runAsNonRoot.Validate(context.TODO(), imageRef)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})

	AssertMetaData(&runAsNonRoot)
})
