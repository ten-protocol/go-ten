# Non-Destructive Deployment Changes

## Summary
Updated the consolidated k8s deployment workflow to properly handle non-destructive deployments by deleting only the host deployments and enclave statefulsets, then syncing ArgoCD apps without deleting the apps themselves or redeploying L1/L2 contracts.

## Key Changes

### 1. L1 Contracts Deployment (Destructive Only)
- **Changed**: L1 contract deployment now only runs for destructive deployments
- **Job**: `deploy-l1-contracts`
- **Condition**: Added `needs.validate-inputs.outputs.IS_DESTRUCTIVE == 'true'`

### 2. New Job: ArgoCD Delete Resources (Non-Destructive Only)
- **Added**: New job `argocd-delete-resources` that runs only for non-destructive deployments
- **Purpose**: Delete host deployments and enclave statefulsets using ArgoCD CLI
- **Resources deleted per app (sequencer, validator-01, validator-02)**:
  - `Deployment`: `{app}-host`
  - `StatefulSet`: `{app}-enclave`
  - `StatefulSet`: `{app}-enclave02` (if exists)
- **Method**: Uses `argocd app delete-resource` command (no kubectl/k8s access needed)
- **Dependencies**: Runs after `update-ten-apps-config` completes

### 3. Updated ArgoCD Sync Job
- **Changed**: `argocd-sync-apps` now depends on both deletion jobs
- **Dependencies**:
  - `argocd-delete-child-apps` (destructive path)
  - `argocd-delete-resources` (non-destructive path)
- **Condition**: Uses `always()` with result checks to run after either deletion path completes

### 4. L2 Contracts Deployment (Already Destructive Only)
- **Unchanged**: Already only runs for destructive deployments
- No changes needed

## Workflow Flow

### Non-Destructive Flow:
1. ✅ Build images (optional)
2. ✅ Update image tags in ten-apps
3. ⏭️ Skip L1 contract deployment
4. ⏭️ Skip ArgoCD app deletion
5. ✅ **Delete host/enclave resources via ArgoCD** (new)
6. ✅ Sync ArgoCD apps
7. ✅ Wait for apps to be healthy
8. ⏭️ Skip L2 contract deployment

### Destructive Flow (Unchanged):
1. ✅ Build images (optional)
2. ✅ Update image tags and node selectors in ten-apps
3. ✅ Deploy L1 contracts
4. ✅ Delete ArgoCD child apps
5. ⏭️ Skip resource deletion
6. ✅ Sync ArgoCD apps (recreates from scratch)
7. ✅ Wait for apps to be healthy
8. ✅ Deploy L2 contracts

## Technical Details

### ArgoCD Resource Deletion
Uses ArgoCD CLI command to delete resources without needing kubectl access to the private cluster:

```bash
argocd app delete-resource "$app" \
  --server "$ARGOCD_SERVER" \
  --kind Deployment \
  --resource-name "${app}-host" \
  --all
```

This approach:
- Works with private clusters (no direct k8s access needed)
- Uses existing ArgoCD authentication via `--server` flag and `ARGOCD_AUTH_TOKEN` env var
- Properly manages resources tracked by ArgoCD
- Allows ArgoCD to sync and recreate resources with new image tags

**Note**: The `--server` flag is required on all ArgoCD CLI commands to specify the ArgoCD server address.

## Benefits

1. **No ArgoCD App Deletion**: Apps remain configured, faster deployment
2. **No L1/L2 Redeployment**: Saves time and prevents contract address changes
3. **Clean Resource Update**: Deletes old pods/containers before syncing new images
4. **Private Cluster Compatible**: Uses ArgoCD API, no kubectl access required
5. **Maintains State**: Secrets, ConfigMaps, and other resources preserved

## Testing Recommendations

1. Test non-destructive deployment on dev-testnet first
2. Verify resources are deleted before sync
3. Confirm new images are pulled and deployed
4. Check that existing secrets/configs are maintained
5. Validate no L1/L2 contract redeployment occurs
