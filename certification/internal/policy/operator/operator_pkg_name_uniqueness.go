package operator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/certification"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/certification/internal/bundle"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// apiRespondData is the response received from the defined API
type apiResponseData struct {
	Data     []packageData `json:"data"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Total    int           `json:"total"`
}

// packageData represents a single package entry in the API response.
type packageData struct {
	Id             string `json:"_id"`
	Association    string `json:"association"`
	CreateDate     string `json:"creation_date"`
	LastUpdateTime string `json:"last_update_date"`
	PackageName    string `json:"package_name"`
	Source         string `json:"source"`
}

// OperatorPkgNameIsUniqueMountedCheck finds the package name as defined in the operator bundle's annotations
// and checks it against Red Hat APIs to confirm that the package name is unique at the time of the
// check.
type OperatorPkgNameIsUniqueCheck struct{}

func (p *OperatorPkgNameIsUniqueCheck) Validate(ctx context.Context, bundleRef certification.ImageReference) (bool, error) {
	// retrieve the operator metadata from bundle image
	annotationsFileName := filepath.Join(bundleRef.ImageFSPath, "metadata", "annotations.yaml")
	annotationsFile, err := os.Open(annotationsFileName)
	if err != nil {
		log.Error(fmt.Errorf("%w: could not open annotations.yaml", err))
		return false, err
	}
	annotations, err := bundle.GetAnnotations(ctx, annotationsFile)
	if err != nil {
		log.Error("unable to get annotations.yaml from the bundle")
		return false, err
	}

	packageName, err := p.getPackageName(ctx, annotations)
	if err != nil {
		log.Error("unable to extract package name from ClusterServiceVersion", err)
		return false, err
	}

	log.Debugf("operator package name is %s", packageName)

	req, err := p.buildRequest(ctx, apiEndpoint, packageName)
	if err != nil {
		log.Error("unable to build API request structure", err)
		return false, err
	}

	resp, err := p.queryAPI(ctx, http.DefaultClient, req)
	if err != nil {
		log.Error("unable to query package name validation API for uniqueness check", err)
		return false, err
	}

	data, err := p.parseAPIResponse(ctx, resp)
	if err != nil {
		log.Error("unable to parse response provided by package name validation API", err)
		return false, err
	}

	return p.validate(ctx, data)
}

// getPackageName accepts the annotations map and searches for the specified annotation corresponding
// with the complete bundle name for an operator, which is then returned.
func (p *OperatorPkgNameIsUniqueCheck) getPackageName(ctx context.Context, annotations map[string]string) (string, error) {
	log.Tracef("searching for package key (%s) in bundle", packageKey)
	log.Trace("bundle data: ", annotations)
	pkg, found := annotations[packageKey]
	if !found {
		return "", fmt.Errorf("did not find package name at the key %s in the annotations.yaml", packageKey)
	}

	return pkg, nil
}

// buildRequest builds the http.Request using the input parameters and returns a client.
func (p *OperatorPkgNameIsUniqueCheck) buildRequest(ctx context.Context, apiURL, packageName string) (*http.Request, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// this endpoint supports a query string, and we use that to determine if a
	// package with the name already exists.
	queryString := req.URL.Query()
	queryString.Add("filter", fmt.Sprintf("package_name==%s", packageName))
	req.URL.RawQuery = queryString.Encode()

	return req, nil
}

// queryAPI uses the provided client to query the remote API, and returns the response if it
// response is successful, or an error if the response was unexpected in any way.
func (p *OperatorPkgNameIsUniqueCheck) queryAPI(ctx context.Context, client apiClient, request *http.Request) (*http.Response, error) {
	log.Trace("making API request to ", request.URL.String())
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	log.Trace("response code: ", resp.Status)

	// The Connect API returns a 200 regardless of whether the package was found or not. Until this
	// assumption changes, we assume any non-200 response is invalid, or due to a malformed request.
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received an unexpected status code for the request")
	}

	return resp, nil
}

// parseAPIResponse reads the response and checks the body for the expected contents, and then
// returns the body content as apiResponseData.
func (p *OperatorPkgNameIsUniqueCheck) parseAPIResponse(ctx context.Context, resp *http.Response) (*apiResponseData, error) {
	var data apiResponseData
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Trace("response body: ", string(body))

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// validate checks the apiResponseData and confirms that the package is unique by confirming that the
// API returned no packages using the same name.
func (p *OperatorPkgNameIsUniqueCheck) validate(ctx context.Context, resp *apiResponseData) (bool, error) {
	// success case - the API returned no entries
	if len(resp.Data) == 0 {
		return true, nil
	} else if viper.IsSet("certprojectid") {
		for i := range resp.Data {
			rawCertProjectID := string(resp.Data[i].Association)
			strippedCertProjectID := strings.ReplaceAll(strings.TrimPrefix(rawCertProjectID, "ospid"), "-", "")
			targetCertProjectID := viper.GetString("certprojectid")
			if strings.EqualFold(targetCertProjectID, rawCertProjectID) || strings.EqualFold(targetCertProjectID, strippedCertProjectID) {
				return true, nil
			}
		}
	}

	log.Error("a package already exists in the Red Hat ecosystem using the same name")
	// we don't expect multiple entries, but resp.Data is a list so we will iterate.
	for _, v := range resp.Data {
		log.Error("found the following entry: ", v)
	}

	return false, nil
}

func (p *OperatorPkgNameIsUniqueCheck) Name() string {
	return "OperatorPackageNameIsUnique"
}

func (p *OperatorPkgNameIsUniqueCheck) Metadata() certification.Metadata {
	return certification.Metadata{
		Description:      "Validating Bundle image package name uniqueness",
		Level:            "best",
		KnowledgeBaseURL: "https://sdk.operatorframework.io/docs/olm-integration/tutorial-bundle/",
		CheckURL:         "https://sdk.operatorframework.io/docs/olm-integration/tutorial-bundle/",
	}
}

func (p *OperatorPkgNameIsUniqueCheck) Help() certification.HelpText {
	return certification.HelpText{
		Message:    "Check encountered an error. It is possible that the bundle package name already exist in the RedHat Catalog registry.",
		Suggestion: "Bundle package name must be unique meaning that it doesn't already exist in the RedHat Catalog registry",
	}
}

// apiClient is a simple interface encompassing the only http.Client method we utilize for preflight checks. This exists to
// enable mock implementations for testing purposes.
type apiClient interface {
	Do(req *http.Request) (*http.Response, error)
}
