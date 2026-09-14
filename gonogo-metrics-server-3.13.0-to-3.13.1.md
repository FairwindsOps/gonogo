# Go/No-Go Analysis: metrics-server Upgrade with Image Override

**Date:** May 1, 2026  
**Component:** metrics-server (Helm Chart + Application Image)  
**Current Version:** Chart 3.13.0 with image v0.8.0  
**Target Version:** Chart 3.13.0 with image v0.8.1 (override)  
**Status:** ⚠️ **CONDITIONAL GO**

---

## Executive Summary

**Helm chart version 3.13.1 does not exist**, but we can upgrade the metrics-server application from v0.8.0 to v0.8.1 by using Helm chart 3.13.0 with an image override.

### Upgrade Approach
- **Helm Chart:** Stay on 3.13.0 (latest available)
- **Application Image:** Override to `registry.k8s.io/metrics-server/metrics-server:v0.8.1`
- **Rationale:** Security updates (Golang 1.24.12, Kubernetes dependencies v0.33.7)

---

## Investigation Details

### Version Availability

1. **Helm Chart Versions (Official Sources)**
   - Repository: https://github.com/kubernetes-sigs/metrics-server
   - Artifact Hub: https://artifacthub.io/packages/helm/metrics-server/metrics-server
   - **Latest chart:** 3.13.0 (released Jul 22, 2025)
   - **Chart 3.13.1:** Does not exist
   
2. **Application Image Versions**
   - Chart 3.13.0 default: `registry.k8s.io/metrics-server/metrics-server:v0.8.0`
   - Latest available: `v0.8.1` (released Jan 29, 2026)
   - Image repository: https://github.com/kubernetes-sigs/metrics-server/releases

### What's New in v0.8.1 (Application Image)

**Release Date:** January 29, 2026  
**Changes since v0.8.0:**

#### Dependency Updates
- ✅ **Golang:** Upgraded from 1.24.11 → **1.24.12**
  - Fixes 6 CVEs including:
    - CVE-2025-61728: archive/zip denial of service
    - CVE-2025-61729: net/http memory exhaustion
    - CVE-2025-68121: crypto/tls session resumption issues
    - CVE-2025-61730: crypto/tls handshake encryption level issues
    - Multiple crypto/tls security improvements

- ✅ **Kubernetes dependencies:** Upgraded to **v0.33.7**
  - Includes security patches and bug fixes
  - Better compatibility with Kubernetes 1.33+

### Security Considerations

#### ✅ Fixes from v0.8.0
The v0.8.0 release (bundled in chart 3.13.0) already fixed:
- CVE-2024-45337 (Critical) - golang.org/x/crypto
- CVE-2024-45338 (High) - golang.org/x/net
- CVE-2024-34158, CVE-2024-34156

#### ⚠️ Known Issues in v0.8.1
As of May 2026, v0.8.1 has reported CVEs in dependencies:
- **CVE-2026-33186** (Critical) - google.golang.org/grpc@v1.72.0
  - Fixed in: grpc v1.79.3
  - Impact: Incorrect Authorization
- **CVE-2026-24051** (High) - go.opentelemetry.io/otel/sdk/resource@v1.35.0
  - Fixed in: otel v1.40.0
  - Impact: Untrusted Search Path

**Note:** These CVEs are in dependency libraries and may not be exploitable in metrics-server's specific usage. However, they will be flagged by security scanners.

---

## Recommendation

### ⚠️ CONDITIONAL GO - Image Override to v0.8.1

Given the requirement to upgrade the application to v0.8.1, we recommend using Helm chart 3.13.0 with an image override.

#### Deployment Method

```yaml
# values.yaml override
image:
  repository: registry.k8s.io/metrics-server/metrics-server
  tag: v0.8.1
```

Or via Helm command:
```bash
helm upgrade metrics-server metrics-server/metrics-server \
  --version 3.13.0 \
  --set image.tag=v0.8.1
```

### Risk Assessment

| Risk Factor | Level | Details |
|-------------|-------|---------|
| **Compatibility** | 🟢 Low | v0.8.1 is a minor patch release with no breaking changes |
| **Testing** | 🟡 Medium | Image override not officially bundled/tested in chart 3.13.0 |
| **Security** | 🟡 Medium | Fixes Golang CVEs but introduces new dependency CVEs |
| **Support** | 🟢 Low | Both chart and image are from official kubernetes-sigs project |
| **Rollback** | 🟢 Low | Easy to rollback to v0.8.0 if issues arise |

