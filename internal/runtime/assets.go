package runtime

import (
	"context"
	"fmt"
	goruntime "runtime"
	"strings"

	"github.com/go-logr/logr"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/log"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/containers/image/v5/manifest"
)

// images maps the images use by preflight with their purpose.
//
// these should have accessor functions made available if they are
// to be used outside of this package.
var images = map[string]string{
	// operator policy, operator-sdk scorecard
	"scorecard": "quay.io/operator-framework/scorecard-test:v1.40.0",
}

// imageList takes the images mapping and represents them using just
// the image URIs.
func imageList(ctx context.Context) []string {
	logger := logr.FromContextOrDiscard(ctx)
	
	// Create system context for containers/image
	sys := &types.SystemContext{
		ArchitectureChoice: goruntime.GOARCH,
		OSChoice:           "linux",
	}

	imageList := make([]string, 0, len(images))

	for _, image := range images {
		base := strings.Split(image, ":")[0]
		
		// Parse image reference
		imageRef := image
		if !strings.HasPrefix(imageRef, "docker://") {
			imageRef = "docker://" + imageRef
		}
		ref, err := alltransports.ParseImageName(imageRef)
		if err != nil {
			logger.Error(fmt.Errorf("could not parse image reference: %w", err), "image reference error")
			continue
		}
		
		// Get image source and manifest to calculate digest
		src, err := ref.NewImageSource(ctx, sys)
		if err != nil {
			logger.Error(fmt.Errorf("could not create image source: %w", err), "image source error")
			continue
		}
		
		manifestBytes, _, err := src.GetManifest(ctx, nil)
		if err != nil {
			src.Close()
			logger.Error(fmt.Errorf("could not get manifest: %w", err), "manifest error")
			continue
		}
		src.Close()
		
		// Calculate manifest digest
		digest, err := manifest.Digest(manifestBytes)
		if err != nil {
			logger.Error(fmt.Errorf("could not calculate manifest digest: %w", err), "digest error")
			continue
		}
		
		imageList = append(imageList, fmt.Sprintf("%s@%s", base, digest.String()))
	}

	return imageList
}

// Assets returns a full collection of assets used in Preflight.
func Assets(ctx context.Context) AssetData {
	return AssetData{
		Images: imageList(ctx),
	}
}

// ScorecardImage returns the container image used for OperatorSDK
// Scorecard based checks. If userProvidedScorecardImage is set, it is
// returned, otherwise, the default is returned.
func ScorecardImage(ctx context.Context, userProvidedScorecardImage string) string {
	logger := logr.FromContextOrDiscard(ctx)
	if userProvidedScorecardImage != "" {
		logger.V(log.DBG).Info("user provided scorecard test image", "image", userProvidedScorecardImage)
		return userProvidedScorecardImage
	}
	return images["scorecard"]
}

// Assets is the publicly accessible representation of Preflight's
// used assets. This struct will be serialized to JSON and presented
// to the end-user when requested.
type AssetData struct {
	Images []string `json:"images"`
}
