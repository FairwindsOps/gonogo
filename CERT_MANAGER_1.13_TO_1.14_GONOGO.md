# GoNoGo Assessment: cert-manager 1.13.0 to 1.14.0 Upgrade

## Executive Summary

**Recommendation: GO with caution ⚠️**

The cert-manager upgrade from 1.13.0 to 1.14.x is feasible and recommended, but requires careful attention to specific version selection and breaking changes. The 1.13.0 release reached End of Life on June 5, 2024, making the upgrade necessary for continued support.

**CRITICAL: Skip versions 1.14.0, 1.14.1, 1.14.2, and 1.14.3** - Target version 1.14.4 or later (latest stable is v1.14.6).

---

## Version Compatibility

### Kubernetes Version Support

| cert-manager Version | Compatible Kubernetes Versions | Compatible OpenShift Versions |
|---------------------|-------------------------------|------------------------------|
| **1.13.0** (current) | 1.21 → 1.27 | 4.8 → 4.14 |
| **1.14.x** (target) | 1.24 → 1.31 | 4.11 → 4.16 |

**Compatibility Assessment:**
- ✅ Kubernetes 1.24-1.27: Fully compatible with both versions
- ⚠️ Kubernetes 1.21-1.23: Only supported by 1.13.0, upgrade to K8s first
- ✅ Kubernetes 1.28-1.31: New support added in 1.14.x

### Kubernetes Client Libraries
- cert-manager 1.13.0 uses Kubernetes libraries v0.27.4
- cert-manager 1.14.0 maintains compatibility with similar library versions

---

## Critical Breaking Changes

### 1. Startupapicheck Image Change (HIGH PRIORITY)

**Impact:** Air-gapped or restricted environments

**Change:** The startupapicheck job now uses a new OCI image called `startupapicheck` instead of the `ctl` image.

**Action Required:**
- **Air-gapped environments:** Pre-pull and make available:
  ```
  quay.io/jetstack/cert-manager-startupapicheck:v1.14.6
  ```
- **Online environments:** No action needed (automatic pull)

**Detection:** The GoNoGo bundle includes OPA checks to detect outdated startupapicheck jobs

---

### 2. KeyUsage and BasicConstraints Encoding (MEDIUM PRIORITY)

**Impact:** External certificate validation

**Change:** KeyUsage and BasicConstraints extensions are now encoded as **critical** in the CertificateRequest's CSR blob (per RFC 5280 compliance).

