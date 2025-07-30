package container

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/check"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/image"
	"github.com/redhat-openshift-ecosystem/openshift-preflight/internal/log"

	"github.com/go-logr/logr"
	rpmdb "github.com/knqyf263/go-rpmdb/pkg"
	"github.com/spf13/afero"
)

var _ check.Check = &HasModifiedFilesCheck{}

// HasModifiedFilesCheck evaluates that no files from the base layer have been modified by
// subsequent layers by comparing the file list installed by Packages against the file list
// modified in subsequent layers.
type HasModifiedFilesCheck struct{}

type packageMeta struct {
	Name        string
	Version     string
	Release     string
	Arch        string
	Vendor      string
	InstallTime int
}

func (pm packageMeta) Compare(other packageMeta) int {
	return 0
}

type packageFilesRef struct {
	// LayerFiles contains a slice of files created/modified in layer
	LayerFiles map[string]fileInfo
	// LayerPackages is a map of the packages in a layer. The key is
	// the NVR of the package. The value is metadata about the package
	// that we use for processing
	LayerPackages map[string]packageMeta
	// LayerPackageFiles maps files to a package name-version-release
	LayerPackageFiles map[string]string
	HasRPMDB          bool
}

// Validate runs the check of whether any Red Hat files were modified
func (p *HasModifiedFilesCheck) Validate(ctx context.Context, imgRef image.ImageReference) (bool, error) {
	fs := afero.NewOsFs()
	layerIDs, packageFiles, err := p.gatherDataToValidate(ctx, imgRef, fs)
	if err != nil {
		return false, fmt.Errorf("could not generate modified files list: %v", err)
	}

	packageDist, err := p.parsePackageDist(ctx, imgRef.ImageFSPath, fs)
	if err != nil {
		return false, fmt.Errorf("could not generate modified files list: %v", err)
	}

	return p.validate(ctx, layerIDs, packageFiles, packageDist)
}

