package container

import (
	"context"
	"fmt"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/check"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/image"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/log"

	"github.com/go-logr/logr"
)

const (
	acceptableLayerMax = 40
)

var _ check.Check = &MaxLayersCheck{}

// UnderLayerMaxCheck ensures that the image has less layers in its assembly than a predefined maximum.
type MaxLayersCheck struct{}

func (p *MaxLayersCheck) Validate(ctx context.Context, imgRef image.ImageReference) (bool, error) {
	layerCount := p.getDataToValidate(imgRef)
	return p.validate(ctx, layerCount)
}

func (p *MaxLayersCheck) getDataToValidate(imgRef image.ImageReference) int {
	return imgRef.GetLayerCount()
}

func (p *MaxLayersCheck) validate(ctx context.Context, layerCount int) (bool, error) {
	logr.FromContextOrDiscard(ctx).V(log.DBG).Info("number of layers detected in image", "layerCount", layerCount)
	return layerCount <= acceptableLayerMax, nil
}

func (p *MaxLayersCheck) Name() string {
	return "LayerCountAcceptable"
}

func (p *MaxLayersCheck) Metadata() check.Metadata {
	return check.Metadata{
		Description:      fmt.Sprintf("Checking if container has less than %d layers.  Too many layers within the container images can degrade container performance.", acceptableLayerMax),
		Level:            "better",
		KnowledgeBaseURL: certDocumentationURL,
		CheckURL:         certDocumentationURL,
	}
}

func (p *MaxLayersCheck) Help() check.HelpText {
	return check.HelpText{
		Message:    "Check LayerCountAcceptable encountered an error. Please review the preflight.log file for more information.",
		Suggestion: "Optimize your Dockerfile to consolidate and minimize the number of layers. Each RUN command will produce a new layer. Try combining RUN commands using && where possible.",
	}
}
