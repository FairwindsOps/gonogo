# ArgoCD Helm Chart Upgrade: Go/No-Go Assessment
## Version: 9.4.15 → 9.5.6

**Generated:** 2026-05-01  
**Assessment Status:** ✅ **GO** with considerations

---

## Executive Summary

The upgrade from ArgoCD Helm chart version 9.4.15 to 9.5.6 is assessed as **LOW RISK** and recommended to proceed. This upgrade includes:

- **Application Version**: v3.3.4 → v3.3.8 (4 patch releases)
- **Chart Changes**: 47 commits across 50 files
- **Breaking Changes**: None identified
- **New Features**: VPA support, HTTPRoute support for ApplicationSet
- **Kubernetes Requirements**: >= 1.25.0-0 (unchanged)

---

## Change Analysis

### Application Version Changes (v3.3.4 → v3.3.8)

The chart upgrade includes 4 ArgoCD application patch releases:
- v3.3.5 (released March 25, 2026)
- v3.3.6 (released March 27, 2026)
- v3.3.7 (released April 17, 2026)
- v3.3.8 (released April 21, 2026)

**Recommendation**: Review the [ArgoCD release notes](https://github.com/argoproj/argo-cd/releases) for these versions to identify any security patches or important bug fixes.

### Major Chart Changes

#### 1. VerticalPodAutoscaler (VPA) Support ✨ NEW
- **PR**: #3817
- **Impact**: VPA resources added for all ArgoCD components
- **Components affected**: 
  - Application Controller
  - Repo Server
  - Server
  - Dex
  - Redis
  - ApplicationSet
  - Commit Server
  - Notifications Controller
- **Risk**: LOW - Disabled by default
- **Action Required**: 
  - If enabling VPA: Ensure VPA CRDs are installed in your cluster
  - Install from: https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler
  - Test values configuration files for VPA settings if you plan to enable

#### 2. HTTPRoute Support for ApplicationSet Webhook ✨ NEW
- **PR**: #3859
- **Impact**: Gateway API HTTPRoute support added for ApplicationSet webhook
- **Risk**: LOW - Opt-in feature
- **Action Required**:
  - If using Gateway API: Verify Gateway configuration compatibility
  - Review new values under `applicationSet.httproute`

#### 3. VPA updateMode Quoting Fix 🐛 FIXED
- **PR**: #3790
- **Impact**: Fixed YAML boolean coercion where "Off" was being converted to `false`
- **Risk**: VERY LOW - Bug fix that improves correctness
- **Action Required**: None - This is a fix for existing functionality

#### 4. PrometheusRule Annotation Fix 🐛 FIXED
- **PR**: #3856
- **Impact**: Fixed template syntax in ArgoAppNotSynced PrometheusRule annotation
- **Risk**: VERY LOW - Template syntax correction
- **Action Required**: None if not using PrometheusRule

#### 5. Repo Server Git Copy Optimization ⚡ OPTIMIZATION
- **PR**: #3834
- **Impact**: Added `repoServer.copyutil.extraArgs` with default `--update=none`
- **Purpose**: Optimizes git operations by skipping unnecessary updates
- **Risk**: LOW - Performance optimization
- **Action Required**: 
  - Review if you have custom git workflows that might be affected
  - Can override with custom `repoServer.copyutil.extraArgs` if needed

#### 6. Extension Installer Image Update 📦 DEPENDENCY UPDATE
- **Change**: quay.io/argoprojlabs/argocd-extension-installer:v0.0.9 → v1.0.0
- **Risk**: MEDIUM - Major version bump (but pre-1.0 to 1.0)
- **Action Required**:
  - If using ArgoCD UI extensions: Test compatibility with v1.0.0
  - Review extension documentation for any breaking changes

---

## Compatibility Assessment

### ✅ Kubernetes Version Requirements
- **Required**: >= 1.25.0-0
- **Status**: Unchanged from 9.4.15
- **Action**: Verify your cluster is running Kubernetes 1.25 or later

### ✅ Required API Versions
- `apps/v1` - Deployments, StatefulSets
- `v1` - Core resources (ConfigMaps, Secrets, Services)
- `networking.k8s.io/v1` - Ingress, NetworkPolicy
- **Status**: Standard APIs, should be available in all Kubernetes 1.25+ clusters

### ⚠️ Optional API Versions (for new features)
- `autoscaling.k8s.io/v1` - VerticalPodAutoscaler (if enabling VPA)
- `gateway.networking.k8s.io/v1` - HTTPRoute (if using Gateway API)

---

## Security Considerations

### 1. Patch Release Updates
The upgrade includes 4 patch releases of the ArgoCD application (v3.3.4 → v3.3.8). While these are typically bug fixes and security patches, you should:

**Action Items**:
- [ ] Review [ArgoCD v3.3.5 release notes](https://github.com/argoproj/argo-cd/releases/tag/v3.3.5)
- [ ] Review [ArgoCD v3.3.6 release notes](https://github.com/argoproj/argo-cd/releases/tag/v3.3.6)
- [ ] Review [ArgoCD v3.3.7 release notes](https://github.com/argoproj/argo-cd/releases/tag/v3.3.7)
- [ ] Review [ArgoCD v3.3.8 release notes](https://github.com/argoproj/argo-cd/releases/tag/v3.3.8)
- [ ] Check if any CVEs are addressed in these releases

### 2. Workflow Security Hardening
The chart repository itself has added step-security/harden-runner to all GitHub workflow jobs, indicating improved security posture in the chart development process.

---

## Testing Recommendations

### Pre-Upgrade Testing

#### 1. Non-Production Environment Testing
Deploy the upgrade to a non-production environment first:

```bash
# Add the ArgoCD Helm repository
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

# Review the new values
helm show values argo/argo-cd --version 9.5.6 > values-9.5.6.yaml

# Perform a dry-run
helm upgrade argocd argo/argo-cd \
  --version 9.5.6 \
  --namespace argocd \
  --values your-values.yaml \
  --dry-run --debug

# Apply to test cluster
helm upgrade argocd argo/argo-cd \
  --version 9.5.6 \
  --namespace argocd \
  --values your-values.yaml
```

#### 2. Functional Testing Checklist
- [ ] Verify all ArgoCD components are running (controller, repo-server, server, dex, redis)
- [ ] Test application sync operations
- [ ] Verify SSO/authentication if configured
- [ ] Test ApplicationSet functionality if in use
- [ ] Verify webhook functionality
- [ ] Check Prometheus metrics collection
- [ ] Test UI extensions if configured
- [ ] Verify notification controllers if configured

#### 3. Values File Review
Compare your current values against the new defaults:

```bash
# Get current values
helm get values argocd -n argocd > current-values.yaml

# Compare with new defaults
helm show values argo/argo-cd --version 9.5.6 > new-defaults.yaml
diff current-values.yaml new-defaults.yaml
```

#### 4. New Configuration Options to Consider

If you want to enable VPA (optional):
```yaml
controller:
  vpa:
    enabled: true
    updateMode: "Initial"  # or "Auto", "Recreate", "Off"

repoServer:
  vpa:
    enabled: true
    updateMode: "Initial"

server:
  vpa:
    enabled: true
    updateMode: "Initial"
```

If you want to enable ApplicationSet HTTPRoute (optional):
```yaml
applicationSet:
  httproute:
    enabled: true
    parentRefs:
      - name: your-gateway
        namespace: gateway-namespace
    hostnames:
      - argocd-applicationset.example.com
```

---

## Upgrade Procedure

### Recommended Upgrade Steps

1. **Backup Current State**
   ```bash
   # Backup current Helm release
   helm get values argocd -n argocd > backup-values.yaml
   helm get manifest argocd -n argocd > backup-manifest.yaml
   
   # Backup ArgoCD applications (optional but recommended)
   kubectl get applications -n argocd -o yaml > backup-applications.yaml
   ```

2. **Update Helm Repository**
   ```bash
   helm repo update argo
   ```

3. **Review Changes**
   ```bash
   helm diff upgrade argocd argo/argo-cd \
     --version 9.5.6 \
     --namespace argocd \
     --values your-values.yaml
   ```
   
   *Note: Requires [helm-diff plugin](https://github.com/databus23/helm-diff)*

4. **Perform Upgrade**
   ```bash
   helm upgrade argocd argo/argo-cd \
     --version 9.5.6 \
     --namespace argocd \
     --values your-values.yaml \
     --wait \
     --timeout 10m
   ```

5. **Verify Deployment**
   ```bash
   # Check all pods are running
   kubectl get pods -n argocd
   
   # Check ArgoCD server version
   kubectl exec -n argocd deployment/argocd-server -- argocd version
   
   # Check application health
   kubectl get applications -n argocd
   ```

### Rollback Procedure (if needed)

If issues are encountered, rollback is straightforward:

```bash
helm rollback argocd -n argocd
```

Or rollback to specific revision:
```bash
# List releases
helm history argocd -n argocd

# Rollback to specific revision
helm rollback argocd [REVISION] -n argocd
```

---

## Risk Assessment Matrix

| Category | Risk Level | Severity | Mitigation |
|----------|-----------|----------|------------|
| Breaking Changes | 🟢 NONE | N/A | None identified |
| Application Patches | 🟡 LOW | Minor | Review release notes for CVEs |
| VPA Support | 🟢 VERY LOW | None | Feature disabled by default |
| HTTPRoute Support | 🟢 VERY LOW | None | Feature disabled by default |
| Extension Installer Update | 🟡 MEDIUM | Low | Test if using extensions |
| Git Copy Optimization | 🟢 LOW | None | Can override if needed |
| Values Schema Changes | 🟢 VERY LOW | None | Backward compatible |

**Overall Risk Assessment**: 🟢 **LOW RISK**

---

## Go/No-Go Decision Checklist

### Pre-Upgrade Checklist

- [ ] Kubernetes cluster version is >= 1.25.0
- [ ] Current ArgoCD installation is healthy and operational
- [ ] Backup of current Helm values completed
- [ ] Backup of ArgoCD applications completed (optional)
- [ ] ArgoCD v3.3.5 - v3.3.8 release notes reviewed
- [ ] No known critical issues with target version
- [ ] Maintenance window scheduled (recommended)
- [ ] Rollback procedure documented and tested
- [ ] Team notified of upcoming upgrade

### Optional Feature Checklist (if enabling)

#### If Enabling VPA:
- [ ] VPA CRDs installed in cluster
- [ ] VPA admission controller running
- [ ] VPA configuration values prepared
- [ ] Resource limits/requests implications understood

#### If Using Gateway API/HTTPRoute:
- [ ] Gateway API CRDs installed
- [ ] Gateway resources configured
- [ ] HTTPRoute configuration prepared

#### If Using Extensions:
- [ ] Extension compatibility with v1.0.0 verified
- [ ] Extension testing completed

### Post-Upgrade Verification

- [ ] All ArgoCD pods are running
- [ ] ArgoCD server is accessible
- [ ] Application sync operations working
- [ ] Webhooks functioning correctly
- [ ] SSO/authentication working (if configured)
- [ ] Monitoring/metrics collection active
- [ ] No error logs in ArgoCD components

---

## Monitoring During Upgrade

### Key Metrics to Watch

1. **Pod Status**
   ```bash
   watch kubectl get pods -n argocd
   ```

2. **Application Health**
   ```bash
   watch kubectl get applications -n argocd
   ```

3. **Component Logs**
   ```bash
   # Application Controller
   kubectl logs -n argocd deployment/argocd-application-controller -f
   
   # Repo Server
   kubectl logs -n argocd deployment/argocd-repo-server -f
   
   # Server
   kubectl logs -n argocd deployment/argocd-server -f
   ```

4. **Resource Usage**
   ```bash
   kubectl top pods -n argocd
   ```

---

## Known Issues & Limitations

### Issues Fixed in This Release
1. ✅ VPA updateMode boolean coercion (PR #3790)
2. ✅ PrometheusRule annotation template syntax (PR #3856)
3. ✅ Duplicate nodePort key in argocd-image-updater (PR #3823)

### Potential Considerations
1. **Extension Installer Major Version**: The extension installer has moved from v0.0.9 to v1.0.0. While this is technically a major version bump, it's moving from pre-release to stable. Review extension documentation if you use custom extensions.

2. **Git Copy Optimization**: The new default `--update=none` flag for copyutil may affect workflows that depend on git update behavior. Override if needed.

---

## References

### Official Documentation
- [ArgoCD Helm Chart Repository](https://github.com/argoproj/argo-helm)
- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [ArgoCD Release Notes](https://github.com/argoproj/argo-cd/releases)

### Relevant Pull Requests
- [PR #3817 - Add VPA support for all components](https://github.com/argoproj/argo-helm/pull/3817)
- [PR #3859 - Add HTTPRoute support to ApplicationSet webhook](https://github.com/argoproj/argo-helm/pull/3859)
- [PR #3790 - Quote VPA updateMode to prevent YAML boolean coercion](https://github.com/argoproj/argo-helm/pull/3790)
- [PR #3834 - Add repoServer.copyutil.extraArgs with default '--update=none'](https://github.com/argoproj/argo-helm/pull/3834)

### Version Comparison
- [Full Changelog: 9.4.15...9.5.6](https://github.com/argoproj/argo-helm/compare/argo-cd-9.4.15...argo-cd-9.5.6)

---

## Conclusion

### Final Recommendation: ✅ **GO**

The upgrade from ArgoCD Helm chart 9.4.15 to 9.5.6 is **recommended to proceed** with the following confidence level:

**Confidence Score: 85/100**

**Reasoning**:
- No breaking changes identified
- Includes important bug fixes (VPA boolean coercion, PrometheusRule templates)
- New features are opt-in and disabled by default
- Four patch releases of the application should be reviewed for security updates
- Backward compatible with existing configurations
- Standard upgrade procedure applies

**Conditions for Approval**:
1. Test in non-production environment first
2. Review ArgoCD application release notes (v3.3.5-v3.3.8)
3. Have rollback plan ready
4. Schedule during maintenance window (recommended but not required)

**Timeline Recommendation**:
- **Testing Phase**: Deploy to staging/dev environment, perform functional testing
- **Production Upgrade**: After successful testing and release note review

---

**Assessment Prepared By**: GoNoGo Cloud Agent  
**Date**: 2026-05-01  
**Bundle Specification**: `argocd-9.4.15-to-9.5.6-bundle.yaml`
