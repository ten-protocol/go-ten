# Consolidated K8s Deployment Guide

## Overview

The `manual-deploy-k8s-consolidated.yml` workflow consolidates three separate operations into a single unified GitHub Action:

1. **Image Build & Push** (merged from `build-release-images.yml`)
2. **L1 Contract Deployment**
3. **ArgoCD App Management & Sync**
4. **L2 Contract Deployment**

This eliminates manual steps and provides a single audit trail for deployments.

---

## Deployment Strategies

### Non-Destructive Deployment
- ✅ Builds & pushes images (if enabled)
- ✅ Deploys L1 contracts
- ✅ **Syncs ArgoCD apps** (incremental update)
- ❌ **SKIPS L2 deployment** (no reset needed)
- ❌ Does NOT delete child apps
- **Use case**: Image updates, config changes, minor updates

### Destructive Deployment
- ✅ Builds & pushes images (if enabled)
- ✅ Deploys L1 contracts
- ✅ **Deletes all child ArgoCD apps** (not parent)
- ✅ **Syncs ArgoCD apps** (fresh creation)
- ✅ **Waits for apps to be healthy**
- ✅ **Deploys L2 contracts** (full reset)
- ❌ Requires approval for sepolia/mainnet
- **Use case**: Full network reset, major version upgrades

---

## Input Parameters

| Parameter | Type | Required | Options | Description |
|-----------|------|----------|---------|-------------|
| `testnet_type` | choice | ✅ | dev-testnet, uat-testnet, sepolia-testnet, mainnet | Target environment |
| `deployment_strategy` | choice | ✅ | non-destructive, destructive | Deployment approach |
| `image_build` | choice | ✅ | yes, no | Build & push images |
| `image_tag` | string | ❌ | v1.5.8.0 or empty | Docker image tag (empty = latest) |
| `confirmation` | string | ❌ | "confirm" | Required for sepolia/mainnet |
| `log_level` | number | ❌ | 1-5 | Application log level |
| `max_gas_gwei` | string | ❌ | 1.5 | Max gas price for transactions |
| `sync_timeout` | choice | ❌ | 5m, 10m, 15m, 20m | ArgoCD sync timeout |
| `sequencer_node_selector` | string | ❌ | aks-sgxpool01-xxx | Target node for sequencer (empty = no change) |
| `validator_01_node_selector` | string | ❌ | aks-sgxpool01-xxx | Target node for validator-01 (empty = no change) |
| `validator_02_node_selector` | string | ❌ | aks-sgxpool02-xxx | Target node for validator-02 (empty = no change) |
| `gateway_node_selector` | string | ❌ | aks-pool-xxx | Target node for gateway (empty = no change) |

---

## Usage Examples

### Example 1: Non-Destructive Image Update (Dev)
```
testnet_type: dev-testnet
deployment_strategy: non-destructive
image_build: yes
image_tag: v1.5.8.0
```
**What happens:**
- Builds images with tag `v1.5.8.0`
- Updates ten-apps YAML values with new tag
- Deploys L1 contracts
- Syncs ArgoCD child apps with new image tags
- ✅ No L2 deployment (non-destructive)

### Example 2: Destructive Full Reset (Dev)
```
testnet_type: dev-testnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.8.0
```
**What happens:**
- Builds images with tag `v1.5.8.0`
- Updates ten-apps YAML values with new tag
- Deploys L1 contracts
- **Deletes all child apps:** dev-sequencer, dev-validator-01, dev-validator-02, dev-gateway, dev-tools
- Syncs ArgoCD child apps (fresh creation with new images)
- Waits for apps to reach Healthy state
- Deploys L2 contracts
- Triggers ten-test repository dispatch

### Example 3: Destructive Deployment to Mainnet (Requires Approval)
```
testnet_type: mainnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.8.0
confirmation: confirm
```
**What happens:**
- 🔐 **APPROVAL GATE**: Workflow pauses, requires human approval
- Once approved, proceeds with full destructive deployment
- All mainnet apps deleted and recreated
- L2 contracts deployed

### Example 4: Non-Destructive Sync Only (Config Update)
```
testnet_type: uat-testnet
deployment_strategy: non-destructive
image_build: no
confirmation: ""
```
**What happens:**
- ❌ No image build
- ✅ Deploys fresh L1 contracts
- ✅ Syncs ArgoCD apps (picks up config from ten-apps)
- ✅ **No L2 deployment**

---

## Approval Gates

Automatic approval required for:
- ✅ Destructive deployments to **sepolia-testnet**
- ✅ Destructive deployments to **mainnet**

**Requirements:**
1. Type `"confirm"` in the `confirmation` field
2. Approval must be granted by another team member in GitHub Actions UI

**Non-destructive deployments:**
- No approval required for any environment

**Destructive to dev/uat:**
- No approval required

---

## Environment Configuration

### ArgoCD Setup

Each environment requires:

1. **ArgoCD Server Secret** (per environment):
   - `ARGOCD_SERVER_PROD`: Production ArgoCD instance URL
   - `ARGOCD_SERVER_NONPROD`: Non-production ArgoCD instance URL

2. **ArgoCD Token Secret**:
   - `ARGOCD_TOKEN`: Service account token with app management permissions

3. **Create ArgoCD Service Account**:
```bash
# In each ArgoCD instance
kubectl create serviceaccount github-actions -n argocd
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

# Get token
kubectl -n argocd create token github-actions
```

### GitHub Environment Setup

Set these in GitHub Environments (`Settings > Environments`):

**Environment: dev-testnet**
- Secrets: `ARGOCD_SERVER_NONPROD`, `ARGOCD_TOKEN`, etc.

**Environment: uat-testnet**
- Secrets: `ARGOCD_SERVER_NONPROD`, `ARGOCD_TOKEN`, etc.

