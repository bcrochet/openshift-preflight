// Package types contains all types relevant to this PoC.
//
// This is organized into a single place just for PoC purposes.
// These are copied from preflight because preflight contains these in
// an internal package.
package preflight

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/version"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// CheckEngine defines the functionality necessary to run all checks for a policy,
// and return the results of that check execution.
type CheckEngine interface {
	// ExecuteChecks should execute all checks in a policy and internally
	// store the results. Errors returned by ExecuteChecks should reflect
	// errors in pre-validation tasks, and not errors in individual check
	// execution itself.
	ExecuteChecks(context.Context) error
	// Results returns the outcome of executing all checks.
	Results(context.Context) Results
}

// Check as an interface containing all methods necessary
// to use and identify a given check.
type Check interface {
	// Validate will test the provided image and determine whether the
	// image complies with the check's requirements.
	Validate(ctx context.Context, imageReference ImageReference) (result bool, err error)
	// Name returns the name of the check.
	Name() string
	// Metadata returns the check's metadata.
	Metadata() Metadata
	// Help return the check's help information
	Help() HelpText
}

// ImageReference holds all things image-related
type ImageReference struct {
	ImageURI        string
	ImageFSPath     string
	ImageInfo       v1.Image
	ImageRepository string
	ImageRegistry   string
	ImageTagOrSha   string
}

type Result struct {
	Check
	ElapsedTime time.Duration
}

type Results struct {
	TestedImage       string
	PassedOverall     bool
	TestedOn          OpenshiftClusterVersion
	CertificationHash string
	Passed            []Result
	Failed            []Result
	Errors            []Result
}

// Metadata contains useful information regarding the check.
type Metadata struct {
	// Description contains a brief text detailing the overall goal of the check.
	Description string `json:"description" xml:"description"`
	// Level describes the certification level associated with the given check.
	//
	// TODO: define this more explicitly when requirements surrounding this metadata
	// text.
	Level string `json:"level" xml:"level"`
	// KnowledgeBaseURL is a URL detailing how to resolve a check failure.
	KnowledgeBaseURL string `json:"knowledge_base_url,omitempty" xml:"knowledgeBaseURL"`
	// CheckURL is a URL pointing to the official policy documentation from Red Hat, containing
	// information on exactly what is being tested and why.
	CheckURL string `json:"check_url,omitempty" xml:"checkURL"`
}

// HelpText is the help message associated with any given check
type HelpText struct {
	// Message is text provided to the user indicating where they should look
	// to find out why they failed or encountered an error in validation.
	Message string `json:"message" xml:"message"`
	// Suggestion is text provided to the user indicating what might need to
	// change in order to pass a check.
	Suggestion string `json:"suggestion" xml:"suggestion"`
}

type OpenshiftClusterVersion struct {
	Name    string
	Version string
}

func UnknownOpenshiftClusterVersion() OpenshiftClusterVersion {
	return OpenshiftClusterVersion{
		Name:    "unknown",
		Version: "unknown",
	}
}

// ResponseFormatter describes the expected methods a formatter
// must implement.
type ResponseFormatter interface {
	// PrettyName is the name used to represent this formatter.
	PrettyName() string
	// FileExtension represents the file extension one might use when creating
	// a file with the contents of this formatter.
	FileExtension() string
	// Format takes Results, formats it as needed, and returns the formatted
	// results ready to write as a byte slice.
	Format(context.Context, Results) (response []byte, formattingError error)
}

// ResultWriter defines methods associated with writing check results.
type ResultWriter interface {
	OpenFile(name string) (io.WriteCloser, error)
	io.WriteCloser
}

// ResultSubmitter defines methods associated with submitting results to Red HAt.
type ResultSubmitter interface {
	Submit(context.Context) error
}

// UserResponse is the standard user-facing response.
type UserResponse struct {
	Image             string                 `json:"image" xml:"image"`
	Passed            bool                   `json:"passed" xml:"passed"`
	CertificationHash string                 `json:"certification_hash,omitempty" xml:"certification_hash,omitempty"`
	LibraryInfo       version.VersionContext `json:"test_library" xml:"test_library"`
	Results           resultsText            `json:"results" xml:"results"`
}

// resultsText represents the results of check execution against the asset.
type resultsText struct {
	Passed []checkExecutionInfo `json:"passed" xml:"passed"`
	Failed []checkExecutionInfo `json:"failed" xml:"failed"`
	Errors []checkExecutionInfo `json:"errors" xml:"errors"`
}

// checkExecutionInfo contains all possible output fields that a user might see in their result.
// Empty fields will be omitted.
type checkExecutionInfo struct {
	Name             string  `json:"name,omitempty" xml:"name,omitempty"`
	ElapsedTime      float64 `json:"elapsed_time" xml:"elapsed_time"`
	Description      string  `json:"description,omitempty" xml:"description,omitempty"`
	Help             string  `json:"help,omitempty" xml:"help,omitempty"`
	Suggestion       string  `json:"suggestion,omitempty" xml:"suggestion,omitempty"`
	KnowledgeBaseURL string  `json:"knowledgebase_url,omitempty" xml:"knowledgebase_url,omitempty"`
	CheckURL         string  `json:"check_url,omitempty" xml:"check_url,omitempty"`
}

// Borrowed directly from Preflight's internal/runtime
// ResultWriterFile implements a ResultWriter for use at preflight runtime.
type ResultWriterFile struct {
	file *os.File
}

// OpenFile will open the expected file for writing.
func (f *ResultWriterFile) OpenFile(name string) (io.WriteCloser, error) {
	file, err := os.OpenFile(
		name,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0o600)
	if err != nil {
		return nil, err
	}

	f.file = file // so we can close it later.
	return f, nil
}

func (f *ResultWriterFile) Close() error {
	return f.file.Close()
}

func (f *ResultWriterFile) Write(p []byte) (int, error) {
	return f.file.Write(p)
}