### Pros ✅

1. **Security improvements:**
   - Golang 1.24.12 fixes 6 CVEs from previous versions
   - Updated Kubernetes dependencies (v0.33.7)
   - Resolves CVE scanner alerts from v0.8.0

2. **Minimal risk:**
   - Only dependency version bumps, no functional changes
   - Official image from kubernetes-sigs
   - Helm chart 3.13.0 is stable and tested

3. **Better K8s compatibility:**
   - Improved support for Kubernetes 1.33+
   - More recent dependency versions

### Cons ⚠️

1. **Not officially bundled:**
   - Image v0.8.1 not officially tested with chart 3.13.0
   - Requires manual image override

2. **New CVEs present:**
   - CVE-2026-33186 (Critical) in grpc dependency
   - CVE-2026-24051 (High) in otel dependency
   - Security scanners will flag these

3. **Future chart release:**
   - A new chart version (3.13.1 or 3.14.0) may be released soon
   - Would require another upgrade to stay on official versions

### Alternative Options

#### Option 1: Image Override to v0.8.1 (Recommended)
- **Use case:** Need Golang CVE fixes, acceptable to have new dependency CVEs
- **Implementation:** Chart 3.13.0 + image.tag override to v0.8.1
- **Timeline:** Can deploy immediately

#### Option 2: Stay on v0.8.0
- **Use case:** Risk-averse, prefer officially bundled versions
- **Implementation:** Chart 3.13.0 + default image v0.8.0
- **Timeline:** Stay current until chart 3.13.1/3.14.0 is released
- **Trade-off:** Keep Golang CVEs, avoid new dependency CVEs

#### Option 3: Wait for Official Chart Release
- **Use case:** Need fully tested, officially bundled solution
- **Implementation:** Wait for chart 3.13.1 or 3.14.0
- **Timeline:** Unknown (could be weeks or months)
- **Trade-off:** Delay getting security updates

### Testing Recommendations

Before deploying to production:

1. **Deploy to non-production environment**
   - Verify metrics collection works
   - Check HPA (Horizontal Pod Autoscaler) functionality
   - Monitor for errors in logs

2. **Run security scans**
   - Confirm which CVEs are actually present
   - Assess exploitability in your environment

3. **Validate compatibility**
   - Test with your Kubernetes version
   - Verify RBAC permissions still work
   - Check monitoring/alerting integrations

4. **Performance testing**
   - Monitor resource usage
   - Verify metrics accuracy
   - Check API response times

### Monitoring After Deployment

- Watch for new chart releases: https://github.com/kubernetes-sigs/metrics-server/releases
- Monitor security advisories: https://github.com/kubernetes-sigs/metrics-server/security/advisories
- Track CVE fixes: Issues #1780 (CVE-2026-33186) and #1774 (CVE-2026-24051)

---

## Version Matrix

| Component | Current | Proposed | Notes |
|-----------|---------|----------|-------|
| **Helm Chart** | 3.13.0 | 3.13.0 | Latest available (no 3.13.1) |
| **Application Image** | v0.8.0 | v0.8.1 | Via image override |
| **Golang Version** | 1.24.11 | 1.24.12 | Fixes 6 CVEs |
| **K8s Dependencies** | v0.33.x | v0.33.7 | Security & bug fixes |
| **Chart Release Date** | Jul 22, 2025 | Jul 22, 2025 | Unchanged |
| **Image Release Date** | Jul 22, 2025 | Jan 29, 2026 | Official release |

### CVE Status

| CVE | Version | Status | Severity |
|-----|---------|--------|----------|
| CVE-2024-45337 | v0.8.0+ | ✅ Fixed | Critical |
| CVE-2024-45338 | v0.8.0+ | ✅ Fixed | High |
| CVE-2024-34158 | v0.8.0+ | ✅ Fixed | High |
| CVE-2025-61728 | v0.8.1+ | ✅ Fixed | High |
| CVE-2025-61729 | v0.8.1+ | ✅ Fixed | High |
| CVE-2025-68121 | v0.8.1+ | ✅ Fixed | Medium |
| CVE-2026-33186 | v0.8.1 | ⚠️ Present | Critical |
| CVE-2026-24051 | v0.8.1 | ⚠️ Present | High |

