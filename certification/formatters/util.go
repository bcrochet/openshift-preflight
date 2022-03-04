package formatters

import (
	"github.com/redhat-openshift-ecosystem/openshift-preflight/certification/runtime"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/version"
)

// getResponse will extract the runtime's results and format it to fit the
// UserResponse definition in a way that can then be formatted.
func getResponse(r runtime.Results) UserResponse {
	passedChecks := make([]checkExecutionInfo, 0, len(r.Passed))
	failedChecks := make([]checkExecutionInfo, 0, len(r.Failed))
	erroredChecks := make([]checkExecutionInfo, 0, len(r.Errors))

	if len(r.Passed) > 0 {
		for _, check := range r.Passed {
			passedChecks = append(passedChecks, checkExecutionInfo{
				Name:        check.Name(),
				ElapsedTime: float64(check.ElapsedTime.Milliseconds()),
				Description: check.Metadata().Description,
			})
		}
	}

	if len(r.Failed) > 0 {
		for _, check := range r.Failed {
			failedChecks = append(failedChecks, checkExecutionInfo{
				Name:             check.Name(),
				ElapsedTime:      float64(check.ElapsedTime.Milliseconds()),
				Description:      check.Metadata().Description,
				Help:             check.Help().Message,
				Suggestion:       check.Help().Suggestion,
				KnowledgeBaseURL: check.Metadata().KnowledgeBaseURL,
				CheckURL:         check.Metadata().CheckURL,
			})
		}
	}

	if len(r.Errors) > 0 {
		for _, check := range r.Errors {
			erroredChecks = append(erroredChecks, checkExecutionInfo{
				Name:        check.Name(),
				ElapsedTime: float64(check.ElapsedTime.Milliseconds()),
				Description: check.Metadata().Description,
				Help:        check.Help().Message,
			})
		}
	}

	response := UserResponse{
		Image:             r.TestedImage,
		Passed:            r.PassedOverall,
		LibraryInfo:       version.Version,
		CertificationHash: r.CertificationHash,
		// TestedOn:          r.TestedOn,
		Results: resultsText{
			Passed: passedChecks,
			Failed: failedChecks,
			Errors: erroredChecks,
		},
	}

	return response
}

// UserResponse is the standard user-facing response.
type UserResponse struct {
	Image             string                 `json:"image" xml:"image" graphql:"image"`
	Passed            bool                   `json:"passed" xml:"passed" graphql:"passed"`
	CertificationHash string                 `json:"certification_hash,omitempty" xml:"certification_hash,omitempty" graphql:"certification_hash"`
	LibraryInfo       version.VersionContext `json:"test_library" xml:"test_library" graphql:"test_library"`
	Results           resultsText            `json:"results" xml:"results" graphql:"results"`
	// TestedOn          runtime.OpenshiftClusterVersion `json:"tested_on" xml:"tested_on"`
}

// resultsText represents the results of check execution against the asset.
type resultsText struct {
	Passed []checkExecutionInfo `json:"passed" xml:"passed" graphql:"passed"`
	Failed []checkExecutionInfo `json:"failed" xml:"failed" graphql:"failed"`
	Errors []checkExecutionInfo `json:"errors" xml:"errors" graphql:"errors"`
}

// checkExecutionInfo contains all possible output fields that a user might see in their result.
// Empty fields will be omitted.
type checkExecutionInfo struct {
	Name             string  `json:"name,omitempty" xml:"name,omitempty" graphql:"name"`
	ElapsedTime      float64 `json:"elapsed_time" xml:"elapsed_time" graphql:"elapsed_time"`
	Description      string  `json:"description,omitempty" xml:"description,omitempty" graphql:"description"`
	Help             string  `json:"help,omitempty" xml:"help,omitempty" graphql:"help"`
	Suggestion       string  `json:"suggestion,omitempty" xml:"suggestion,omitempty" graphql:"suggestion"`
	KnowledgeBaseURL string  `json:"knowledgebase_url,omitempty" xml:"knowledgebase_url,omitempty" graphql:"knowledgebase_url"`
	CheckURL         string  `json:"check_url,omitempty" xml:"check_url,omitempty" graphql:"check_url"`
}
