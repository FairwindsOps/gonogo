# cert-manager GoNoGo Bundle

This bundle assesses the upgrade readiness for cert-manager from version 1.13.0 to 1.14.x.

## Quick Start

```bash
# Run the bundle check against your cluster
gonogo check -b pkg/bundle/bundles/cert-manager.yaml

# Or run all bundles including cert-manager
gonogo check -d pkg/bundle/bundles/
```

## What This Bundle Checks

### 1. Version Compatibility
- **From:** cert-manager 1.13.0 (EOL: June 5, 2024)
- **To:** cert-manager 1.14.x (recommend v1.14.6+)
- **Kubernetes:** 1.24 - 1.31
- **OpenShift:** 4.11 - 4.16

### 2. API Version Requirements
Validates that your cluster has the necessary APIs:
- `apps/v1`
- `v1`
- `apiextensions.k8s.io/v1`
- `admissionregistration.k8s.io/v1`

### 3. OPA Policy Checks

#### Check 1: Deprecated Ingress Annotations
Scans all Ingress resources for deprecated `certmanager.k8s.io/*` annotations:
- `certmanager.k8s.io/issuer`
- `certmanager.k8s.io/cluster-issuer`
- `certmanager.k8s.io/acme-challenge-type`
- `certmanager.k8s.io/acme-dns01-provider`

**Severity:** 0.3 (Low-Medium)  
**Remediation:** Update to `cert-manager.io/*` namespace

#### Check 2: Deprecated Certificate Annotations
Scans all Certificate resources for deprecated annotations:
- `certmanager.k8s.io/issuer`
- `certmanager.k8s.io/cluster-issuer`

**Severity:** 0.3 (Low-Medium)  
**Remediation:** Update to `cert-manager.io/*` namespace

#### Check 3: Outdated Startupapicheck Image
Identifies startupapicheck Jobs using the old `ctl` image instead of the new `startupapicheck` image.

**Severity:** 0.5 (Medium)  
**Remediation:** cert-manager 1.14+ requires the new startupapicheck OCI image

### 4. Resource Scanning
The bundle scans these resource types in your cluster:
- Ingresses (`networking.k8s.io/v1/ingresses`)
- Secrets (`v1/secrets`)
- Certificates (`cert-manager.io/v1/certificates`)
- CertificateRequests (`cert-manager.io/v1/certificaterequests`)
- Issuers (`cert-manager.io/v1/issuers`)
- ClusterIssuers (`cert-manager.io/v1/clusterissuers`)

## Critical Warnings

The bundle includes important warnings about known issues:

### 🔴 Version-Specific Bugs
**DO NOT** install these versions:
- v1.14.0 - Wrong cainjector image, SAN issues
- v1.14.1 - CA/SelfSigned SAN critical flag bug
- v1.14.2 - Additional bugs
- v1.14.3 - Additional bugs

**Target:** v1.14.4 or later (recommend v1.14.6)

### ⚠️ Breaking Changes

1. **Startupapicheck Image**
   - New image: `quay.io/jetstack/cert-manager-startupapicheck`
   - **Critical for air-gapped environments**

2. **CSR Encoding**
   - KeyUsage and BasicConstraints now marked as critical
   - May affect external CA validation

3. **ACME Solver Annotation**
   - New default: `cluster-autoscaler.kubernetes.io/safe-to-evict: "true"`
   - Override in podTemplate if needed

