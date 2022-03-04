package pyxis

import (
	"github.com/hasura/go-graphql-client"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/certification/formatters"
)

type CertImage struct {
	ID                     string       `json:"_id,omitempty" graphql:"_id"`
	Certified              bool         `json:"certified" default:"false" graphql:"certified"`
	Deleted                bool         `json:"deleted" default:"false" graphql:"deleted"`
	DockerImageDigest      string       `json:"docker_image_digest,omitempty" graphql:"docker_image_digest"`
	DockerImageID          string       `json:"docker_image_id,omitempty" graphql:"docker_image_id"`
	ImageID                string       `json:"image_id,omitempty" graphql:"image_id"`
	ISVPID                 string       `json:"isv_pid,omitempty" graphql:"isv_pid"` // required
	ParsedData             *ParsedData  `json:"parsed_data,omitempty" graphql:"parsed_data"`
	Architecture           string       `json:"architecture" default:"amd64" graphql:"architecture"`
	RawConfig              string       `json:"raw_config,omitempty" graphql:"raw_config"`
	Repositories           []Repository `json:"repositories,omitempty" graphql:"repositories"`
	SumLayerSizeBytes      int64        `json:"sum_layer_size_bytes,omitempty" graphql:"sum_layer_size_bytes"`
	UncompressedTopLayerId string       `json:"uncompressed_top_layer_id,omitempty" graphql:"uncompressed_top_layer_id"` // TODO: figure out how to populate this, it is not required
}

type ParsedData struct {
	Architecture           string  `json:"architecture,omitempty" graphql:"architecture"`
	Command                string  `json:"command,omitempty" graphql:"command"`
	Comment                string  `json:"comment,omitempty" graphql:"comment"`
	Container              string  `json:"container,omitempty" graphql:"container"`
	Created                string  `json:"created,omitempty" graphql:"created"`
	DockerVersion          string  `json:"docker_version,omitempty" graphql:"docker_version"`
	ImageID                string  `json:"image_id,omitempty" graphql:"image_id"`
	Labels                 []Label `json:"labels,omitempty" graphql:"labels"` // required
	OS                     string  `json:"os,omitempty" graphql:"os"`
	Ports                  string  `json:"ports,omitempty" graphql:"ports"`
	Size                   int64   `json:"size,omitempty" graphql:"size"`
	UncompressedLayerSizes []Layer `json:"uncompressed_layer_sizes,omitempty" graphql:"uncompressed_layer_sizes"`
}

type Repository struct {
	Published  bool   `json:"published" default:"false" graphql:"published"`
	PushDate   string `json:"push_date,omitempty" graphql:"push_date"` // time.Now
	Registry   string `json:"registry,omitempty" graphql:"registry"`
	Repository string `json:"repository,omitempty" graphql:"repository"`
	Tags       []Tag  `json:"tags,omitempty" graphql:"tags"`
}

type Label struct {
	Name  string `json:"name" graphql:"name"`
	Value string `json:"value" graphql:"value"`
}

type Tag struct {
	AddedDate string `json:"added_date,omitempty" graphql:"added_date"` // time.Now
	Name      string `json:"name,omitempty" graphql:"name"`
}

type RPMManifest struct {
	ID      string `json:"_id,omitempty" graphql:"_id"`
	ImageID string `json:"image_id,omitempty" graphql:"image_id"`
	RPMS    []RPM  `json:"rpms,omitempty" graphql:"rpms"`
}

type RPM struct {
	Architecture string `json:"architecture,omitempty" graphql:"architecture"`
	Gpg          string `json:"gpg,omitempty" graphql:"gpg"`
	Name         string `json:"name,omitempty" graphql:"name"`
	Nvra         string `json:"nvra,omitempty" graphql:"nvra"`
	Release      string `json:"release,omitempty" graphql:"release"`
	SrpmName     string `json:"srpm_name,omitempty" graphql:"srpm_name"`
	SrpmNevra    string `json:"srpm_nevra,omitempty" graphql:"srpm_nevra"`
	Summary      string `json:"summary,omitempty" graphql:"summary"`
	Version      string `json:"version,omitempty" graphql:"version"`
}

type CertificationProjectInput struct {
	ID                  graphql.String `json:"_id,omitempty" graphql:"_id"`
	CertificationStatus string         `json:"certification_status" default:"In Progress" graphql:"certification_status"`
}

type (
	StringFilter                           map[string]string
	ContainerImageFilter                   map[string]StringFilter
	CertificationImageInput                CertImage
	ContainerImageRPMManifestInput         RPMManifest
	CertProjectTestResultInput             TestResults
	ContainerImageInput                    CertImage
	ContainerImageResponse                 CertImage
	ContainerImageRPMManifestResponse      RPMManifest
	ContainerImagePaginatedResponse        struct{}
	CertificationProjectResponse           CertProject
	CertProjectTestResultResponse          TestResults
	CertProjectArtifactResponse            struct{}
	CertProjectTestResultFilter            struct{}
	CertProjectTestResultPaginatedResponse struct{}
	CertProjectArtifactFilter              struct{}
	CertProjectArtifactPaginatedResponse   struct{}
	SortBy                                 struct{}
)

type CertProject struct {
	ID                  string    `json:"_id,omitempty" graphql:"_id"`
	CertificationStatus string    `json:"certification_status" default:"In Progress" graphql:"certification_status"`
	Container           Container `json:"container" graphql:"container"`
	Name                string    `json:"name" graphql:"name"`                      // required
	ProjectStatus       string    `json:"project_status" graphql:"project_status"`  // required
	Type                string    `json:"type" default:"Containers" graphql:"type"` // required
}

type Container struct {
	DockerConfigJSON string `json:"docker_config_json" graphql:"docker_config_json"`
	Type             string `json:"type" default:"Containers" graphql:"type"` // conditionally required
	ISVPID           string `json:"isv_pid,omitempty" graphql:"isv_pid"`      // required
	OsContentType    string `json:"os_content_type,omitempty" graphql:"os_content_type"`
}

type GraphqlError struct {
	Status int    `graphql:"status"`
	Detail string `graphql:"detail"`
}

type Layer struct {
	LayerId string `json:"layer_id" graphql:"layer_id"`
	Size    int64  `json:"size_bytes" graphql:"size_bytes"`
}

type TestResults struct {
	ID          string `json:"_id,omitempty" graphql:"_id"`
	CertProject string `json:"cert_project,omitempty" graphql:"cert_project"`
	OrgID       int    `json:"org_id,omitempty" graphql:"org_id"`
	Version     string `json:"version,omitempty" graphql:"version"`
	ImageID     string `json:"image_id,omitempty" graphql:"image_id"`
	formatters.UserResponse
}