// parsePackageDist returns the platform's distribution value from the
// os-release contents in the extracted image.
func (p *HasModifiedFilesCheck) parsePackageDist(_ context.Context, extractedImageFSPath string, fs afero.Fs) (string, error) {
	osRelease, err := fs.Open(filepath.Join(extractedImageFSPath, "etc", "os-release"))
	if err != nil {
		return "", fmt.Errorf("could not open os-release: %v", err)
	}
	defer osRelease.Close()
	scanner := bufio.NewScanner(osRelease)
	packageDist := "unknown"
	r, err := regexp.Compile(`PLATFORM_ID="platform:([[:alnum:]]+)"`)
	if err != nil {
		return "", fmt.Errorf("error while compiling regexp: %w", err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		m := r.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		packageDist = m[1]
		break
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error while scanning for package dist: %v", err)
	}

	return packageDist, nil
}

// gatherDataToValidate returns a map from layer digest to a struct containing the list of files
// (packageFilesRef.LayerPackageFiles) installed via packages (packageFilesRef.LayerPackages)
// from the container image, and the list of files (packageFilesRef.LayerFiles) modified/added
// via layers in the image.
func (p *HasModifiedFilesCheck) gatherDataToValidate(ctx context.Context, imgRef image.ImageReference, fs afero.Fs) ([]string, map[string]packageFilesRef, error) {
	logger := logr.FromContextOrDiscard(ctx)

	layerDir, err := afero.TempDir(fs, "", "rpm-layers-")
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = fs.RemoveAll(layerDir)
	}()

	if imgRef.ImageSource == nil {
		return nil, nil, fmt.Errorf("image reference invalid")
	}

	// Get layer information from manifest
	layerInfos := imgRef.Manifest.LayerInfos()

	layerIDs := make([]string, 0, len(layerInfos))
	layerRefs := make(map[string]packageFilesRef, len(layerInfos))

	// Get diff IDs from the config for debugging
	diffIDs, err := imgRef.GetLayerDiffIDs()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get layer diff IDs: %w", err)
	}

	// Uncompress each layer and build maps containing the packages,
	// the package files, and the files modified by each layer
	// Also generate a list of the layerIDs so we can keep the
	// order within the maps.
	for idx, layerInfo := range layerInfos {
		layerIDHash := layerInfo.Digest

		// Capture the diff ID to aid in debugging. We don't technically care if
		// there's an error returned here because we don't use the layerDiffID
		// value for anything meaningful.
		layerDiffID := "unknown"
		if idx < len(diffIDs) {
			layerDiffID = diffIDs[idx].String()
		}

		rawLayerID := layerIDHash.String()
		// Map everything using a combination of the layer index and the layer
		// ID to avoid problems when images have multiple scattered layers with
		// the same ID.
		layerID := fmt.Sprintf("%02d-%s", idx, rawLayerID)
		logger.V(log.TRC).Info("generating unique layer ID", "uniqueLayerID", layerID, "layerID", rawLayerID, "layerDiffID", layerDiffID)

		layerDir := filepath.Join(layerDir, layerID)
		err = fs.Mkdir(layerDir, 0o755)
		if err != nil {
			return nil, nil, fmt.Errorf("could not create layer directory: %w", err)
		}

		layerIDs = append(layerIDs, layerID)

		// For now, skip layer processing until we implement the blob-based functions
		// TODO: Implement generateChangesForBlob and findRPMDBInBlob
		files := make(map[string]fileInfo)
		found := false
		pkgList := make([]*rpmdb.PackageInfo, 0)
		if !found {
			logger.V(log.TRC).Info("could not find rpm database in layer", "layer", layerID)
			if idx > 0 {
				// Just make this is the same as last layer, since the RPM db was not modified
				lastLayer := layerIDs[idx-1]
				layerRefs[layerID] = packageFilesRef{
					LayerFiles:        files,
					LayerPackages:     layerRefs[lastLayer].LayerPackages,
					LayerPackageFiles: layerRefs[lastLayer].LayerPackageFiles,
					HasRPMDB:          false,
				}
				continue
			}

			// If it's the first layer, just make the pkgList empty.
			pkgList = make([]*rpmdb.PackageInfo, 0)
		}

		pkgNameList := extractPackageNameVersionRelease(pkgList)

		packageFiles, err := installedFileMapWithExclusions(ctx, pkgList)
		if err != nil {
			return nil, nil, err
		}

		layerRefs[layerID] = packageFilesRef{
			LayerFiles:        files,
			LayerPackages:     pkgNameList,
			LayerPackageFiles: packageFiles,
			HasRPMDB:          true,
		}
	}

	return layerIDs, layerRefs, nil
}

