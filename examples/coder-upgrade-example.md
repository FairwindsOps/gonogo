# Coder 2.24.1 to 2.25.0 Upgrade Check Example

This example demonstrates how to use GoNoGo to assess the upgrade confidence for Coder from version 2.24.1 to 2.25.0.

## Overview

Coder v2.25.0 (released August 5, 2025) is a mainline release with several breaking changes and new features. This bundle spec helps you validate your cluster and configuration before performing the upgrade.

## Key Changes in v2.25.0

### Breaking Changes

1. **Connection Logs Migration**: Connection logs have been moved from the Audit Log to a new Connection Log entity. This affects:
   - SSH connections
   - Workspace app connections
   - Browser port-forwarding connections
   
2. **Audit Log Cleanup**: Connection events older than 90 days will be deleted from the Audit Log during upgrade.

3. **CLI Preset Flag**: The `coder create` command now requires `--preset` flag for templates with presets but no default, breaking non-interactive workflows.

4. **Workspace Update Behavior**: Updates now explicitly trigger a stop build followed by a start build (from v2.24.0 onwards).

### New Features

- **Dynamic Parameters GA**: Generally available with immutable parameters, masking, and search capabilities
- **DevContainer Integration**: Automatic detection and auto-start support
- **OAuth2 Provider**: Experimental OAuth2 provider functionality
- **MCP Server**: External Coder MCP server for AI agent workspace creation
- **Connection Logs**: New dedicated connection logs page and API

## Prerequisites

- GoNoGo installed ([installation instructions](https://gonogo.docs.fairwinds.com/installation/))
- kubectl configured with access to your cluster
- Coder currently deployed via Helm (version 2.24.x)

## Usage

### Using the Built-in Bundle

If the Coder bundle is included in your GoNoGo installation:

```bash
gonogo check
```

This will automatically load all built-in bundles including the Coder bundle.

### Using the Bundle File Directly

```bash
gonogo check -b pkg/bundle/bundles/coder.yaml
```

### Specifying Multiple Bundles

```bash
gonogo check -b pkg/bundle/bundles/coder.yaml -b pkg/bundle/bundles/metrics-server.yaml
```

## Understanding the Results

GoNoGo will check the following aspects:

### 1. Kubernetes Version Compatibility

The bundle verifies your cluster is running Kubernetes 1.19 or newer (up to 1.36).

### 2. Helm Release Detection

GoNoGo will detect if you have a Helm release named `coder` in the version range 2.24.1 to 2.25.0.

### 3. OPA Policy Checks

Four OPA checks are included:

#### Connection Log Retention Configuration
- **Severity**: 0.3
- **Check**: Verifies if `CODER_CONNECTION_LOG_RETENTION` is configured
- **Remediation**: Set this environment variable in your Helm values to manage log growth

#### PostgreSQL Version Compatibility
- **Severity**: 0.2
- **Check**: Ensures PostgreSQL is version 11 or newer
- **Remediation**: Upgrade PostgreSQL before upgrading Coder if needed

#### Audit Log Backup Warning
- **Severity**: 0.4
- **Check**: Warns about connection events being deleted from Audit Log
- **Remediation**: Back up your database or query connection events via REST API before upgrading

#### Non-Interactive Workflow Compatibility
- **Severity**: 0.5
- **Check**: Identifies automated workflows that may break due to preset changes
- **Remediation**: Update scripts to explicitly pass `--preset` flag

### 4. Resource Validation

The bundle checks for the following Kubernetes resources:
- Deployments (apps/v1)
- Services (v1)
- ConfigMaps (v1)
- Secrets (v1)
- PersistentVolumeClaims (v1)
- Ingresses (networking.k8s.io/v1)
- Roles (rbac.authorization.k8s.io/v1)
- RoleBindings (rbac.authorization.k8s.io/v1)

## Pre-Upgrade Checklist

Before upgrading to Coder 2.25.0:

- [ ] Review all warnings in the GoNoGo output
- [ ] Back up PostgreSQL database containing audit logs if historical connection data is needed
- [ ] Verify PostgreSQL is version 11 or newer with contrib package
- [ ] Update automated scripts using `coder create` to include `--preset` flag
- [ ] Review templates with presets and set defaults where appropriate
- [ ] Plan for connection log retention policy
- [ ] Test the upgrade in a staging environment if available
- [ ] Review the full changelog: https://github.com/coder/coder/releases/tag/v2.25.0

## Upgrading Coder

After GoNoGo gives a positive confidence score:

```bash
# Update Helm repo
helm repo update coder-v2

# Upgrade to v2.25.0
helm upgrade coder coder-v2/coder \
  --namespace coder \
  --values values.yaml \
  --version 2.25.0
```

Or using OCI registry:

```bash
helm upgrade coder oci://ghcr.io/coder/chart/coder \
  --namespace coder \
  --values values.yaml \
  --version 2.25.0
```

## Post-Upgrade Verification

After upgrading:

1. Verify Coder pods are running:
   ```bash
   kubectl get pods -n coder
   ```

2. Check Coder version:
   ```bash
   kubectl exec -n coder deployment/coder -- /coder version
   ```

3. Access the new Connection Log page in the dashboard

4. Test workspace creation with presets

5. Verify dynamic parameters are working correctly

## Troubleshooting

### Coder Pods Not Starting

Check pod logs:
```bash
kubectl logs -n coder deployment/coder
```

Common issues:
- PostgreSQL connection issues
- Missing environment variables
- TLS certificate problems

### Connection Logs Not Appearing

- Verify you have Premium Coder license
- Check that users have appropriate permissions
- Review CODER_CONNECTION_LOG_RETENTION configuration

### Preset Prompts in Automated Workflows

Update scripts to include explicit preset selection:
```bash
coder create myworkspace --template mytemplate --preset default
```

## Additional Resources

- [Coder v2.25.0 Release Notes](https://github.com/coder/coder/releases/tag/v2.25.0)
- [Coder v2.25 Changelog](https://coder.com/changelog/coder-2-25)
- [Coder Documentation](https://coder.com/docs)
- [Connection Logs Documentation](https://coder.com/docs/admin/monitoring/connection-logs)
- [Dynamic Parameters Documentation](https://coder.com/docs/admin/templates/extending-templates/dynamic-parameters)
- [DevContainers Integration](https://coder.com/docs/admin/templates/managing-templates/devcontainers)

## License

This bundle spec follows the same Apache 2.0 license as the GoNoGo project.
