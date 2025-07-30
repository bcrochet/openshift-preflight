package image

import (
	"encoding/json"
	"fmt"

	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/types"
	"github.com/opencontainers/go-digest"
)

// ImageReference holds all things image-related
type ImageReference struct {
	ImageURI           string
	ImageFSPath        string
	ImageSource        types.ImageSource
	Manifest           manifest.Manifest
	ConfigBytes        []byte
	ImageRepository    string
	ImageRegistry      string
	ImageTagOrSha      string
	ManifestListDigest string
	// For testing purposes - mock layer count
	MockLayerCount int
}

// ParsedConfig represents the parsed image configuration
type ParsedConfig struct {
	Architecture  string `json:"architecture"`
	OS            string `json:"os"`
	Created       string `json:"created"`
	DockerVersion string `json:"docker_version"`
	Config        struct {
		Labels map[string]string `json:"Labels"`
		Cmd    []string          `json:"Cmd"`
		User   string            `json:"User"`
	} `json:"config"`
	RootFS struct {
		DiffIDs []string `json:"diff_ids"`
	} `json:"rootfs"`
}

// GetConfig parses and returns the image configuration
func (ir *ImageReference) GetConfig() (*ParsedConfig, error) {
	var config ParsedConfig
	if err := json.Unmarshal(ir.ConfigBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse image config: %w", err)
	}
	return &config, nil
}

// GetLayerDiffIDs returns the layer diff IDs from the image config
func (ir *ImageReference) GetLayerDiffIDs() ([]digest.Digest, error) {
	config, err := ir.GetConfig()
	if err != nil {
		return nil, err
	}

	diffIDs := make([]digest.Digest, 0, len(config.RootFS.DiffIDs))
	for _, diffID := range config.RootFS.DiffIDs {
		digest, err := digest.Parse(diffID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse diff ID %s: %w", diffID, err)
		}
		diffIDs = append(diffIDs, digest)
	}
	return diffIDs, nil
}

// GetImageDigest returns the digest of the image config
func (ir *ImageReference) GetImageDigest() (digest.Digest, error) {
	configInfo := ir.Manifest.ConfigInfo()
	if configInfo.Digest == "" {
		return "", fmt.Errorf("image has no config info")
	}
	return configInfo.Digest, nil
}

// GetLayerCount returns the number of layers in the image
func (ir *ImageReference) GetLayerCount() int {
	// Use mock layer count for testing if set
	if ir.MockLayerCount > 0 {
		return ir.MockLayerCount
	}
	// Otherwise use actual manifest
	if ir.Manifest != nil {
		return len(ir.Manifest.LayerInfos())
	}
	return 0
}
