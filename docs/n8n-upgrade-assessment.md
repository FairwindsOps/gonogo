# n8n Helm Chart Upgrade Assessment: 1.16.9 to 1.16.37

## Overview

This document provides a comprehensive go/no-go assessment for upgrading the n8n Helm chart from version 1.16.9 to 1.16.37. This upgrade represents a significant change as it crosses the major version boundary from n8n application version 1.x to 2.x.

## Chart Version Details

- **Starting Version**: 1.16.9 (Released: December 9, 2025)
  - n8n Application Version: 1.123.4
  - Repository: community-charts/helm-charts

- **Target Version**: 1.16.37 (Released: April 22, 2026)
  - n8n Application Version: 2.18.x
  - Repository: community-charts/helm-charts

## Critical Breaking Changes

### 1. n8n Version 2.0 Breaking Changes

This upgrade crosses the n8n v2.0 boundary, which introduces significant breaking changes focused on security hardening and platform stability.

#### Security Changes

- **Task Runners Enabled by Default**
  - Code node executions now run in isolated environments
  - Improved security posture but may affect Code nodes relying on certain modules or environment access
  - **Action Required**: Review all Code nodes for compatibility with task runner isolation

- **Environment Variable Access Blocked in Code Nodes**
  - `N8N_BLOCK_ENV_ACCESS_IN_NODE` now defaults to `true`
  - Code nodes can no longer access environment variables by default
  - **Action Required**: Set `N8N_BLOCK_ENV_ACCESS_IN_NODE=false` if required, or refactor workflows to pass variables explicitly

- **Disabled Nodes by Default**
  - ExecuteCommand and LocalFileTrigger nodes are disabled by default
  - **Action Required**: Set `N8N_NODES_INCLUDE` to enable these nodes if needed

- **OAuth Callback Authentication**
  - OAuth callback URLs now require authentication by default
  - **Action Required**: Review your OAuth configurations for compatibility

#### Data Storage Changes

- **MySQL/MariaDB Support Dropped**
  - n8n v2.0 no longer supports MySQL or MariaDB databases
  - **Action Required**: Migrate to PostgreSQL before upgrading
  - **Severity**: CRITICAL - This is a blocking issue

- **In-Memory Binary Data Mode Removed**
  - Binary data must now use filesystem or S3-compatible storage
  - **Action Required**: Configure `N8N_BINARY_DATA_MODE` to either 'filesystem' or 's3'
  - **Severity**: HIGH - Workflows will fail without proper configuration

#### Workflow and Node Changes

- **Start Node Removed**
  - The Start node is no longer supported
  - **Action Required**: Replace Start nodes with specific trigger nodes (Manual Trigger, Webhook, Schedule Trigger, etc.)
  - **Severity**: HIGH - Workflows using Start nodes will fail

- **Python Code Node Migration**
  - Pyodide-based Python Code node has been removed
  - **Action Required**: Update Python Code nodes to use the `pythonNative` parameter with the new task runner
  - **Severity**: MEDIUM - Affects workflows using Python Code nodes

- **Workflow Publishing System**
  - Activate/Deactivate toggles replaced with Publish/Unpublish buttons
  - Provides better control over workflow deployment
  - **Action Required**: Review deployment workflows

#### Configuration Changes

- **CLI Changes**
  - `n8n --tunnel` option has been removed
  - **Action Required**: Use alternative webhook configuration methods

- **Dotenv Parsing Changes**
  - Upgraded dotenv library changes how `.env` files are parsed
  - Multiline values now supported
  - `#` marks the beginning of comments
  - **Action Required**: Review your `.env` files for compatibility

- **Release Channels Renamed**
  - `latest` → `stable`
  - `next` → `beta`
  - **Action Required**: Update image tags if using channel-based tags

### 2. Kubernetes Compatibility

- **Supported Kubernetes Versions**: 1.23 to 1.35
- **Required API Versions**:
  - `apps/v1`
  - `v1`
  - `batch/v1`

### 3. Required Resources

The following Kubernetes resources are used by n8n:
- Deployments
- StatefulSets
- Services
- PersistentVolumeClaims

## Pre-Upgrade Checklist

### Critical (Must Do Before Upgrade)

- [ ] **Database Migration**: If using MySQL/MariaDB, migrate to PostgreSQL
- [ ] **Binary Data Storage**: Configure filesystem or S3 storage mode
- [ ] **Run Migration Report Tool**: Use the official n8n migration tool at https://docs.n8n.io/migration-tool/
- [ ] **Backup Data**: Backup all workflow data, credentials, and configurations
- [ ] **Review Start Nodes**: Identify and plan replacement for any Start nodes in workflows

### Important (Strongly Recommended)

- [ ] **Test Environment**: Deploy v2.0 in a test environment first
- [ ] **Review Code Nodes**: Audit all Code nodes for:
  - Environment variable access
  - Module dependencies
  - Python code compatibility
- [ ] **Review OAuth Configurations**: Test OAuth flows with new authentication requirements
- [ ] **Update .env Files**: Review for multiline values and comment syntax
- [ ] **Document Custom Nodes**: Ensure ExecuteCommand and LocalFileTrigger usage is documented

