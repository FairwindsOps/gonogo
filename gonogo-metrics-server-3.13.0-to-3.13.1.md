# Go/No-Go Analysis: metrics-server 3.13.0 → 3.13.1

**Date:** May 1, 2026  
**Component:** metrics-server (Helm Chart)  
**Current Version:** 3.13.0  
**Target Version:** 3.13.1  
**Status:** ❌ **NO-GO**

---

## Executive Summary

**The requested upgrade from metrics-server 3.13.0 to 3.13.1 is not possible because version 3.13.1 does not exist.**

After thorough investigation of official sources, the latest available version of the metrics-server Helm chart is **3.13.0**, released on July 22, 2025.

---

## Investigation Details

### Sources Verified

1. **GitHub Repository (Official)**
   - Repository: https://github.com/kubernetes-sigs/metrics-server
   - Changelog: https://github.com/kubernetes-sigs/metrics-server/blob/master/charts/metrics-server/CHANGELOG.md
   - Latest Helm chart version: **3.13.0**
   - Latest metrics-server application version: **v0.8.1** (released Jan 29, 2026)

2. **Artifact Hub (Official Registry)**
   - URL: https://artifacthub.io/packages/helm/metrics-server/metrics-server
   - Latest chart version: **3.13.0** (Jul 22, 2025)
   - No 3.13.1 version listed

3. **Release History**
   - 3.13.0 (Jul 22, 2025) ← Current latest
   - 3.12.2 (Oct 7, 2024)
   - 3.12.1 (Apr 5, 2024)
   - 3.12.0 (Feb 7, 2024)

### What's in metrics-server 3.13.0?

The current latest version (3.13.0) includes:

- **Security enhancements**: Chart options to secure the connection between Metrics Server and the Kubernetes API Server (#1288)
- **PDB improvements**: `unhealthyPodEvictionPolicy` in PodDisruptionBudget as a user-enabled feature (#1574)
- **Image updates**:
  - Addon Resizer OCI image: `1.8.23` (#1626)
  - Metrics Server OCI image: `v0.8.0` (#1683)

### Latest Application Version

While the Helm chart is at 3.13.0, the metrics-server application itself has a newer version:
- **v0.8.1** (released Jan 29, 2026)

This suggests that a future Helm chart release (potentially 3.13.1 or 3.14.0) may be planned to bundle the v0.8.1 application, but it has not been released yet as of May 1, 2026.

---

## Recommendation

### ❌ NO-GO for 3.13.1 Upgrade

**Reason:** Version 3.13.1 does not exist.

### Alternative Actions

1. **Stay on 3.13.0**: The current version is the latest available and was released less than a year ago.

2. **Monitor for future releases**: Watch the official repository for:
   - Chart version 3.13.1 (if released)
   - Chart version 3.14.0 (next major release)
   - These may include the metrics-server v0.8.1 application

3. **Manual upgrade option**: If you need metrics-server v0.8.1 application features, you could:
   - Use Helm chart 3.13.0 with a custom `image.tag` override to `v0.8.1`
   - **Warning:** This is not officially tested and may have compatibility issues

### Monitoring Resources

- GitHub Releases: https://github.com/kubernetes-sigs/metrics-server/releases
- Helm Chart Changelog: https://github.com/kubernetes-sigs/metrics-server/blob/master/charts/metrics-server/CHANGELOG.md
- Artifact Hub: https://artifacthub.io/packages/helm/metrics-server/metrics-server

---

## Version Matrix

| Component | Current | Target | Status |
|-----------|---------|--------|--------|
| Helm Chart | 3.13.0 | 3.13.1 | ❌ Does not exist |
| Application | v0.8.0 | - | v0.8.1 available but not in chart |
| Chart Release Date | Jul 22, 2025 | N/A | - |

---

## Conclusion

The upgrade to metrics-server Helm chart version 3.13.1 **cannot proceed** as this version has not been released by the kubernetes-sigs project. The current version 3.13.0 remains the latest stable release and should be maintained until a newer version is officially published.

If version 3.13.1 was mentioned in documentation or planning materials, it may have been:
- A typo or misunderstanding
- A planned but unreleased version
- Confusion with the application version (v0.8.1)

**Recommendation: Verify the source of the 3.13.1 version reference and confirm the actual upgrade target.**
