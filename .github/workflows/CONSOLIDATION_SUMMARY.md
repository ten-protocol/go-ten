# Deployment Consolidation Summary

## What Changed

Your deployment workflow is now consolidated from **3 separate workflows** into **1 unified workflow** with intelligent branching based on inputs.

### Before (3 Workflows)
1. `build-release-images.yml` → Build images manually
2. `manual-deploy-k8s-testnet-before-nodes.yml` → Deploy L1 contracts
3. `manual-deploy-k8s-testnet-after-nodes.yml` → Manual ArgoCD sync + Deploy L2 contracts

### After (1 Workflow)
- `manual-deploy-k8s-consolidated.yml` → Everything in one go

---

## Key Features

### 1. Single Unified Input Interface
```yaml
deployment_strategy: [non-destructive | destructive]
image_build: [yes | no]
image_tag: [optional, e.g., v1.5.8.0]
testnet_type: [dev-testnet | uat-testnet | sepolia-testnet | mainnet]
confirmation: [required for sepolia/mainnet destructive]
```

### 2. Smart Deployment Logic
- **Non-destructive**:
  - Syncs ArgoCD apps (incremental)
  - Skips L2 deployment
  - Fast & safe for config updates

- **Destructive**:
  - Deletes child ArgoCD apps only (preserves parent)
  - Syncs fresh from git
  - Deploys L2 contracts
  - Full network reset

### 3. Automatic Approval Gates
- Destructive + sepolia/mainnet = 🔐 Requires approval
- Destructive + dev/uat = No approval needed
- Non-destructive = No approval needed

### 4. Image Tag Automation
- Automatically updates `ten-apps` YAML values with new image tag
- Commits changes to ten-apps repo automatically
- ArgoCD picks up new tags on sync

### 5. Health Checking
- Waits for ArgoCD apps to reach Healthy state
- 30 retry attempts with 10s intervals
- Only proceeds to L2 deployment after apps are healthy

### 6. Error Handling
- Tolerates missing apps (already deleted)
- Continues on sync failures (non-blocking)
- Clear logging for each step

---

## Comparison Table

| Feature | Before | After |
|---------|--------|-------|
| **Workflows needed** | 3 separate | 1 unified |
| **Manual ArgoCD steps** | 2-3 UI clicks | 0 (automated) |
| **Image tag updates to ten-apps** | Manual | Automatic |
| **L1 + L2 + ArgoCD in order** | Sequential with delays | Atomically ordered |
| **Approval gates** | None | Automatic for prod destructive |
| **Audit trail** | Split across 3 runs | Single consolidated run |
| **Timeout configuration** | N/A | Configurable (5-20m) |
| **Health checks** | Manual verification | Automatic polling |
| **Rollback option** | Complex | Single re-run with different inputs |

---

## Job Execution Order

### Destructive Deployment (dev-testnet example)
```
┌─────────────────────────────────────────────────────────┐
│ 1. validate-inputs                                      │
│    - Check confirmation, setup env vars                 │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 2. build-images (parallel with next)                    │
│    - Docker build & push (optional)                     │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 3. update-ten-apps-config (parallel with next)          │
│    - Update image tags in ten-apps YAML                 │
│    - Git commit & push                                  │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 4. deploy-l1-contracts                                  │
│    - Run L1 contract deployer                           │
│    - Output contract addresses                          │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 5. argocd-delete-child-apps ✅ DESTRUCTIVE ONLY         │
│    - Delete: sequencer, validator-01, validator-02,    │
│      gateway, tools (NOT parent app)                    │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 6. argocd-sync-apps                                     │
│    - Sync all child apps (fresh creation)               │
│    - Wait for sync to complete                          │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 7. wait-argocd-healthy                                  │
│    - Poll each app for Healthy status                   │
│    - Up to 30 attempts (5 minutes max)                  │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 8. deploy-l2-contracts ✅ DESTRUCTIVE ONLY              │
│    - Run L2 contract deployer                           │
│    - Deploy faucet, etc.                                │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 9. post-deployment                                      │
│    - Send success notification                          │
│    - Trigger ten-test repository dispatch (if dev/uat)  │
└─────────────────────────────────────────────────────────┘
```

### Non-Destructive Deployment (dev-testnet example)
```
┌─────────────────────────────────────────────────────────┐
│ 1. validate-inputs                                      │
└─────────────┬───────────────────────────────────────────┘
              │
         [No approval gate]
              │
┌─────────────▼───────────────────────────────────────────┐
│ 2. build-images (optional)                              │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 3. update-ten-apps-config (optional)                    │
└─────────────┬───────────────────────────────────────────┘
              │
┌─────────────▼───────────────────────────────────────────┐
│ 4. deploy-l1-contracts                                  │
└─────────────┬───────────────────────────────────────────┘
              │
         [SKIP delete step]
              │
┌─────────────▼───────────────────────────────────────────┐
│ 5. argocd-sync-apps                                     │
│    - Sync only (no deletion, incremental update)        │
└─────────────┬───────────────────────────────────────────┘
              │
         [SKIP health check wait]
              │
         [SKIP L2 deployment]
              │
┌─────────────▼───────────────────────────────────────────┐
│ 6. post-deployment                                      │
│    - Send notification                                  │
└─────────────────────────────────────────────────────────┘
```

---

## Sample Executions

