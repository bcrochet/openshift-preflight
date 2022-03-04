package pyxis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hasura/go-graphql-client"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/certification/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type pyxisGraphqlEngine struct {
	ProjectId string
	Client    *graphql.Client
}

type graphqlClient interface {
	Query(context.Context, interface{}, map[string]interface{}, ...graphql.Option)
	Mutate(context.Context, interface{}, map[string]interface{}, ...graphql.Option)
}

func GetPyxisGraphqlUrl() string {
	return fmt.Sprintf("https://%s/graphql/", viper.GetString("pyxis_host"))
}

func NewPyxisGraphqlEngine(graphqlUrl string, apiToken string, projectId string, httpClient *http.Client) *pyxisGraphqlEngine {
	client := graphql.NewClient(graphqlUrl, httpClient).WithRequestModifier(func(r *http.Request) {
		r.Header.Add("X-API-KEY", apiToken)
		// log.Debugf("%+v", r)
	})
	return &pyxisGraphqlEngine{
		ProjectId: projectId,
		Client:    client,
	}
}

func (p *pyxisGraphqlEngine) GetProject(ctx context.Context) (*CertProject, error) {
	var query struct {
		CertifiationProject struct {
			Data struct {
				CertProject
			} `graphql:"data"`
			GraphqlError `graphql:"error"`
		} `graphql:"get_certification_project(id: $id)"`
	}
	variables := map[string]interface{}{
		"id": graphql.String(p.ProjectId),
	}

	if err := p.Client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.CertifiationProject.Data.CertProject, nil
}

func (p *pyxisGraphqlEngine) updateProject(ctx context.Context, certProject *CertProject) (*CertProject, error) {
	var m struct {
		UpdateCertificationProject struct {
			Data struct {
				CertProject
				// ID                  string `graphql:"_id"`
				// CertificationStatus string `graphql:"certification_status"`
			} `graphql:"data"`
			GraphqlError `graphql:"error"`
		} `graphql:"update_certification_project(id: $id, input: $certificationProject)"`
	}
	variables := map[string]interface{}{
		"id": graphql.String(p.ProjectId),
		"certificationProject": CertificationProjectInput{
			ID:                  graphql.String(certProject.ID),
			CertificationStatus: certProject.CertificationStatus,
		},
	}

	err := p.Client.Mutate(ctx, &m, variables)
	if err != nil {
		return nil, err
	}

	// newProject, err := p.GetProject(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// return newProject, nil
	return &m.UpdateCertificationProject.Data.CertProject, nil
}

func (p *pyxisGraphqlEngine) createImage(ctx context.Context, certImage *CertImage) (*CertImage, *GraphqlError, error) {
	var m struct {
		CreateImage struct {
			Data struct {
				CertImage
			} `graphql:"data"`
			GraphqlError `graphql:"error"`
		} `graphql:"create_image(input: $containerImage)"`
	}
	variables := map[string]interface{}{
		"containerImage": CertificationImageInput(*certImage),
	}

	if err := p.Client.Mutate(ctx, &m, variables); err != nil {
		log.Errorf("error in createImage mutate: %+v", m)
		return nil, nil, err
	}

	if m.CreateImage.GraphqlError.Status < 200 || m.CreateImage.GraphqlError.Status > 299 {
		log.Debugf("%s", m.CreateImage.GraphqlError.Detail)
		return nil, &m.CreateImage.GraphqlError, fmt.Errorf("%s: %w", m.CreateImage.GraphqlError.Detail, errors.ErrNon200StatusCode)
	}

	return &m.CreateImage.Data.CertImage, nil, nil
}

func (p *pyxisGraphqlEngine) createRPMManifest(ctx context.Context, rpmManifest *RPMManifest) (*RPMManifest, *GraphqlError, error) {
	var m struct {
		CreateImageRPMManifest struct {
			Data struct {
				RPMManifest
			} `graphql:"data"`
			GraphqlError `graphql:"error"`
		} `graphql:"create_image_rpm_manifest(id: $id, input: $rpmManifest)"`
	}
	variables := map[string]interface{}{
		"id":          graphql.String(rpmManifest.ImageID),
		"rpmManifest": ContainerImageRPMManifestInput(*rpmManifest),
	}

	if err := p.Client.Mutate(ctx, &m, variables); err != nil {
		return nil, nil, err
	}

	if m.CreateImageRPMManifest.Status < 200 || m.CreateImageRPMManifest.Status > 299 {
		log.Debugf("%s", m.CreateImageRPMManifest.GraphqlError.Detail)
		return nil, &m.CreateImageRPMManifest.GraphqlError, fmt.Errorf("%s: %w", m.CreateImageRPMManifest.GraphqlError.Detail, errors.ErrNon200StatusCode)
	}

	return &m.CreateImageRPMManifest.Data.RPMManifest, &m.CreateImageRPMManifest.GraphqlError, nil
}

func (p *pyxisGraphqlEngine) getRPMManifest(ctx context.Context, imageId string) (*RPMManifest, error) {
	var m struct {
		GetImageRPMManifest struct {
			Data struct {
				RPMManifest
			} `graphql:"data"`
			GraphqlError `graphql:"error"`
		} `graphql:"get_image_rpm_manifest(id: $id)"`
	}
	variables := map[string]interface{}{
		"id": graphql.String(imageId),
	}

	if err := p.Client.Query(ctx, &m, variables); err != nil {
		return nil, err
	}

	return &m.GetImageRPMManifest.Data.RPMManifest, nil
}

func (p *pyxisGraphqlEngine) createTestResults(ctx context.Context, testResults *TestResults) (*TestResults, error) {
	var m struct {
		CreateCertifciationProjectTestResult struct {
			Data struct {
				TestResults
			} `graphql:"data"`
			GraphqlError `graphql:"error"`
		} `graphql:"create_certification_project_test_result(id: $id, input: $testResults)"`
	}
	variables := map[string]interface{}{
		"id":          graphql.String(p.ProjectId),
		"testResults": CertProjectTestResultInput(*testResults),
	}

	if err := p.Client.Mutate(ctx, &m, variables); err != nil {
		return nil, err
	}

	return &m.CreateCertifciationProjectTestResult.Data.TestResults, nil
}

func (p *pyxisGraphqlEngine) getImage(ctx context.Context, imageDigest string) (*CertImage, *GraphqlError, error) {
	var query struct {
		GetContainerImage struct {
			Data struct {
				ContainerImage []CertImage
			} `graphql:"data"`
			GraphqlError `graphql:"error"`
		} `graphql:"find_images(filter: $filter)"`
	}
	variables := map[string]interface{}{
		"filter": ContainerImageFilter{
			"docker_image_digest": {
				"eq": imageDigest,
			},
		},
	}

	if err := p.Client.Query(ctx, &query, variables); err != nil {
		return nil, &query.GetContainerImage.GraphqlError, err
	}

	return &query.GetContainerImage.Data.ContainerImage[0], nil, nil
}
