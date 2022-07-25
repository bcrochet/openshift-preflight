package pyxis

import (
	"strings"

	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
)

var _ = Describe("SHA calc", func() {
	When("given a certain imput", func() {
		It("will return the correct output", func() {
			expected := "fe813dd3baf1ccc8f1f82fb66b009fad7e82be4d842029dbb2e8d70d657f3a49"
			actual, err := BinarySHA(strings.NewReader("this is my test string that needs to be long enough"))
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})
	})
	When("a bad reader is passed", func() {
		It("will fail", func() {
			_, err := BinarySHA(errReader(0))
			Expect(err).To(HaveOccurred())
		})
	})
})
