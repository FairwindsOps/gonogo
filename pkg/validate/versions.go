package validate

import (
	"github.com/blang/semver/v4"
	clusterVersion "k8s.io/apimachinery/pkg/version"
	"k8s.io/klog"
)

func (m *match) validateClusterVersion(cv *clusterVersion.Info) error {

	maxVer, err := semver.ParseTolerant(m.Bundle.CompatibleK8sVersions.Max)
	if err != nil {
		klog.Error(err)
	}

	minVer, err := semver.ParseTolerant(m.Bundle.CompatibleK8sVersions.Min)
	if err != nil {
		klog.Error(err)
	}

	clusterVer, err := semver.ParseTolerant(cv.String())
	if err != nil {
		klog.Error(err)
	}

	if clusterVer.LT(minVer) || clusterVer.GT(maxVer) {
		m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, &ActionItem{
			ResourceNamespace: m.Release.Namespace,
			ResourceName:      m.Release.Name,
			Title:             "Unsupported cluster version",
			Description:       "The Kubernetes cluster version is less or greater than the minimum and maximum versions specified in the bundle spec",
		})

	}
	return nil
}
