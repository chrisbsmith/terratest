package helm

import (
	"path/filepath"

	"github.com/gruntwork-io/gruntwork-cli/errors"
	"github.com/stretchr/testify/require"

	"github.com/chrisbsmith/terratest/modules/files"
	"github.com/chrisbsmith/terratest/modules/testing"
)

// Dependency will manage the dependencies of a parent chart. The test will fail if there is an error
func Dependency(t testing.TestingT, options *Options, chart string, releaseName string) {
	require.NoError(t, DependencyE(t, options, chart, releaseName))
}

// DependencyE will install the selected helm chart with the provided options under the given release name.
func DependencyE(t testing.TestingT, options *Options, chart string, releaseName string) error {
	// If the chart refers to a path, convert to absolute path. Otherwise, pass straight through as it may be a remote
	// chart.
	if files.FileExists(chart) {
		absChartDir, err := filepath.Abs(chart)
		if err != nil {
			return errors.WithStackTrace(err)
		}
		chart = absChartDir
	}

	// Now call out to helm dependency to install the charts with the provided options
	// Declare err here so that we can update args later
	var err error
	args := []string{}
	args = append(args, getNamespaceArgs(options)...)
	if options.Version != "" {
		args = append(args, "--version", options.Version)
	}
	args, err = getValuesArgsE(t, options, args...)
	if err != nil {
		return err
	}
	args = append(args, releaseName, chart)
	_, err = RunHelmCommandAndGetOutputE(t, options, "dependency", args...)
	return err
}