// validate compares the list of LayerFiles and PackageFiles to see what PackageFiles
// have been modified within the additional layers. packageDist is the value we expect
// to find in the base package's Release field.
func (p *HasModifiedFilesCheck) validate(ctx context.Context, layerIDs []string, packageFiles map[string]packageFilesRef, packageDist string) (bool, error) {
	logger := logr.FromContextOrDiscard(ctx)

	disallowedModifications := false
	for idx, layerID := range layerIDs {
		logger := logger.WithValues("layer", layerID)
		ref := packageFiles[layerID]
		for modifiedFile, modifiedFileInfo := range ref.LayerFiles {
			logger := logger.WithValues("file", modifiedFile)

			// If it's a modifiedFile but this layer has an RPM db, that's allowed, but only if the
			// package itself is updated.
			if idx == 0 && ref.HasRPMDB {
				// Just skip these in the first layer
				continue
			}
			if _, found := ref.LayerPackageFiles[modifiedFile]; !found {
				// Far as we can tell, this isn't from an RPM
				continue
			}
			previousPackageVersion, prevFound := packageFiles[layerIDs[idx-1]].LayerPackageFiles[modifiedFile]
			if !prevFound && ref.HasRPMDB {
				// This is a net-new package file. Pass it.
				continue
			}
			currentPackageVersion := ref.LayerPackageFiles[modifiedFile]
			previousPackage := packageFiles[layerIDs[idx-1]].LayerPackages[previousPackageVersion]
			currentPackage := ref.LayerPackages[currentPackageVersion]

			if previousPackageVersion == currentPackageVersion {
				previousFileInfo := fileInfo{}
				// Since the modified file will not necessarily be present in the immediately previous layer, we need
				// to go backwards through the layers to look for the last time this file was in a layer, and get the
				// mode from there.
				for layerIdx := idx - 1; layerIdx > -1; layerIdx-- {
					if pfi, found := packageFiles[layerIDs[layerIdx]].LayerFiles[modifiedFile]; found {
						previousFileInfo.Mode = pfi.Mode
						break
					}
				}
				setUIDRemoved := previousFileInfo.Mode&fs.ModeSetuid > 0 && modifiedFileInfo.Mode&fs.ModeSetuid == 0
				setGIDRemoved := previousFileInfo.Mode&fs.ModeSetgid > 0 && modifiedFileInfo.Mode&fs.ModeSetgid == 0

				// Something in the mode changed. The only thing we support is removal of setuid/setgid bits
				if setUIDRemoved || setGIDRemoved {
					logger.V(log.DBG).Info("setuid/setgid bit removed")
					continue
				}

				if !strings.Contains(currentPackage.Release, packageDist) && packageDist != "unknown" {
					// This means it's _probably_ not a RH package. If the file is changed, warn, but don't fail
					logger.Info("WARN: an rpm-installed file was modified outside of rpm, but appears to be from a third-party. This could be a failure in the future")
					continue
				}

				if currentPackage.Vendor != "Red Hat, Inc." && previousPackage.Vendor != "Red Hat, Inc." {
					// This means it's _probably_ not a RH package. If the file is changed, warn, but don't fail
					logger.Info("WARN: an rpm-installed file was modified outside of rpm, but appears to be from a third-party. This could be a failure in the future")
					continue
				}

				if currentPackage.InstallTime > previousPackage.InstallTime {
					// This _probably_ means that the package was either:
					// a) explicitly rpm -e then rpm -i
					// b) dnf reinstall
					// This should not trigger. Going to trace log this, but not always report
					logger.V(log.TRC).Info("package appears to have been re-installed or removed and installed in the same layer", "package", currentPackage.Name)
					continue
				}

				// Nope, nope, nope. File was modified without using RPM
				logger.Info("found disallowed modification in layer", "file", modifiedFile)
				disallowedModifications = true
				continue
			}

			// Check that release contains the same arch, this is to ensure that a package did not get rebuilt with
			// a different architecture
			previousOsRelease := strings.Contains(previousPackage.Release, packageDist)
			currentOsRelease := strings.Contains(currentPackage.Release, packageDist)

			if previousOsRelease && !currentOsRelease {
				logger.Info("mismatch in OS release", "file", modifiedFile)
				disallowedModifications = true
				continue
			}

			// Check that the architectures for previous version and current version of a given package match
			if previousPackage.Arch != currentPackage.Arch {
				logger.Info("mismatch in package architecture", "file", modifiedFile)
				disallowedModifications = true
				continue
			}

			// This appears like an update. This is allowed.
			// No further action required
		}
	}
	return !disallowedModifications, nil
}

func (p HasModifiedFilesCheck) Name() string {
	return "HasModifiedFiles"
}

func (p HasModifiedFilesCheck) Help() check.HelpText {
	return check.HelpText{
		Message:    "Check HasModifiedFiles encountered an error. Please review the preflight.log file for more information.",
		Suggestion: "Do not modify any files installed by RPM in the base Red Hat layer",
	}
}

func (p HasModifiedFilesCheck) Metadata() check.Metadata {
	return check.Metadata{
		Description:      "Checks that no files installed via RPM in the base Red Hat layer have been modified",
		Level:            "best",
		KnowledgeBaseURL: certDocumentationURL,
		CheckURL:         certDocumentationURL,
	}
}

func extractPackageNameVersionRelease(pkgList []*rpmdb.PackageInfo) map[string]packageMeta {
	pkgNameList := make(map[string]packageMeta, len(pkgList))
	for _, pkg := range pkgList {
		pkgNameList[strings.Join([]string{pkg.Name, pkg.Version, pkg.Release, pkg.Arch}, "-")] = packageMeta{
			Name:        pkg.Name,
			Version:     pkg.Version,
			Release:     pkg.Release,
			Arch:        pkg.Arch,
			Vendor:      pkg.Vendor,
			InstallTime: pkg.InstallTime,
		}
	}
	return pkgNameList
}

