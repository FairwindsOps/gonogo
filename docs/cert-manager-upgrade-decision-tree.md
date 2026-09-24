# cert-manager 1.13.0 → 1.14.x Upgrade Decision Tree

## Decision Flow

```
┌─────────────────────────────────────────────────┐
│   Starting Point: cert-manager 1.13.0          │
│   Question: Should we upgrade to 1.14.x?       │
└────────────────┬────────────────────────────────┘
                 │
                 ▼
         ┌───────────────┐
         │ Is K8s version│
         │ >= 1.24?      │
         └───┬───────┬───┘
             │       │
         NO  │       │  YES
             │       │
             ▼       ▼
    ┌────────────┐  ┌────────────────────────────┐
    │ STOP       │  │ Which version are you      │
    │ Upgrade K8s│  │ planning to target?        │
    │ to 1.24+   │  └────────┬───────────────────┘
    │ first      │           │
    └────────────┘           ▼
                    ┌─────────────────────┐
                    │ Is it v1.14.0,      │
                    │ v1.14.1, v1.14.2,   │
                    │ or v1.14.3?         │
                    └────┬───────────┬────┘
                         │           │
                     YES │           │ NO
                         │           │
                         ▼           ▼
               ┌──────────────┐  ┌───────────────────┐
               │ STOP         │  │ Is it v1.14.4+?   │
               │ Change target│  │ (recommend v1.14.6)│
               │ to v1.14.4+  │  └────┬──────────────┘
               │ Known bugs!  │       │
               └──────────────┘       │ YES
                                      ▼
                            ┌──────────────────────┐
                            │ Air-gapped/          │
                            │ restricted network?  │
                            └────┬─────────────┬───┘
                                 │             │
                             YES │             │ NO
                                 │             │
                                 ▼             ▼
                    ┌─────────────────────┐  ┌────────────────┐
                    │ Have you pre-pulled │  │ Run GoNoGo     │
                    │ startupapicheck     │  │ bundle check   │
                    │ image?              │  └────┬───────────┘
                    └───┬─────────────┬───┘       │
                        │             │           │
                    NO  │             │  YES      │
                        │             │           │
                        ▼             └───────────┤
           ┌──────────────────┐                  │
           │ STOP             │                  │
           │ Pre-pull image:  │                  │
           │ quay.io/jetstack/│                  │
           │ cert-manager-    │                  │
           │ startupapicheck  │                  │
           └──────────────────┘                  │
                                                 ▼
                                    ┌───────────────────────┐
                                    │ Any CRITICAL action   │
                                    │ items from GoNoGo?    │
                                    └────┬──────────────┬───┘
                                         │              │
                                     YES │              │ NO
                                         │              │
                                         ▼              ▼
                            ┌──────────────────┐   ┌───────────────┐
                            │ PAUSE            │   │ Tested in     │
                            │ Remediate issues │   │ staging/dev?  │
                            │ found by GoNoGo  │   └───┬───────┬───┘
                            │ Re-run check     │       │       │
                            └──────────────────┘   NO  │       │  YES
                                                       │       │
                                                       ▼       ▼
                                          ┌──────────────┐  ┌──────────┐
                                          │ RECOMMEND    │  │ ✅ GO   │
                                          │ Test first   │  │ Proceed  │
                                          │ (best        │  │ with     │
                                          │  practice)   │  │ upgrade  │
                                          └──────────────┘  └──────────┘
```

## Version Selection Matrix

| Your K8s Version | Current cert-manager | Can Upgrade to 1.14? | Recommended Target | Notes |
|------------------|---------------------|----------------------|-------------------|-------|
| 1.21-1.23 | 1.13.0 | ❌ NO | Upgrade K8s first | K8s 1.24 min required |
| 1.24-1.27 | 1.13.0 | ✅ YES | v1.14.6 | Both versions supported |
| 1.28-1.31 | 1.13.0 | ✅ YES | v1.14.6 | Only 1.14+ supports these |
| 1.32+ | 1.13.0 | ❌ NO | v1.15+ or v1.16+ | Need newer cert-manager |

## Risk Assessment by Scenario

### Scenario A: Online Environment, K8s 1.24-1.31
```
Risk Level: 🟢 LOW
Upgrade Path: Direct to v1.14.6
Blockers: None
Time to Upgrade: 15-30 minutes
```

### Scenario B: Air-gapped Environment, K8s 1.24-1.31
```
Risk Level: 🟡 MEDIUM
Upgrade Path: Direct to v1.14.6 (after image prep)
Blockers: startupapicheck image pull
Time to Upgrade: 1-2 hours (including image prep)
```

### Scenario C: Custom CA Integration
```
Risk Level: 🟡 MEDIUM
Upgrade Path: Staging → Validation → Production
Blockers: CSR criticality testing needed
Time to Upgrade: 1-3 days (including testing)
```

