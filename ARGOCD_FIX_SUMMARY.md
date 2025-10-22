# ArgoCD Deployment Failure - Root Cause Analysis

## Failed Run: 18713271026

### Failure Summary:
- ✅ Build images: SUCCESS
- ✅ Update ten-apps config: SUCCESS  
- ✅ Deploy L1 contracts: SUCCESS
- ❌ ArgoCD delete child apps: **FAILED** (permission denied on parent sync)
- ❌ Deploy L2 contracts: **FAILED** (couldn't connect - apps not running)

---

## Root Cause: Wrong Parent App Name

### The Problem:
```bash
# Line 715 in workflow:
PARENT_APP="${TESTNET_SHORT_NAME}-testnet"

# For UAT environment:
TESTNET_SHORT_NAME = "uat"
PARENT_APP = "uat-testnet"  ❌ WRONG - This app doesn't exist!
```

### The Fix:
```bash
PARENT_APP="ten-${TESTNET_SHORT_NAME}-apps"

# For UAT environment:
PARENT_APP = "ten-uat-apps"  ✅ CORRECT
```

### Correct App Structure (confirmed by user):
```
Environment: UAT
├─ Parent App:  ten-uat-apps
├─ Child Apps:
│   ├─ uat-sequencer
│   ├─ uat-validator-01
│   ├─ uat-validator-02
│   ├─ uat-gateway
│   └─ uat-tools

Environment: DEV
├─ Parent App:  ten-dev-apps
├─ Child Apps:
│   ├─ dev-sequencer
│   ├─ dev-validator-01
│   ├─ dev-validator-02
│   ├─ dev-gateway
│   └─ dev-tools

Environment: SEPOLIA
├─ Parent App:  ten-sepolia-apps
├─ Child Apps:
│   ├─ sepolia-sequencer
│   ├─ sepolia-validator-01
│   ├─ sepolia-validator-02
│   ├─ sepolia-gateway
│   └─ sepolia-tools

Environment: MAINNET
├─ Parent App:  ten-mainnet-apps
├─ Child Apps:
│   ├─ mainnet-sequencer
│   ├─ mainnet-validator-01
│   ├─ mainnet-validator-02
│   ├─ mainnet-gateway
│   └─ mainnet-tools
```

---

## Secondary Issue: Missing --grpc-web Flag

### Log Evidence:
```
{"level":"warning","msg":"Failed to invoke grpc call. Use flag --grpc-web in grpc calls..."}
{"level":"fatal","msg":"rpc error: code = PermissionDenied desc = permission denied"}
```

### Current Code (Line 772):
```bash
argocd app sync "$PARENT_APP" --server "${ARGOCD_SERVER}"
```

### Should Be:
```bash
argocd app sync "$PARENT_APP" --server "${ARGOCD_SERVER}" --grpc-web
```

---

## Changes Required

### 1. Fix Parent App Name (Line 715)
```bash
# BEFORE
PARENT_APP="${TESTNET_SHORT_NAME}-testnet"

# AFTER
PARENT_APP="ten-${TESTNET_SHORT_NAME}-apps"
```

### 2. Add --grpc-web Flag to Sync (Line 772)
```bash
# BEFORE
if argocd app sync "$PARENT_APP" --server "${ARGOCD_SERVER}"; then

# AFTER
if argocd app sync "$PARENT_APP" --server "${ARGOCD_SERVER}" --grpc-web; then
```

### 3. Optional: Add --grpc-web to Delete Commands
The delete commands work but show warnings. Adding --grpc-web would clean up logs:
```bash
# Line ~730 (example)
argocd app delete "$app" --server "${ARGOCD_SERVER}" --grpc-web --cascade --yes
```

---

## Impact Analysis

### What Happened:
1. Workflow tried to sync app "uat-testnet" ❌
2. This app doesn't exist
3. ArgoCD returned "permission denied" (really "not found")
4. Child apps were deleted but not recreated
5. L2 deployment failed because validator-01 wasn't running

### What Will Happen After Fix:
1. Workflow deletes child apps ✅
2. Workflow syncs "ten-uat-apps" parent ✅
3. Parent recreates child apps ✅
4. Child apps deploy with new images ✅
5. L2 contracts deploy successfully ✅

---

## Testing Recommendation

After applying fix:
1. Trigger UAT destructive deployment
2. Verify parent app name in logs: "ten-uat-apps"
3. Verify child apps are recreated
4. Verify L2 deployment succeeds

---

## Files to Update
- `.github/workflows/manual-deploy-k8s-consolidated.yml`
  - Line 715: Parent app name
  - Line 772: Add --grpc-web flag


---

## ADDITIONAL ISSUE: L2 Deployer Docker Image

### Discovery:
The L2 contract deployer **DOES use Docker** internally! The `go run` command actually starts a Docker container.

### Log Evidence:
```
DEPLOY_DOCKERIMAGE: testnetobscuronet.azurecr.io/obscuronet/uat_hardhatdeployer:latest
```

But we built:
```
testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:689b04c
```

### The Problem:
```yaml
# Line 1237 in workflow:
env:
  DEPLOY_DOCKERIMAGE: ${{ vars.DOCKER_BUILD_TAG_L2_HARDHAT_DEPLOYER }}
  # This resolves to: testnetobscuronet.azurecr.io/obscuronet/uat_hardhatdeployer:latest
```

### The Fix:
We need to use the IMAGE_TAG from build-images output:
```yaml
env:
  DEPLOY_DOCKERIMAGE: ${{ env.REGISTRY }}/${{ env.REGISTRY_ORG }}/hardhatdeployer:${{ needs.build-images.outputs.IMAGE_TAG }}
  # This will be: testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:689b04c
```

### Why L2 Failed (Two Reasons):
1. **Primary**: Couldn't connect to `uat-validator-01.ten.xyz:3000` because validator wasn't running (ArgoCD parent sync failed)
2. **Secondary**: Using old `:latest` image instead of new commit hash image

### Same Issue in L1 Deployer?
No! L1 deployer doesn't use DEPLOY_DOCKERIMAGE at all. It's purely `go run` with no Docker.

---

## Complete Fix List

### 1. Fix Parent App Name (Line 715)
```bash
# BEFORE
PARENT_APP="${TESTNET_SHORT_NAME}-testnet"

# AFTER  
PARENT_APP="ten-${TESTNET_SHORT_NAME}-apps"
```

### 2. Add --grpc-web to Parent Sync (Line 772)
```bash
# BEFORE
if argocd app sync "$PARENT_APP" --server "${ARGOCD_SERVER}"; then

# AFTER
if argocd app sync "$PARENT_APP" --server "${ARGOCD_SERVER}" --grpc-web; then
```

### 3. Fix L2 Deployer Image Tag (Line 1237)
```yaml
# BEFORE
env:
  DEPLOY_DOCKERIMAGE: ${{ vars.DOCKER_BUILD_TAG_L2_HARDHAT_DEPLOYER }}

# AFTER
env:
  DEPLOY_DOCKERIMAGE: ${{ env.REGISTRY }}/${{ env.REGISTRY_ORG }}/hardhatdeployer:${{ needs.build-images.outputs.IMAGE_TAG }}
```

### 4. Optional: Clean up unused DEPLOY_DOCKERIMAGE in L1
L1 deployer doesn't use it, so can be removed from line 601 (but harmless to keep).

---

## Why Docker Was Building/Pulling in L2

The L2 contract deployer workflow:
```
go run ./testnet/launcher/l2contractdeployer/cmd
  ↓
  Runs Go code that starts Docker container
  ↓
  docker run testnetobscuronet.azurecr.io/obscuronet/uat_hardhatdeployer:latest
  ↓
  Container runs Hardhat deployment scripts
  ↓
  Container tries to connect to uat-validator-01.ten.xyz
  ↓
  FAIL: Connection refused (validator not running due to ArgoCD fail)
```

So it wasn't "building" - it was **pulling** the image from registry to run the deployment.

