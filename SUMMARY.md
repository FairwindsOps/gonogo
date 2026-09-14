# GoNoGo Assessment Summary: cert-manager 1.13.0 → 1.14.x

## Quick Decision: ✅ GO (with version constraints)

---

## Critical Requirements

### ✅ DO: Upgrade to v1.14.4 or later
**Recommended:** v1.14.6 (latest stable as of May 2026)

### 🛑 DON'T: Use v1.14.0, v1.14.1, v1.14.2, or v1.14.3
**Reason:** Known critical bugs (cainjector image, SAN issues, cmctl bugs)

---

## Prerequisites

| Check | Requirement | Status |
|-------|-------------|--------|
| Kubernetes Version | 1.24 - 1.31 | Verify with `kubectl version` |
| cert-manager EOL | 1.13.0 EOL: June 5, 2024 | ⚠️ Upgrade needed |
| Air-gapped Env | Pre-pull startupapicheck image | Required if air-gapped |
| Backup | Configuration backed up | Run before upgrade |

---

## Breaking Changes Summary

### 🔴 HIGH: Startupapicheck Image
- **Old:** Uses `ctl` image
- **New:** Uses `startupapicheck` image
- **Impact:** Air-gapped environments must pre-pull new image
- **Image:** `quay.io/jetstack/cert-manager-startupapicheck:v1.14.6`

### 🟡 MEDIUM: CSR Encoding
- **Change:** KeyUsage/BasicConstraints now marked as critical
- **Impact:** May affect custom CA integrations
- **Action:** Test in staging first

### 🟢 LOW: Cluster Autoscaler Annotation
- **Change:** ACME pods now have `safe-to-evict: true`
- **Impact:** Pods may be evicted during scaling
- **Action:** Override if needed

---

## Installation Commands

### Quick Upgrade (Helm)
```bash
helm repo update
helm upgrade cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.14.6 \
  --reuse-values
```

### Verification
```bash
kubectl -n cert-manager get pods
kubectl get certificates -A
```

### Rollback (if needed)
```bash
helm rollback cert-manager -n cert-manager
```

---

## Using the GoNoGo Bundle

### Run Assessment
```bash
# From the gonogo repository root
gonogo check -b pkg/bundle/bundles/cert-manager.yaml
```

### What It Checks
- ✓ Kubernetes version compatibility (1.24-1.31)
- ✓ Required API versions
- ✓ Deprecated annotations on Ingresses
- ✓ Deprecated annotations on Certificates
- ✓ Outdated startupapicheck jobs

---

## Risk Level: LOW-MEDIUM ⚠️

| Risk Factor | Level | Mitigation |
|-------------|-------|-----------|
| Version bugs (v1.14.0-v1.14.3) | HIGH → LOW | Use v1.14.4+ |
| Air-gapped image issue | MEDIUM | Pre-pull image |
| CSR compatibility | LOW-MEDIUM | Test staging |
| Upgrade downtime | LOW | Zero-downtime upgrade |
| Certificate disruption | LOW | Auto-renewal |

---

## Timeline

### Immediate Actions
1. ✓ Created GoNoGo bundle
2. ✓ Documented assessment
3. Run bundle check against your cluster
4. Review action items from check

### Before Upgrade
1. Backup configurations
2. Verify K8s version (≥1.24)
3. Pre-pull startupapicheck image (air-gapped only)
4. Test in non-production

### During Upgrade
1. Run Helm upgrade to v1.14.6
2. Monitor pod status
3. Verify certificate renewals

### After Upgrade
1. Confirm all pods healthy
2. Test certificate issuance
3. Monitor webhook/cainjector logs

---

## Files Created

1. **`pkg/bundle/bundles/cert-manager.yaml`** - GoNoGo bundle specification
2. **`CERT_MANAGER_1.13_TO_1.14_GONOGO.md`** - Comprehensive assessment (15 pages)
3. **`SUMMARY.md`** - This quick reference (2 pages)

---

## Resources

- **Full Assessment:** [`CERT_MANAGER_1.13_TO_1.14_GONOGO.md`](./CERT_MANAGER_1.13_TO_1.14_GONOGO.md)
- **Bundle:** [`pkg/bundle/bundles/cert-manager.yaml`](./pkg/bundle/bundles/cert-manager.yaml)
- **Official Docs:** https://cert-manager.io/docs/releases/upgrading/upgrading-1.13-1.14/
- **Release Notes:** https://github.com/cert-manager/cert-manager/releases/tag/v1.14.6
- **Pull Request:** https://github.com/FairwindsOps/gonogo/pull/201

---

## Bottom Line

**The cert-manager 1.13.0 → 1.14.x upgrade is APPROVED** when targeting v1.14.4 or later. Main concerns are version selection and air-gapped image management. With proper planning, this is a low-risk upgrade that brings important security updates and new features.

**Next Steps:**
1. Run the GoNoGo bundle check
2. Address any action items
3. Test in staging
4. Proceed with production upgrade to v1.14.6

---

*Assessment completed: May 1, 2026*  
*Target versions: cert-manager 1.13.0 → 1.14.6*  
*Kubernetes compatibility: 1.24 - 1.31*
