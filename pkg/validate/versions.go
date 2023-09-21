package validate

import (
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/thoas/go-funk"
	clusterVersion "k8s.io/apimachinery/pkg/version"
	"k8s.io/klog"
)

func (m *match) validateClusterVersion(cv *clusterVersion.Info) error {
	var maxVer, minVer, clusterVer semver.Version
	var err error

	clusterVer, err = semver.ParseTolerant(cv.String())
	if err != nil {
		return err
	}

	if len(m.Bundle.CompatibleK8sVersions.Max) != 0 {
		maxVer, err = semver.ParseTolerant(m.Bundle.CompatibleK8sVersions.Max)
		if err != nil {
			return err
		}
		if clusterVer.GT(maxVer) {
			m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, &ActionItem{
				ResourceNamespace: m.Release.Namespace,
				ResourceName:      m.Release.Name,
				Title:             "Unsupported cluster version",
				Description:       "The Kubernetes cluster version is greater than the maximum version specified in the bundle spec",
			})
		}
	}

	if len(m.Bundle.CompatibleK8sVersions.Min) != 0 {
		minVer, err = semver.ParseTolerant(m.Bundle.CompatibleK8sVersions.Min)
		if err != nil {
			return err
		}
		if clusterVer.LT(minVer) {
			m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, &ActionItem{
				ResourceNamespace: m.Release.Namespace,
				ResourceName:      m.Release.Name,
				Title:             "Unsupported cluster version",
				Description:       "The Kubernetes cluster version is less than the minimum version specified in the bundle spec",
			})
		}
	}
	return nil
}

func (m *match) validateAPIVersion(v []string) {
	apiVersions := m.Bundle.NecessaryAPIVersions

	for _, av := range apiVersions {
		if !funk.Contains(v, av) {
			m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, &ActionItem{
				ResourceNamespace: m.Release.Namespace,
				ResourceName:      m.Release.Name,
				Title:             fmt.Sprintf("API version %s is not available", av),
				Description:       fmt.Sprintf("The Kubernetes cluster version does not contain the api %s", av),
			})
		} else {
			klog.V(5).Infof("found required apiversion %s", av)
		}
	}

}
