package pyxis

import (
	"context"
	"errors"
	"net/http"

	graphqlserver "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Pyxis Submit", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).ToNot(HaveOccurred())
				Expect(certProject).ToNot(BeNil())
				Expect(certImage).ToNot(BeNil())
				Expect(testResults).ToNot(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit updateProject 401 Unauthorized", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCertProjectUnauthorizedClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
				Expect(certImage).To(BeNil())
				Expect(testResults).To(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit with createImage 409 Conflict", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCreateImageConflictClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).ToNot(HaveOccurred())
				Expect(certProject).ToNot(BeNil())
				Expect(certImage).ToNot(BeNil())
				Expect(testResults).ToNot(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit with createImage 401 Unauthorized", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCreateImageUnauthorizedClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
				Expect(certImage).To(BeNil())
				Expect(testResults).To(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit with createImage 409 Conflict and getImage 401 Unauthorized ", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCreateImageConflictAndUnauthorizedClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
				Expect(certImage).To(BeNil())
				Expect(testResults).To(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit with createRPMManifest 409 Conflict", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCreateRPMManifestConflictClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).ToNot(HaveOccurred())
				Expect(certProject).ToNot(BeNil())
				Expect(certImage).ToNot(BeNil())
				Expect(testResults).ToNot(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit with createRPMManifest 401 Unauthorized", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCreateRPMManifestUnauthorizedClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
				Expect(certImage).To(BeNil())
				Expect(testResults).To(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit with createRPMManifest 409 Conflict and getRPMManifest 401 Unauthorized", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCreateRPMManifestConflictAndUnauthorizedClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
				Expect(certImage).To(BeNil())
				Expect(testResults).To(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis Submit with createTestResults 401 Unauthorized", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCreateTestResultsUnauthorizedClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisEngine.SubmitResults(&CertProject{CertificationStatus: "Started"}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
				Expect(certImage).To(BeNil())
				Expect(testResults).To(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis GetProejct", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, err := pyxisEngine.GetProject(context.Background())
				Expect(err).ToNot(HaveOccurred())
				Expect(certProject).ToNot(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis GetProject 401 Unauthorized", func() {
	var pyxisEngine *pyxisEngine

	BeforeEach(func() {
		pyxisEngine = NewPyxisEngine("my-spiffy-api-token", "my-awseome-project-id", fakeHttpCertProjectUnauthorizedClient{})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, err := pyxisEngine.GetProject(context.Background())
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
			})
		})
	})
})

var _ = Describe("Pyxis GraphQL Submit", func() {
	var pyxisGraphqlEngine *pyxisGraphqlEngine
	schema, err := graphqlserver.ParseSchema(Schema, &Resolver{})
	if err != nil {
		log.Error(err)
	}
	Expect(err).ToNot(HaveOccurred())
	mux := http.NewServeMux()
	mux.Handle("/query", &relay.Handler{Schema: schema})

	BeforeEach(func() {
		pyxisGraphqlEngine = NewPyxisGraphqlEngine("/query", "spiffyapitoken", "cool-project", &http.Client{Transport: localRoundTripper{handler: mux}})
	})
	Context("when a project is submitted", func() {
		Context("and it is not already In Progress", func() {
			It("should switch to In Progress", func() {
				certProject, certImage, testResults, err := pyxisGraphqlEngine.SubmitResults(context.Background(), &CertProject{}, &CertImage{}, &RPMManifest{}, &TestResults{})
				Expect(err).To(MatchError(errors.New("error calling remote API")))
				Expect(certProject).To(BeNil())
				Expect(certImage).To(BeNil())
				Expect(testResults).To(BeNil())
			})
		})
	})
})