**Potential Issues:**
- Strict CAs or validation tools may reject CSRs with different criticality settings
- Most standard CAs (Let's Encrypt, DigiCert, etc.) handle this correctly

**Action Required:**
- Test certificate issuance in staging environment first
- Verify external CA compatibility if using custom/enterprise CAs

---

### 3. ACME HTTP01 Solver Annotation (LOW PRIORITY)

**Impact:** Cluster autoscaler behavior

**Change:** ACME challenge solver Pods now include annotation:
```yaml
cluster-autoscaler.kubernetes.io/safe-to-evict: "true"
```

**Action Required:**
- If you need solver pods to be non-evictable, override in your issuer's `podTemplate`:
```yaml
spec:
  acme:
    solvers:
    - http01:
        ingress:
          podTemplate:
            metadata:
              annotations:
                cluster-autoscaler.kubernetes.io/safe-to-evict: "false"
```

---

## Known Issues (Version-Specific)

### ⛔ v1.14.0
- Helm chart uses wrong cainjector image (installation fails)
- CA and SelfSigned issuers incorrectly copy critical flag from CSR
- cmctl namespace detection bug prevents use outside cert-manager namespace
- cmctl experimental install command panics

**Status:** Fixed in v1.14.1+, but further bugs discovered

### ⛔ v1.14.1
- CA and SelfSigned issuers still have SAN criticality issues

**Status:** Fixed in v1.14.2

### ⛔ v1.14.2 & v1.14.3
- Additional bugs discovered during community testing

**Status:** All known issues resolved in v1.14.4+

### ✅ v1.14.4+ (RECOMMENDED)
- All critical bugs fixed
- JKS and PKCS12 stores now contain full CA set
- cainjector leaderelection configuration corrected
- Go 1.21.8+ (CVE-2024-24783 fixed)

### ✅ v1.14.6 (LATEST STABLE)
- Go 1.21.11 (security fixes for archive/zip and net/netip)
- Additional stability improvements

---

## New Features in 1.14

### Enhancements
1. **X.509 "Other Name" Fields:** Support for creating certificates with Other Name fields
2. **CA Name Constraints:** Support for creating CA certificates with Name Constraints extensions
3. **Authority Information Accessors:** Support for AIA extensions in CA certificates
4. **Security Updates:** Built with Go 1.21.11, addressing multiple CVEs

### Go Package Changes
- Deprecated `pkg/util.RandStringRunes` → use `k8s.io/apimachinery/pkg/util/rand.String`
- Deprecated `pkg/controller/test.RandStringBytes` → use `k8s.io/apimachinery/pkg/util/rand.String`
- Deprecated `sets.String` → use generic `sets.Set`

---

## Pre-Upgrade Checklist

### Required Actions

- [ ] **Verify Kubernetes version compatibility** (1.24-1.31 required for cert-manager 1.14)
- [ ] **Check current cert-manager version** (`kubectl -n cert-manager get deploy cert-manager -o yaml | grep image:`)
- [ ] **Identify target version** (1.14.4 or later, recommend 1.14.6)
- [ ] **Air-gapped environments only:** Pre-pull new startupapicheck image
- [ ] **Run GoNoGo bundle check** to detect deprecated annotations or configurations
- [ ] **Review custom CA integrations** for CSR criticality compatibility
- [ ] **Backup current configuration:**
  ```bash
  kubectl get certificates,certificaterequests,issuers,clusterissuers -A -o yaml > cert-manager-backup.yaml
  helm get values cert-manager -n cert-manager > cert-manager-values.yaml
  ```

### Recommended Actions

- [ ] **Test in non-production environment first**
- [ ] **Review existing ingress annotations** for deprecated certmanager.k8s.io/* usage
- [ ] **Check cluster-autoscaler behavior** if using ACME HTTP01 challenges
- [ ] **Update monitoring/alerts** for new startupapicheck job
- [ ] **Review certificate renewal windows** to minimize impact
- [ ] **Prepare rollback plan** (keep 1.13.x Helm chart available)

---

## Using the GoNoGo Bundle

### Installation

The cert-manager bundle is located at: `pkg/bundle/bundles/cert-manager.yaml`

### Running the Check

```bash
# Using the GoNoGo CLI with the cert-manager bundle
gonogo check -b pkg/bundle/bundles/cert-manager.yaml

# Or scan all bundles in the directory
gonogo check -d pkg/bundle/bundles/
```

### What the Bundle Checks

The GoNoGo bundle for cert-manager 1.13→1.14 includes:

1. **Kubernetes Version Compatibility**
   - Validates cluster version is between 1.24 and 1.31

2. **Required API Versions**
   - `apps/v1`
   - `v1`
   - `apiextensions.k8s.io/v1`
   - `admissionregistration.k8s.io/v1`

3. **OPA Policy Checks**
   - Detects deprecated `certmanager.k8s.io/*` annotations on Ingress resources
   - Detects deprecated `certmanager.k8s.io/*` annotations on Certificate resources
   - Identifies outdated startupapicheck jobs using old `ctl` image

4. **Resource Scanning**
   - Ingresses (`networking.k8s.io/v1/ingresses`)
   - Secrets (`v1/secrets`)
   - Certificates (`cert-manager.io/v1/certificates`)
   - CertificateRequests (`cert-manager.io/v1/certificaterequests`)
   - Issuers (`cert-manager.io/v1/issuers`)
   - ClusterIssuers (`cert-manager.io/v1/clusterissuers`)

### Interpreting Results

The bundle will output JSON with action items for any issues found:

```json
{
  "Addons": [
    {
      "Name": "cert-manager",
      "Versions": {
        "Current": "v1.13.0",
        "Upgrade": "1.14.6"
      },
      "UpgradeConfidence": 0.7,
      "ActionItems": [
        {
          "ResourceNamespace": "default",
          "ResourceKind": "Ingress",
          "ResourceName": "example-ingress",
          "Title": "Deprecated cert-manager annotation found on Ingress",
          "Description": "Ingress default/example-ingress uses deprecated cert-manager annotation: certmanager.k8s.io/issuer",
          "Remediation": "Update to use current cert-manager.io/* annotations...",
          "Severity": "0.3",
          "Category": "Reliability"
        }
      ],
      "Warnings": [
        "CRITICAL: Do NOT install v1.14.0, v1.14.1, v1.14.2, or v1.14.3..."
      ]
    }
  ]
}
```

---

## Upgrade Procedure

### Using Helm (Recommended)

```bash
# Add/update the cert-manager Helm repository
helm repo add jetstack https://charts.jetstack.io
helm repo update

# Check current version
helm list -n cert-manager

# Upgrade to 1.14.6 (or latest stable)
helm upgrade cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.14.6 \
  --reuse-values

# Verify the upgrade
kubectl -n cert-manager get pods
kubectl -n cert-manager get deploy cert-manager -o jsonpath='{.spec.template.spec.containers[0].image}'
```

### Using kubectl (Static Manifests)

```bash
# Apply the CRDs (if needed)
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.6/cert-manager.crds.yaml

# Apply the cert-manager manifests
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.6/cert-manager.yaml
```

### Post-Upgrade Verification

```bash
# Check all pods are running
kubectl -n cert-manager get pods

# Verify startupapicheck job completed successfully
kubectl -n cert-manager get jobs

# Test certificate issuance
kubectl get certificates -A
kubectl describe certificate <cert-name> -n <namespace>

# Check webhook is working
kubectl -n cert-manager logs deploy/cert-manager-webhook

# Verify cainjector is functioning
kubectl -n cert-manager logs deploy/cert-manager-cainjector
```

---

## Rollback Plan

If issues occur, rollback to cert-manager 1.13.x:

```bash
# Using Helm
helm rollback cert-manager -n cert-manager

# Or downgrade to specific version
helm upgrade cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.13.6 \
  --reuse-values
```

**Note:** Rolling back may require re-applying 1.13.x CRDs if they were modified during upgrade.

---

## Risk Assessment

| Risk Area | Level | Mitigation |
|-----------|-------|-----------|
| Version selection (wrong patch) | HIGH | **Use v1.14.4+, avoid v1.14.0-v1.14.3** |
| Air-gapped startupapicheck image | HIGH | Pre-pull new image before upgrade |
| CSR criticality with custom CAs | MEDIUM | Test in staging with actual CA first |
| Downtime during upgrade | LOW | Helm upgrade typically zero-downtime |
| Certificate renewal failures | LOW | Certificates renew automatically after upgrade |
| Cluster autoscaler eviction | LOW | Override annotation if needed |
| Deprecated annotations | LOW | Detected by GoNoGo OPA checks |

---

## Additional Resources

- [Official cert-manager 1.13→1.14 Upgrade Guide](https://cert-manager.io/docs/releases/upgrading/upgrading-1.13-1.14/)
- [cert-manager 1.14 Release Notes](https://cert-manager.io/docs/releases/release-notes/release-notes-1.14/)
- [cert-manager v1.14.6 GitHub Release](https://github.com/cert-manager/cert-manager/releases/tag/v1.14.6)
- [Supported Releases Matrix](https://cert-manager.io/docs/installation/supported-releases/)
- [GoNoGo Documentation](https://gonogo.docs.fairwinds.com)

---

## Decision Summary

### ✅ GO - Proceed with Upgrade

**Conditions:**
- Kubernetes version is 1.24 or higher
- Target version is 1.14.4 or later (recommend 1.14.6)
- Air-gapped environments have pre-pulled startupapicheck image
- Pre-upgrade checks completed successfully
- Tested in non-production environment
- No critical deprecated annotations found (or remediated)

### ⚠️ GO WITH CAUTION

**Additional requirements:**
- Custom CA integration requires staging validation
- Air-gapped environment needs careful image management
- Active cluster autoscaling during ACME challenges requires configuration review

### 🛑 NO GO - Do Not Upgrade Yet

**Blocking conditions:**
- Kubernetes version below 1.24 (upgrade K8s first)
- Planning to use v1.14.0, v1.14.1, v1.14.2, or v1.14.3 (use v1.14.4+)
- Critical production period without maintenance window
- Unable to test in non-production environment
- Air-gapped environment without ability to pull new images

---

## Conclusion

The cert-manager upgrade from 1.13.0 to 1.14.x is **recommended and safe** when following the guidelines above. The primary concerns are:

1. **Version selection** - Skip v1.14.0 through v1.14.3
2. **Image management** - Ensure new startupapicheck image is available
3. **Testing** - Validate in non-production first, especially with custom CAs

With proper planning and the correct target version (1.14.4+), this upgrade should proceed smoothly with minimal risk. The 1.14 series brings important security updates and new features while maintaining broad compatibility with modern Kubernetes versions.

**Final Recommendation: GO ✅** (target v1.14.6, follow checklist above)