### Recommended

- [ ] **Review Task Runner Configuration**: Understand task runner isolation impact
- [ ] **Update CI/CD**: Update deployment scripts for new Publish/Unpublish workflow
- [ ] **Review Monitoring**: Ensure monitoring is compatible with new architecture
- [ ] **Update Documentation**: Document the upgrade process for your team

## Upgrade Process

### Step 1: Preparation

1. Review all warnings and breaking changes in this document
2. Run the n8n Migration Report tool to identify specific issues in your instance
3. Create a comprehensive backup of your n8n data

### Step 2: Pre-Upgrade Configuration

1. **If using MySQL/MariaDB**: Migrate to PostgreSQL
2. **Configure binary data storage**:
   ```yaml
   env:
     N8N_BINARY_DATA_MODE: filesystem  # or 's3'
     # If using S3:
     N8N_BINARY_DATA_STORAGE_S3_HOST: s3.amazonaws.com
     N8N_BINARY_DATA_STORAGE_S3_BUCKET_NAME: your-bucket
     N8N_BINARY_DATA_STORAGE_S3_BUCKET_REGION: us-east-1
   ```
3. **Configure task runners** (if needed):
   ```yaml
   env:
     N8N_RUNNERS_MODE: external  # or adjust as needed
   ```
4. **Enable specific nodes** (if needed):
   ```yaml
   env:
     N8N_NODES_INCLUDE: '["n8n-nodes-base.executeCommand", "n8n-nodes-base.localFileTrigger"]'
   ```

### Step 3: Test Environment Upgrade

1. Deploy n8n v2.0 in a test environment
2. Import a copy of your production workflows
3. Test critical workflows thoroughly
4. Verify OAuth integrations work correctly
5. Test Code nodes with task runner isolation

### Step 4: Production Upgrade

1. Schedule maintenance window
2. Create final backup
3. Upgrade Helm chart: `helm upgrade n8n community-charts/n8n --version 1.16.37`
4. Monitor logs for errors
5. Test critical workflows
6. Verify all integrations are working

### Step 5: Post-Upgrade Tasks

1. Replace any Start nodes in workflows
2. Update Python Code nodes to use pythonNative
3. Review and adjust Code nodes for task runner compatibility
4. Update team documentation
5. Train team on new Publish/Unpublish workflow

## Rollback Plan

If issues arise during upgrade:

1. **Immediate Rollback**: 
   ```bash
   helm rollback n8n
   ```
2. **Restore Database Backup**: If database migration occurred
3. **Verify Data Integrity**: Check that all workflows and data are restored
4. **Document Issues**: Record what went wrong for troubleshooting

## OPA Policy Checks

The n8n bundle includes Open Policy Agent (OPA) checks that will automatically validate:

1. **MySQL/MariaDB Detection**: Warns if MySQL is detected in configuration
2. **Start Node Detection**: Identifies workflows using the deprecated Start node
3. **Code Node Environment Variable Access**: Checks for proper configuration
4. **Binary Data Mode**: Validates binary data storage configuration
5. **Task Runner Configuration**: Ensures task runners are properly configured

## Resources and References

- [n8n v2.0 Breaking Changes Documentation](https://docs.n8n.io/2-0-breaking-changes/)
- [n8n Migration Report Tool](https://docs.n8n.io/migration-tool/)
- [n8n Hosting Configuration](https://docs.n8n.io/hosting/configuration/)
- [Task Runner Documentation](https://docs.n8n.io/hosting/securing/hardening-task-runners/)
- [Binary Data Storage Configuration](https://docs.n8n.io/hosting/configuration/environment-variables/binary-data/)
- [Database Configuration](https://docs.n8n.io/hosting/configuration/database/)

## Risk Assessment

### High Risk Items
- Database migration from MySQL/MariaDB to PostgreSQL
- Workflows using Start nodes
- Binary data storage configuration
- Code nodes with environment variable access

### Medium Risk Items
- Python Code nodes requiring migration
- OAuth integrations
- Custom node configurations
- Task runner isolation impacts

### Low Risk Items
- Workflow publishing UI changes
- Release channel naming
- CLI option changes

## Recommendation

**GO** with caution: This upgrade can proceed but requires careful planning and execution due to significant breaking changes. The upgrade is **CRITICAL** only if you:

1. Successfully migrate from MySQL/MariaDB to PostgreSQL (if applicable)
2. Configure binary data storage properly
3. Update all workflows using Start nodes
4. Run the migration report tool and address all critical issues

**NO-GO** conditions:
- Currently using MySQL/MariaDB and cannot migrate to PostgreSQL
- Unable to allocate sufficient testing time
- Critical workflows rely heavily on Start nodes without replacement plan
- No backup and rollback plan in place

## Summary

This upgrade represents a significant platform evolution with important security and reliability improvements. While it introduces breaking changes, they are manageable with proper planning and testing. The benefits include:

- Enhanced security through task runner isolation
- More secure default configurations
- Better control over workflow deployment
- Improved platform stability and performance

Allocate sufficient time for testing and migration activities to ensure a smooth upgrade process.