## Example Output

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
          "ResourceNamespace": "production",
          "ResourceKind": "Ingress",
          "ResourceName": "app-ingress",
          "Title": "Deprecated cert-manager annotation found on Ingress",
          "Description": "Ingress production/app-ingress uses deprecated cert-manager annotation: certmanager.k8s.io/issuer",
          "Remediation": "Update to use current cert-manager.io/* annotations...",
          "Severity": "0.3",
          "Category": "Reliability"
        }
      ],
      "Warnings": [
        "CRITICAL: Do NOT install v1.14.0, v1.14.1, v1.14.2, or v1.14.3...",
        "The startupapicheck job now uses a new OCI image..."
      ]
    }
  ]
}
```

## Understanding Action Items

### Severity Levels
- **0.1 - 0.3:** Low severity - Should be addressed but not blocking
- **0.4 - 0.6:** Medium severity - Recommended to fix before upgrade
- **0.7 - 1.0:** High severity - Must be addressed before upgrade

### Categories
- **Reliability:** Issues affecting service stability
- **Security:** Security-related concerns
- **Compatibility:** API or version compatibility issues

## Remediation Examples

### Fix Deprecated Ingress Annotation

**Before:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example
  annotations:
    certmanager.k8s.io/issuer: "letsencrypt-prod"  # Deprecated
spec:
  tls:
  - hosts:
    - example.com
    secretName: example-tls
```

**After:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example
  annotations:
    cert-manager.io/issuer: "letsencrypt-prod"  # Current
spec:
  tls:
  - hosts:
    - example.com
    secretName: example-tls
```

### Override ACME Solver Eviction

If you need ACME HTTP01 solver pods to be non-evictable:

```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    solvers:
    - http01:
        ingress:
          podTemplate:
            metadata:
              annotations:
                cluster-autoscaler.kubernetes.io/safe-to-evict: "false"
```

## Pre-Upgrade Checklist

Use this checklist before running your upgrade:

- [ ] Run GoNoGo bundle check
- [ ] Review all action items
- [ ] Fix critical (≥0.7 severity) issues
- [ ] Address medium (0.4-0.6 severity) issues
- [ ] Document low (≤0.3 severity) issues for future work
- [ ] Verify Kubernetes version ≥ 1.24
- [ ] Backup current configuration
- [ ] Pre-pull startupapicheck image (air-gapped only)
- [ ] Test in non-production environment
- [ ] Plan maintenance window
- [ ] Prepare rollback procedure

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: cert-manager Upgrade Readiness
on:
  schedule:
    - cron: '0 0 * * 1'  # Weekly check
  workflow_dispatch:

jobs:
  gonogo-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup kubectl
        uses: azure/setup-kubectl@v3
      
      - name: Configure kubeconfig
        run: |
          echo "${{ secrets.KUBECONFIG }}" > kubeconfig.yaml
          export KUBECONFIG=kubeconfig.yaml
      
      - name: Run GoNoGo Check
        run: |
          curl -L https://github.com/FairwindsOps/gonogo/releases/latest/download/gonogo-linux-amd64.tar.gz | tar xz
          ./gonogo check -b pkg/bundle/bundles/cert-manager.yaml > results.json
      
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: gonogo-results
          path: results.json
```

## Related Documentation

- **Comprehensive Assessment:** [`CERT_MANAGER_1.13_TO_1.14_GONOGO.md`](../../../CERT_MANAGER_1.13_TO_1.14_GONOGO.md)
- **Quick Summary:** [`SUMMARY.md`](../../../SUMMARY.md)
- **Decision Tree:** [`docs/cert-manager-upgrade-decision-tree.md`](../../../docs/cert-manager-upgrade-decision-tree.md)
- **Official cert-manager Upgrade Guide:** https://cert-manager.io/docs/releases/upgrading/upgrading-1.13-1.14/

## Support

- **GoNoGo Issues:** https://github.com/FairwindsOps/gonogo/issues
- **cert-manager Issues:** https://github.com/cert-manager/cert-manager/issues
- **cert-manager Slack:** https://cert-manager.io/docs/contributing/

## Contributing

To improve this bundle:

1. Test the bundle against real clusters
2. Submit findings via GitHub issues
3. Propose additional OPA checks
4. Update version ranges as new releases come out
5. Add new resource types to scan

## License

This bundle is part of the GoNoGo project, licensed under Apache 2.0.

---

**Bundle Version:** 1.0  
**Created:** May 2026  
**Target Versions:** cert-manager 1.13.0 → 1.14.6  
**Kubernetes Compatibility:** 1.24 - 1.31
