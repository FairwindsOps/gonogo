# cert-manager 1.13.0 → 1.14.x Upgrade Assessment - Document Index

This index helps you navigate all documentation created for the cert-manager upgrade assessment.

---

## 🎯 Start Here

**New to this assessment?** Read these documents in order:

1. 📄 **[SUMMARY.md](./SUMMARY.md)** (2 pages)
   - Quick decision matrix: GO/NO-GO
   - Critical requirements at a glance
   - Command cheat sheet
   - **Read this first for executive summary**

2. 📊 **[docs/cert-manager-upgrade-decision-tree.md](./docs/cert-manager-upgrade-decision-tree.md)** (5 pages)
   - Visual decision flow diagram
   - Version selection matrix
   - Risk assessment by scenario
   - Common pitfalls and best practices
   - **Read this for decision-making guidance**

3. 📘 **[CERT_MANAGER_1.13_TO_1.14_GONOGO.md](./CERT_MANAGER_1.13_TO_1.14_GONOGO.md)** (15 pages)
   - Comprehensive technical assessment
   - Breaking changes with detailed impact analysis
   - Pre-upgrade checklist
   - Upgrade and rollback procedures
   - **Read this for complete technical details**

---

## 🔧 Implementation Resources

### GoNoGo Bundle

**Location:** [`pkg/bundle/bundles/cert-manager.yaml`](./pkg/bundle/bundles/cert-manager.yaml)

The executable bundle specification that checks your cluster for upgrade readiness.

**Quick Start:**
```bash
gonogo check -b pkg/bundle/bundles/cert-manager.yaml
```

**Documentation:** [`pkg/bundle/bundles/README_cert-manager.md`](./pkg/bundle/bundles/README_cert-manager.md)

---

## 📚 Complete File Listing

| File | Purpose | Pages | Read If... |
|------|---------|-------|-----------|
| **SUMMARY.md** | Quick reference | 2 | You need the bottom line |
| **CERT_MANAGER_1.13_TO_1.14_GONOGO.md** | Full assessment | 15 | You're planning the upgrade |
| **docs/cert-manager-upgrade-decision-tree.md** | Decision guide | 5 | You need help deciding |
| **pkg/bundle/bundles/cert-manager.yaml** | GoNoGo bundle | - | You're running the check |
| **pkg/bundle/bundles/README_cert-manager.md** | Bundle docs | 7 | You're using the bundle |
| **INDEX_cert-manager-upgrade.md** | This file | 1 | You need navigation |

**Total Documentation:** ~30 pages  
**Total Lines Added:** 1,159 lines  
**Bundle Checks:** 3 OPA policies + version/API checks

---

## 🎓 Document Use Cases

### Use Case 1: "Should we upgrade?"
**Read:** SUMMARY.md → Decision Tree  
**Time:** 15 minutes  
**Output:** GO/NO-GO decision

### Use Case 2: "How do we upgrade safely?"
**Read:** CERT_MANAGER_1.13_TO_1.14_GONOGO.md  
**Time:** 45 minutes  
**Output:** Detailed upgrade plan

### Use Case 3: "What issues exist in our cluster?"
**Run:** `gonogo check -b pkg/bundle/bundles/cert-manager.yaml`  
**Read:** README_cert-manager.md for interpretation  
**Time:** 5 minutes (run) + 15 minutes (interpret)  
**Output:** List of action items

### Use Case 4: "I found an issue, how do I fix it?"
**Read:** README_cert-manager.md → Remediation Examples  
**Time:** 10 minutes  
**Output:** Fix implementation

### Use Case 5: "Planning a production upgrade"
**Read:** All documents + run bundle check  
**Time:** 2-3 hours  
**Output:** Complete upgrade runbook

---

## 📋 Quick Reference Tables

### Version Compatibility
| cert-manager | Kubernetes | Status |
|--------------|-----------|--------|
| 1.13.0 | 1.21 - 1.27 | Current |
| 1.14.6 | 1.24 - 1.31 | Target |

### Critical Versions to AVOID
| Version | Issue | Status |
|---------|-------|--------|
| v1.14.0 | Wrong cainjector image | ❌ Skip |
| v1.14.1 | SAN critical flag bug | ❌ Skip |
| v1.14.2 | Multiple bugs | ❌ Skip |
| v1.14.3 | Multiple bugs | ❌ Skip |
| v1.14.4+ | All fixes applied | ✅ Use |

### Breaking Changes
| Change | Impact | Priority |
|--------|--------|----------|
| startupapicheck image | Air-gapped envs | 🔴 HIGH |
| CSR criticality | Custom CAs | 🟡 MEDIUM |
| ACME eviction | Autoscaling | 🟢 LOW |

---

## 🔍 Finding Specific Information

### "What are the breaking changes?"
- **Short answer:** SUMMARY.md → "Breaking Changes Summary"
- **Detailed:** CERT_MANAGER_1.13_TO_1.14_GONOGO.md → "Critical Breaking Changes"