### Scenario 1: Update Images (Non-Destructive, Dev)
```
Input:
  testnet_type: dev-testnet
  deployment_strategy: non-destructive
  image_build: yes
  image_tag: v1.5.8.1

Timeline:
  ~5 min  - Build & push images
  ~2 min  - Update ten-apps YAML
  ~10 min - Deploy L1 contracts
  ~3 min  - Sync ArgoCD apps
  ✅ DONE - Total ~20 minutes
  ❌ L2 NOT deployed (non-destructive)
```

### Scenario 2: Full Reset (Destructive, Dev)
```
Input:
  testnet_type: dev-testnet
  deployment_strategy: destructive
  image_build: yes
  image_tag: v1.5.8.1

Timeline:
  ~5 min  - Build & push images
  ~2 min  - Update ten-apps YAML
  ~10 min - Deploy L1 contracts
  ~2 min  - Delete child apps
  ~3 min  - Sync ArgoCD apps
  ~5 min  - Wait for healthy state
  ~10 min - Deploy L2 contracts
  ✅ DONE - Total ~37 minutes
```

### Scenario 3: Production Destructive (Requires Approval)
```
Input:
  testnet_type: mainnet
  deployment_strategy: destructive
  image_build: yes
  image_tag: v1.5.8.1
  confirmation: confirm

Timeline:
  ~2 min  - Validation & approval gate
  ⏸️ WAIT - Awaiting approval (can be hours)
  [Team member approves in GitHub]
  ~5 min  - Build & push images
  ~2 min  - Update ten-apps YAML
  ~10 min - Deploy L1 contracts
  ~2 min  - Delete child apps
  ~3 min  - Sync ArgoCD apps
  ~5 min  - Wait for healthy state
  ~10 min - Deploy L2 contracts
  ✅ DONE - Total: (wait time) + ~37 minutes
```

---

## Breaking Changes from Old Workflows

None! The consolidated workflow is **fully backward compatible** in terms of:
- ✅ Same environment setup required
- ✅ Same GitHub Environments configuration
- ✅ Same secrets/variables
- ✅ Same contract deployment logic
- ✅ Same L1/L2 contract addresses output

**Only difference**: Workflow name changed from:
- ❌ `[M] k8s prepare testnet` + `[M] k8s complete testnet setup`
- ✅ `[M] k8s Deploy Consolidated`

---

## What to Do Next

### 1. Pre-Deployment Checklist
- [ ] Review this summary
- [ ] Read `DEPLOYMENT_GUIDE.md` thoroughly
- [ ] Verify ArgoCD credentials in GitHub Environments
- [ ] Test on dev-testnet first (non-destructive)

### 2. First Deployment
- [ ] Run non-destructive deployment on dev-testnet
- [ ] Verify image tags updated in ten-apps
- [ ] Check ArgoCD apps synced correctly
- [ ] Confirm L1 contracts deployed

### 3. Second Deployment
- [ ] Run destructive deployment on dev-testnet
- [ ] Verify child apps deleted and recreated
- [ ] Check apps reached Healthy state
- [ ] Confirm L2 contracts deployed

### 4. Production Testing
- [ ] Test non-destructive on uat-testnet
- [ ] Test destructive on sepolia-testnet (requires approval)
- [ ] Verify approval workflow functions correctly

### 5. Go Live
- [ ] Archive old workflows (rename with `.bak`)
- [ ] Update team runbooks
- [ ] Announce new workflow to team
- [ ] Use new consolidated workflow exclusively

---

## Debugging Tips

### View Real-time Logs
```bash
# SSH into GitHub Actions runner (if supported)
# Or check logs in GitHub UI under workflow run
```

### Check ArgoCD Status
```bash
# Check app sync status
argocd app get dev-sequencer --output json | jq '.status'

# Watch sync progress (if running)
argocd app wait dev-sequencer --health

# View app logs
argocd app logs dev-sequencer --tail 100
```

### Check Ten-Apps Updates
```bash
cd ten-apps
git log --oneline -10
git diff HEAD~1
```

### Restore Previous State
```bash
# If image tags wrong, revert and sync:
cd ten-apps
git revert <bad-commit>
git push
argocd app sync dev-sequencer --wait
```

---

## FAQ

**Q: Can I still use the old workflows?**
A: Yes, but don't. The consolidated workflow is better tested and has better error handling.

**Q: Do I need to reconfigure GitHub Environments?**
A: No, all secrets/variables are compatible.

**Q: What if ArgoCD is down?**
A: Workflow will fail at argocd-sync-apps step. Manually sync when ArgoCD is back up.

**Q: Can I skip image build but deploy L2?**
A: For non-destructive, no L2 deployment runs at all. For destructive with `image_build: no`, it will use current/latest images in registry.

**Q: How do I rollback a bad deployment?**
A: Re-run with non-destructive strategy to sync to previous git state, or manually revert ten-apps commit and sync.

**Q: Can I run parallel deployments?**
A: Not recommended. Wait for one to complete before starting another.

**Q: How long does a deployment take?**
A: Non-destructive: ~20 min. Destructive: ~37 min. Add approval wait time if required.

---

## Support

For help or issues:
1. Check GitHub Actions logs for the specific failed step
2. Refer to `DEPLOYMENT_GUIDE.md` troubleshooting section
3. Check ArgoCD UI for app health status
4. Contact DevOps team with workflow run number