---

## Implementation Guide

### Step 1: Backup Current Configuration

```bash
# Export current metrics-server configuration
helm get values metrics-server -n kube-system > metrics-server-values-backup.yaml

# Export current deployment
kubectl get deployment metrics-server -n kube-system -o yaml > metrics-server-deployment-backup.yaml
```

### Step 2: Update Values File

Create or update your `values.yaml`:

```yaml
# values.yaml
image:
  repository: registry.k8s.io/metrics-server/metrics-server
  tag: v0.8.1
  pullPolicy: IfNotPresent

# Keep existing chart 3.13.0 configurations
# (add your existing custom values here)
```

### Step 3: Perform Upgrade

```bash
# Update Helm repo
helm repo update

# Dry-run first
helm upgrade metrics-server metrics-server/metrics-server \
  --version 3.13.0 \
  --namespace kube-system \
  --values values.yaml \
  --dry-run --debug

# Perform actual upgrade
helm upgrade metrics-server metrics-server/metrics-server \
  --version 3.13.0 \
  --namespace kube-system \
  --values values.yaml

# Verify deployment
kubectl rollout status deployment/metrics-server -n kube-system
```

### Step 4: Validation

```bash
# Check pod is running with new image
kubectl get pod -n kube-system -l app.kubernetes.io/name=metrics-server
kubectl describe pod -n kube-system -l app.kubernetes.io/name=metrics-server | grep "Image:"

# Verify metrics are being collected
kubectl top nodes
kubectl top pods -A

# Check metrics-server logs
kubectl logs -n kube-system -l app.kubernetes.io/name=metrics-server --tail=50

# Test HPA functionality (if applicable)
kubectl get hpa -A
```

### Step 5: Rollback (if needed)

```bash
# Quick rollback to previous version
helm rollback metrics-server -n kube-system

# Or explicitly set back to v0.8.0
helm upgrade metrics-server metrics-server/metrics-server \
  --version 3.13.0 \
  --namespace kube-system \
  --set image.tag=v0.8.0 \
  --reuse-values
```

---

## Conclusion

### Decision: ⚠️ CONDITIONAL GO

The upgrade to metrics-server application v0.8.1 via image override is **recommended** with the following conditions:

✅ **Proceed if:**
- Security scanning alerts for Golang CVEs are a concern
- You have a non-production environment for testing
- You can tolerate new dependency CVEs (CVE-2026-33186, CVE-2026-24051)
- Rollback capability is available

⚠️ **Proceed with caution if:**
- You have strict security scanner requirements (will flag new CVEs)
- Production-only deployment without testing environment
- Limited rollback capabilities

❌ **Do not proceed if:**
- Zero-tolerance policy for Critical CVEs in dependencies
- Cannot accept unofficial image/chart combinations
- Prefer to wait for officially bundled chart release

### Final Recommendation

**Deploy to staging/test environment first**, validate functionality, assess security scanner results, then proceed to production if acceptable. The security improvements from Golang 1.24.12 likely outweigh the risk from the new dependency CVEs, which may not be exploitable in metrics-server's usage context.

### Next Steps

1. Review this analysis with security team
2. Test image override in non-production environment
3. Run security scans on v0.8.1 image
4. Assess CVE exploitability in your environment
5. Make go/no-go decision based on your risk tolerance
6. Monitor for official chart release (3.13.1 or 3.14.0)

---

## References

- **Helm Chart Repository:** https://github.com/kubernetes-sigs/metrics-server/tree/master/charts/metrics-server
- **Release v0.8.1:** https://github.com/kubernetes-sigs/metrics-server/releases/tag/v0.8.1
- **Release v0.8.0:** https://github.com/kubernetes-sigs/metrics-server/releases/tag/v0.8.0
- **CVE-2026-33186 Issue:** https://github.com/kubernetes-sigs/metrics-server/issues/1780
- **CVE-2026-24051 Issue:** https://github.com/kubernetes-sigs/metrics-server/issues/1774
- **Golang 1.24.12 Release Notes:** https://go.dev/doc/devel/release#go1.24.12