### Scenario D: K8s < 1.24
```
Risk Level: 🔴 HIGH
Upgrade Path: Upgrade K8s → cert-manager 1.14.6
Blockers: Kubernetes version
Time to Upgrade: Depends on K8s upgrade timeline
```

## Action Priority Matrix

| Priority | Action | When | Who |
|----------|--------|------|-----|
| 🔴 P0 | Verify K8s version ≥ 1.24 | Before planning | Platform team |
| 🔴 P0 | Select target version (v1.14.4+) | Before planning | Platform team |
| 🟡 P1 | Run GoNoGo bundle check | Before upgrade | DevOps/SRE |
| 🟡 P1 | Backup configurations | Before upgrade | DevOps/SRE |
| 🟡 P1 | Pre-pull images (air-gapped only) | Before upgrade | Platform team |
| 🟢 P2 | Test in staging | Before production | QA/DevOps |
| 🟢 P2 | Review deprecated annotations | Before upgrade | App teams |
| 🟢 P3 | Update monitoring/alerts | After upgrade | SRE |
| 🟢 P3 | Document changes | After upgrade | Platform team |

## Common Pitfalls to Avoid

### ❌ DON'T
1. ❌ Use v1.14.0, v1.14.1, v1.14.2, or v1.14.3
2. ❌ Upgrade directly to 1.14 from versions below 1.12 (upgrade to 1.12 first)
3. ❌ Skip staging/dev testing with custom CAs
4. ❌ Forget to pre-pull startupapicheck image in air-gapped environments
5. ❌ Ignore GoNoGo action items about deprecated annotations

### ✅ DO
1. ✅ Target v1.14.6 or later stable patch release
2. ✅ Run GoNoGo bundle check before upgrade
3. ✅ Backup current configuration
4. ✅ Test certificate issuance after upgrade
5. ✅ Have rollback plan ready
6. ✅ Monitor webhook and cainjector logs post-upgrade

## Emergency Rollback Decision

```
┌─────────────────────────────────────┐
│  Upgrade completed to 1.14.x        │
└────────────┬────────────────────────┘
             │
             ▼
    ┌────────────────────┐
    │ Are certificates   │
    │ being issued?      │
    └────┬───────────┬───┘
         │           │
     NO  │           │  YES
         │           │
         ▼           ▼
┌──────────────┐  ┌────────────────────┐
│ Check logs   │  │ Monitor for 1-2    │
│ for errors   │  │ hours              │
└────┬─────────┘  └────┬───────────────┘
     │                 │
     ▼                 ▼
┌─────────────────┐  ┌───────────────────┐
│ Critical errors?│  │ All working?      │
└────┬────────┬───┘  └────┬──────────┬───┘
     │        │           │          │
 YES │        │ NO    YES │          │ NO
     │        │           │          │
     ▼        ▼           ▼          ▼
┌─────────┐ ┌──────┐  ┌─────────┐ ┌──────────┐
│ROLLBACK │ │Debug │  │SUCCESS! │ │Continue  │
│Now!     │ │and   │  │Document │ │monitoring│
│         │ │fix   │  │upgrade  │ │          │
└─────────┘ └──────┘  └─────────┘ └──────────┘
```

## Quick Command Cheat Sheet

```bash
# Pre-Upgrade Checks
kubectl version --short                              # Check K8s version
kubectl -n cert-manager get deploy cert-manager \
  -o jsonpath='{.spec.template.spec.containers[0].image}'  # Current version
gonogo check -b pkg/bundle/bundles/cert-manager.yaml # Run GoNoGo check

# Backup
helm get values cert-manager -n cert-manager > backup-values.yaml
kubectl get certificates,issuers,clusterissuers -A -o yaml > backup-certs.yaml

# Upgrade
helm repo update
helm upgrade cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.14.6 \
  --reuse-values

# Verify
kubectl -n cert-manager get pods
kubectl get certificates -A
kubectl -n cert-manager logs deploy/cert-manager-webhook --tail=50

# Rollback (if needed)
helm rollback cert-manager -n cert-manager
```

## Support and Resources

- **GoNoGo Bundle:** `pkg/bundle/bundles/cert-manager.yaml`
- **Full Assessment:** `CERT_MANAGER_1.13_TO_1.14_GONOGO.md`
- **Quick Summary:** `SUMMARY.md`
- **Official Docs:** https://cert-manager.io/docs/releases/upgrading/upgrading-1.13-1.14/
- **Community:** cert-manager Slack: https://cert-manager.io/docs/contributing/
- **Issues:** https://github.com/cert-manager/cert-manager/issues

---

*Last Updated: May 1, 2026*  
*GoNoGo Project: https://github.com/FairwindsOps/gonogo*