### "Which version should I target?"
- **Quick:** SUMMARY.md → "Critical Requirements"
- **Matrix:** Decision Tree → "Version Selection Matrix"

### "What does the GoNoGo check do?"
- **Overview:** README_cert-manager.md → "What This Bundle Checks"
- **Technical:** cert-manager.yaml (bundle source)

### "How do I fix deprecated annotations?"
- **Examples:** README_cert-manager.md → "Remediation Examples"
- **Context:** CERT_MANAGER_1.13_TO_1.14_GONOGO.md → "Pre-Upgrade Checklist"

### "What's the risk level?"
- **Matrix:** Decision Tree → "Risk Assessment by Scenario"
- **Detailed:** CERT_MANAGER_1.13_TO_1.14_GONOGO.md → "Risk Assessment"

### "How do I rollback?"
- **Commands:** SUMMARY.md → "Installation Commands"
- **Detailed:** CERT_MANAGER_1.13_TO_1.14_GONOGO.md → "Rollback Plan"
- **Decision:** Decision Tree → "Emergency Rollback Decision"

---

## 🚀 Quick Commands

```bash
# Run the assessment
gonogo check -b pkg/bundle/bundles/cert-manager.yaml

# View results with jq
gonogo check -b pkg/bundle/bundles/cert-manager.yaml | jq .

# Check for critical issues only
gonogo check -b pkg/bundle/bundles/cert-manager.yaml | \
  jq '.Addons[].ActionItems[] | select(.Severity | tonumber >= 0.7)'

# Export results
gonogo check -b pkg/bundle/bundles/cert-manager.yaml > assessment-results.json
```

---

## 🔗 External Resources

### Official Documentation
- [cert-manager 1.13→1.14 Upgrade Guide](https://cert-manager.io/docs/releases/upgrading/upgrading-1.13-1.14/)
- [cert-manager 1.14 Release Notes](https://cert-manager.io/docs/releases/release-notes/release-notes-1.14/)
- [Supported Releases Matrix](https://cert-manager.io/docs/installation/supported-releases/)

### GitHub
- [cert-manager v1.14.6 Release](https://github.com/cert-manager/cert-manager/releases/tag/v1.14.6)
- [GoNoGo Project](https://github.com/FairwindsOps/gonogo)
- [This Assessment PR](https://github.com/FairwindsOps/gonogo/pull/201)

### Community
- [cert-manager Slack](https://cert-manager.io/docs/contributing/)
- [cert-manager Issues](https://github.com/cert-manager/cert-manager/issues)
- [GoNoGo Issues](https://github.com/FairwindsOps/gonogo/issues)

---

## 📊 Assessment Metadata

| Property | Value |
|----------|-------|
| **Assessment Date** | May 1, 2026 |
| **Source Version** | cert-manager 1.13.0 |
| **Target Version** | cert-manager 1.14.6 |
| **Kubernetes Range** | 1.24 - 1.31 |
| **Overall Risk** | LOW-MEDIUM |
| **Recommendation** | ✅ GO (target v1.14.4+) |
| **Documentation Pages** | ~30 pages |
| **OPA Checks** | 3 policies |
| **Resource Types Scanned** | 6 types |
| **API Versions Checked** | 4 APIs |

---

## 🎯 Success Criteria Checklist

Before marking the upgrade as successful, ensure:

- [ ] GoNoGo bundle check executed
- [ ] All critical (≥0.7) action items resolved
- [ ] Kubernetes version is 1.24+
- [ ] Target version is 1.14.4 or later
- [ ] Air-gapped image pre-pulled (if applicable)
- [ ] Configuration backed up
- [ ] Tested in non-production
- [ ] Upgrade executed
- [ ] All pods running
- [ ] Certificates issuing successfully
- [ ] No errors in logs
- [ ] Monitoring updated
- [ ] Documentation updated

---

## 📞 Support

**Questions about this assessment?**
- Open an issue: https://github.com/FairwindsOps/gonogo/issues
- Tag: `cert-manager`, `upgrade-assessment`

**Questions about cert-manager upgrade?**
- cert-manager docs: https://cert-manager.io/docs/
- cert-manager Slack: https://cert-manager.io/docs/contributing/

**Questions about GoNoGo?**
- GoNoGo docs: https://gonogo.docs.fairwinds.com
- GitHub: https://github.com/FairwindsOps/gonogo

---

## 📝 Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0 | May 1, 2026 | Initial assessment created |

---

## 🤝 Contributing

To improve this assessment:

1. **Test the bundle** against real clusters
2. **Report findings** via GitHub issues
3. **Suggest improvements** to OPA checks
4. **Update** as new cert-manager versions release
5. **Share** your upgrade experience

---

## 📄 License

This assessment is part of the GoNoGo project.  
License: Apache 2.0

---

**Assessment ID:** cert-manager-1.13-to-1.14  
**GoNoGo Version:** Compatible with v0.4.0+  
**Last Updated:** May 1, 2026