**Environment: sepolia-testnet**
- Secrets: `ARGOCD_SERVER_NONPROD`, `ARGOCD_TOKEN`, etc.
- Requires approval (destructive only)

**Environment: mainnet**
- Secrets: `ARGOCD_SERVER_PROD`, `ARGOCD_TOKEN`, etc.
- Requires approval (destructive only)

---

## Child Apps Per Environment

### Dev Testnet
- `dev-sequencer`
- `dev-validator-01`
- `dev-validator-02`
- `dev-gateway`
- `dev-tools`

### UAT Testnet
- `uat-sequencer`
- `uat-validator-01`
- `uat-validator-02`
- `uat-gateway`
- `uat-tools`

### Sepolia Testnet
- `sepolia-sequencer`
- `sepolia-validator-01`
- `sepolia-validator-02`
- `sepolia-gateway`
- `sepolia-gateway-dexynth`
- `sepolia-gateway-pentest`
- `sepolia-tools`

### Mainnet
- `mainnet-sequencer`
- `mainnet-validator-01`
- `mainnet-validator-02`
- `mainnet-gateway`
- `mainnet-postgres-client`
- `mainnet-tools`

---

## Job Flow Diagram

### Non-Destructive Flow
```
validate-inputs
    ↓
build-images (optional)
    ↓
update-ten-apps-config (optional)
    ↓
deploy-l1-contracts
    ↓
argocd-delete-child-apps ❌ (SKIPPED)
    ↓
argocd-sync-apps ✅ (incremental)
    ↓
wait-argocd-healthy ✅
    ↓
deploy-l2-contracts ❌ (SKIPPED)
    ↓
post-deployment
```

### Destructive Flow (dev/uat)
```
validate-inputs
    ↓
build-images (optional)
    ↓
update-ten-apps-config (optional)
    ↓
deploy-l1-contracts
    ↓
argocd-delete-child-apps ✅ (delete all)
    ↓
argocd-sync-apps ✅ (fresh sync)
    ↓
wait-argocd-healthy ✅ (30 attempts, 10s each)
    ↓
deploy-l2-contracts ✅
    ↓
post-deployment
```

### Destructive Flow (sepolia/mainnet)
```
validate-inputs
    ↓
approval ⏳ (REQUIRES HUMAN APPROVAL)
    ↓
... (rest of destructive flow)
```

---

## Monitoring & Debugging

### Check Workflow Status
1. Go to GitHub Actions tab
2. Find "k8s Deploy Consolidated" workflow
3. Check individual job logs

### View ArgoCD Sync Status
```bash
argocd app get dev-sequencer
argocd app get dev-sequencer --refresh
argocd app logs dev-sequencer
```

### Common Issues

**Issue: "App not found or already deleted"**
- During destructive: This is normal if app doesn't exist
- Workflow continues safely

**Issue: "Failed to sync app (timeout)"**
- Increase `sync_timeout` parameter
- Check app health manually: `argocd app get <app-name>`
- Check K8s pod status: `kubectl get pods -n <env>`

**Issue: "Approval gate not appearing"**
- Confirm `confirmation: "confirm"` is set
- Confirm environment is sepolia/mainnet
- Confirm strategy is destructive

**Issue: "L2 deployment not running"**
- Only runs for destructive deployments
- Check `deployment_strategy` is "destructive"
- Check ArgoCD apps reached Healthy state

---

## Migration from Old Workflows

### Step 1: Archive Old Workflows
```bash
# Rename (don't delete yet, keep as reference)
mv manual-deploy-k8s-testnet-before-nodes.yml manual-deploy-k8s-testnet-before-nodes.yml.bak
mv manual-deploy-k8s-testnet-after-nodes.yml manual-deploy-k8s-testnet-after-nodes.yml.bak
```

### Step 2: Test New Workflow
- Run on dev-testnet with non-destructive first
- Compare results with old workflow
- Verify ArgoCD apps sync correctly

### Step 3: Cutover
- Update runbooks and documentation
- Inform team of new workflow
- Use new consolidated workflow going forward

### Step 4: Decommission Old Workflows
- After 2-3 weeks of successful runs
- Delete archived `.bak` files
- Update CI/CD documentation

---

## Rollback Procedure

If deployment fails:

### Option 1: Re-run with Non-Destructive
```
deployment_strategy: non-destructive
image_build: no
```
This will sync apps to latest git state without rebuilding.

### Option 2: Manual Rollback
```bash
# Revert ten-apps commit
cd ten-apps
git revert <commit-hash>
git push

# Sync ArgoCD to pick up previous state
argocd app sync dev-sequencer --wait
argocd app sync dev-validator-01 --wait
# ... repeat for all apps
```

### Option 3: Full Cluster Reset (Dev/UAT only)
```bash
# Re-run with destructive strategy
# This will delete all apps and recreate from scratch
```

---

## Troubleshooting Checklist

- [ ] `confirmation: "confirm"` entered for sepolia/mainnet?
- [ ] Correct `testnet_type` selected?
- [ ] ArgoCD credentials (`ARGOCD_TOKEN`, `ARGOCD_SERVER`) configured?
- [ ] Image tag valid (if specified)?
- [ ] L1 contract deployment succeeded?
- [ ] ArgoCD child apps reached Healthy state (check logs)?
- [ ] Ten-apps repository accessible via `DEPLOY_ACTIONS_PAT`?
- [ ] Sufficient timeout for sync (`sync_timeout`)?

---

## Support & Questions

For issues or questions:
1. Check job logs in GitHub Actions
2. Review this guide's Troubleshooting section
3. Check ArgoCD UI for app status: `https://<argocd-server>/applications`
4. Contact DevOps team with workflow run number and error details