// directoryIsExcluded excludes a directory and any file contained in that directory.
func directoryIsExcluded(ctx context.Context, s string) bool {
	excl := map[string]struct{}{
		"etc":               {},
		"var":               {},
		"run":               {},
		"usr/lib/.build-id": {},
		"usr/tmp":           {},
	}

	for k := range excl {
		if strings.HasPrefix(s, filepath.Clean(k+"/")) || k == s {
			logger := logr.FromContextOrDiscard(ctx)
			logger.V(log.TRC).Info("directory excluded", "directory", s)
			return true
		}
	}

	return false
}

// pathIsExcluded checks if s is excluded explicitly as written.
func pathIsExcluded(ctx context.Context, s string) bool {
	excl := map[string]struct{}{
		"etc/resolv.conf": {},
		"etc/hostname":    {},
		// etc and etc/ are both required as both can present the directory
		// in a tarball. Same goes for other directories.
		"etc":  {},
		"etc/": {},
		"run":  {},
		"run/": {},
	}

	_, found := excl[s]
	if found {
		logger := logr.FromContextOrDiscard(ctx)
		logger.V(log.TRC).Info("file excluded", "file", s)
		return true
	}
	return found
}

// prefixAndSuffixIsExcluded will check both start and end of path
func prefixAndSuffixIsExcluded(ctx context.Context, s string) bool {
	excl := []struct {
		Prefix string
		Suffix string
	}{
		{Prefix: "usr/", Suffix: ".cache"},
	}

	for _, v := range excl {
		if strings.HasPrefix(s, v.Prefix) && strings.HasSuffix(s, v.Suffix) {
			logger := logr.FromContextOrDiscard(ctx)
			logger.V(log.TRC).Info("prefix and suffix excluded", "filename", s, "prefix", v.Prefix, "suffix", v.Suffix)
			return true
		}
	}

	return false
}

// normalize will clean a filepath of extraneous characters like ./, //, etc.
// and strip a leading slash. E.g. /foo/../baz --> baz
func normalize(s string) string {
	// for the root path, return the root path.
	if s == "/" {
		return s
	}
	return filepath.Clean(strings.TrimPrefix(s, "/"))
}

// installedFileMapWithExclusions gets a map of installed filenames that have been cleaned
// of extra slashes, dotslashes, and leading slashes.
func installedFileMapWithExclusions(ctx context.Context, pkglist []*rpmdb.PackageInfo) (map[string]string, error) {
	const okFlags = rpmdb.RPMFILE_CONFIG |
		rpmdb.RPMFILE_DOC |
		rpmdb.RPMFILE_LICENSE |
		rpmdb.RPMFILE_MISSINGOK |
		rpmdb.RPMFILE_README |
		rpmdb.RPMFILE_ARTIFACT |
		rpmdb.RPMFILE_GHOST
	m := map[string]string{}
	for _, pkg := range pkglist {
		files, err := pkg.InstalledFiles()
		if err != nil {
			return m, err
		}

		// converting directories to a map so we can filter them out quicker
		pkgDirNamesMap := make(map[string]struct{})
		for _, dir := range pkg.DirNames {
			pkgDirNamesMap[dir] = struct{}{}
		}

		for _, file := range files {
			if _, found := pkgDirNamesMap[file.Path]; found {
				// The file is a directory. Skip it.
				continue
			}

			if int32(file.Flags)&okFlags > 0 {
				// It is one of the ok flags. Skip it.
				continue
			}

			normalized := normalize(file.Path)
			if pathIsExcluded(ctx, normalized) || directoryIsExcluded(ctx, normalized) || prefixAndSuffixIsExcluded(ctx, normalized) {
				// It is either an explicitly excluded path or directory. Skip it.
				continue
			}

			// checking to see if the file is already in the map.
			// check to see if all attributes of the rpm match except architecture.
			// this is to support cross architecture file ownership,
			// the 2nd architecture we encounter, we can skip it.
			if val, found := m[normalized]; found {
				s := strings.Split(val, "-")
				name, version, release, arch := s[0], s[1], s[2], s[3]

				if name == pkg.Name && version == pkg.Version && release == pkg.Release && arch != pkg.Arch {
					continue
				}
			}

			m[normalized] = strings.Join([]string{pkg.Name, pkg.Version, pkg.Release, pkg.Arch}, "-")
		}
	}

	return m, nil
}

type fileInfo struct {
	Mode os.FileMode
}
