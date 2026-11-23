# Gloo Gateway GoNoGo Bundle

This bundle helps assess upgrade confidence when upgrading Gloo Gateway from version 1.20.8 to 1.21.x.

## Overview

The gloo-gateway bundle checks for potential issues when upgrading from Gloo Gateway 1.20.8 to 1.21.2. This version includes significant changes due to the Envoy upgrade from 1.35.x to 1.36.x.

## Usage

### Using the Default Bundle

Simply run gonogo without arguments to use all curated bundles including gloo-gateway:

```bash
gonogo check
```

### Using Only the Gloo Gateway Bundle

To check only Gloo Gateway upgrades:

```bash
gonogo check -b pkg/bundle/bundles/gloo-gateway.yaml
```

### Example Output

```json
{
  "Addons": [
    {
      "Name": "gloo-gateway",
      "Versions": {
        "Current": "1.20.8",
        "Upgrade": "1.21.2"
      },
      "UpgradeConfidence": 0.7,
      "ActionItems": [
        {
          "ResourceNamespace": "gloo-system",
          "ResourceKind": "Gateway",
          "ResourceName": "gateway-proxy",
          "Title": "HTTP/2 Max Concurrent Streams Default Changed",
          "Description": "Envoy 1.36.x changes the default max concurrent streams from 2147483647 to 1024. This may impact high-traffic scenarios.",
          "Remediation": "Review your traffic patterns. If you need the old behavior, set the runtime guard envoy.reloadable_features.safe_http2_options to false, but this is temporary. Consider tuning HTTP/2 settings explicitly.",
          "Severity": "0.3",
          "Category": "Performance"
        }
      ],
      "Warnings": [
        "Envoy version upgraded from 1.35.x to 1.36.x. This includes breaking changes to ExtProc, HTTP/2 defaults, and HTTP/1 CONNECT behavior."
      ]
    }
  ]
}
```

## Key Changes Checked

### 1. Envoy Version Upgrade (1.35.x → 1.36.x)

**HTTP/2 Default Value Changes:**
- Max concurrent streams: 2,147,483,647 → 1,024
- Initial stream window size: 256 MiB → 16 MiB
- Initial connection window size: 256 MiB → 24 MiB

**Impact:** High-traffic scenarios may be affected
**Mitigation:** Temporarily revert with runtime guard `envoy.reloadable_features.safe_http2_options=false` or tune HTTP/2 settings explicitly

### 2. HTTP/1.1 CONNECT Request Changes

**Change:** CONNECT requests now include RFC 9110 compliant Host header by default

**Impact:** Upstream proxies must handle CONNECT requests with Host headers
**Mitigation:** If needed, set runtime flag `envoy.reloadable_features.http_11_proxy_connect_legacy_format=true`

### 3. ExtProc Configuration Changes

**Change:** Removed support for `fail_open` and `FULL_DUPLEX_STREAMED` configuration combinations

**Impact:** Existing ExtProc configurations may fail
**Mitigation:** Update ExtProc configurations to use supported combinations

### 4. XSLT Transformation Deprecation (Enterprise)

**Change:** XSLT transformations deprecated in v1.21.0, will be removed in v1.22.0

**Impact:** XSLT transformations will stop working in v1.22.0
**Mitigation:** Migrate to external processing server

### 5. Incremental Upgrade Requirement

**Requirement:** If upgrading from v1.19.x or older, you must upgrade incrementally:
1. First upgrade to v1.20.x using the v1.20.x documentation
2. Then upgrade to v1.21.x using the v1.21.x documentation

## Kubernetes Compatibility

- **Minimum:** 1.29
- **Maximum:** 1.34

## Required API Versions

The bundle checks for the presence of:
- `apps/v1`
- `v1`
- `networking.k8s.io/v1`
- `apiextensions.k8s.io/v1`

## Resources Monitored

The OPA checks run against:
- `gateway.networking.k8s.io/v1/gateways`
- `gateway.networking.k8s.io/v1/httproutes`
- `v1/services`

## Additional Resources

- [Gloo Gateway 1.21 Upgrade Guide](https://docs.solo.io/gateway/1.21.x/operations/upgrade/)
- [Gloo Gateway 1.21 Breaking Changes](https://docs.solo.io/gloo-edge/main/operations/upgrading/faq/)
- [Envoy v1.36 Changelog](https://www.envoyproxy.io/docs/envoy/latest/version_history/v1.36/v1.36)
- [Solo.io Version Support Matrix](https://docs.solo.io/gateway/1.21.x/reference/versions/)

## Contributing

To update this bundle with additional checks:

1. Edit `pkg/bundle/bundles/gloo-gateway.yaml`
2. Add new OPA checks following the existing pattern
3. Test the bundle: `gonogo check -b pkg/bundle/bundles/gloo-gateway.yaml`
4. Update this documentation if needed
