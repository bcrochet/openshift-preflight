package container

import "github.com/redhat-openshift-ecosystem/openshift-preflight/internal/image"

// getContainerLabels is a helper function to obtain the labels from an images configfile
func getContainerLabels(imgRef image.ImageReference) (map[string]string, error) {
	config, err := imgRef.GetConfig()
	if err != nil {
		return nil, err
	}
	return config.Config.Labels, nil
}
