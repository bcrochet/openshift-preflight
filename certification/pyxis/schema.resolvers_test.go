package pyxis

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CreateImage(ctx context.Context, input *ContainerImageInput) (*ContainerImageResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateImage(ctx context.Context, id *string, input *ContainerImageInput) (*ContainerImageResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ReplaceImage(ctx context.Context, id *string, input *ContainerImageInput) (*ContainerImageResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateImageRpmManifest(ctx context.Context, id *string, input *ContainerImageRPMManifestInput) (*ContainerImageRPMManifestResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ReplaceImageRpmManifest(ctx context.Context, id *string, input *ContainerImageRPMManifestInput) (*ContainerImageRPMManifestResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateImageRpmManifest(ctx context.Context, id *string, input *ContainerImageRPMManifestInput) (*ContainerImageRPMManifestResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateCertificationProject(ctx context.Context, input *CertificationProjectInput) (*CertificationProjectResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCertificationProject(ctx context.Context, id *string, input *CertificationProjectInput) (*CertificationProjectResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ReplaceCertificationProject(ctx context.Context, id *string, input *CertificationProjectInput) (*CertificationProjectResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateCertificationProjectTestResult(ctx context.Context, id *string, input *CertProjectTestResultInput) (*CertProjectTestResultResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCertificationProjectTestResult(ctx context.Context, id *string, input *CertProjectTestResultInput) (*CertProjectTestResultResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateCertificationProjectArtifact(ctx context.Context, id *string, input *CertProjectTestResultInput) (*CertProjectArtifactResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetImage(ctx context.Context, id *string) (*ContainerImageResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindImages(ctx context.Context, sortBy []*SortBy, page *int, pageSize *int, filter *ContainerImageFilter) (*ContainerImagePaginatedResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCertProjectTestResult(ctx context.Context, id *string) (*CertProjectTestResultResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCertProjectTestResults(ctx context.Context, id *string, sortBy []*SortBy, page *int, pageSize *int, filter *CertProjectTestResultFilter) (*CertProjectTestResultPaginatedResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetTestResultsByImageID(ctx context.Context, id *string, sortBy []*SortBy, page *int, pageSize *int, filter *CertProjectTestResultFilter) (*CertProjectTestResultPaginatedResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCertProjectArtifact(ctx context.Context, id *string) (*CertProjectArtifactResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCertProjectArtifacts(ctx context.Context, id *string, sortBy []*SortBy, page *int, pageSize *int, filter *CertProjectArtifactFilter) (*CertProjectArtifactPaginatedResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCertificationArtifactsByImageID(ctx context.Context, id *string, sortBy []*SortBy, page *int, pageSize *int, filter *CertProjectArtifactFilter) (*CertProjectArtifactPaginatedResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetRpmManifest(ctx context.Context, id *string) (*ContainerImageRPMManifestResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetImageRpmManifest(ctx context.Context, id *string) (*ContainerImageRPMManifestResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCertificationProject(ctx context.Context, id *string) (*CertificationProjectResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCertificationProjectByPid(ctx context.Context, pid *string) (*CertificationProjectResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindCertificationProjectsImages(ctx context.Context, id *string, sortBy []*SortBy, page *int, pageSize *int, filter *ContainerImageFilter) (*ContainerImagePaginatedResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)

type Resolver struct{}
