# Workflow Decision Tree - All Permutations & Combinations

## Input Parameters
- `testnet_type`: dev-testnet, uat-testnet, sepolia-testnet, mainnet
- `deployment_strategy`: non-destructive, destructive
- `image_build`: yes, no
- `use_azure_hsm`: yes, no (only matters if image_build=yes)
- `confirmation`: "" or "confirm"

---

## Decision Tree

```
START
│
├─ validate-inputs ✅ (always runs)
│   ├─ Determines: NEEDS_APPROVAL (sepolia/mainnet=true, dev/uat=false unless confirmation='confirm')
│   ├─ Determines: IS_DESTRUCTIVE (destructive=true, non-destructive=false)
│   └─ Determines: TESTNET_SHORT_NAME, CONFIG_PATH, CHILD_APPS
│
├─ image_build == 'yes'? 
│   ├─ YES ✅
│   │   ├─ build-images ✅
│   │   │   ├─ Builds: enclave, host, hardhatdeployer
│   │   │   ├─ Tags with: commit hash or git tag
│   │   │   ├─ use_azure_hsm == 'yes'?
│   │   │   │   ├─ YES: Build enclave with HSM (Azure TENANT_ID, SUBSCRIPTION_ID)
│   │   │   │   └─ NO: Build enclave without HSM
│   │   │   └─ Push to registry
│   │   │
│   │   ├─ update-ten-apps-config ✅
│   │   │   ├─ Updates image tags in ten-apps repo
│   │   │   ├─ IS_DESTRUCTIVE == 'true'? → Also update node selectors (if provided)
│   │   │   └─ Commits & pushes to ten-apps
│   │   │
│   │   └─ NEEDS_APPROVAL == 'true'?
│   │       ├─ YES ✅ → approval (Environment gate)
│   │       │   ├─ Shows detailed summary:
│   │       │   │   ├─ Environment, strategy, triggered by
│   │       │   │   ├─ Image tag, HSM status
│   │       │   │   ├─ Config update status
│   │       │   │   ├─ List of deployment actions
│   │       │   │   └─ Warnings
│   │       │   ├─ APPROVED ✅ → Continue
│   │       │   └─ REJECTED ❌ → revert-config-on-rejection
│   │       │       └─ Reverts ten-apps commit (images stay in registry)
│   │       └─ NO ⏭️ → Skip approval (dev/uat without confirmation)
│   │
│   └─ NO ⏭️
│       ├─ build-images ⏭️ (skipped)
│       ├─ update-ten-apps-config ⏭️ (skipped - depends on build-images)
│       └─ approval ⏭️ (skipped - depends on update-ten-apps-config)
│
└─ deployment_strategy == 'destructive'?
    │
    ├─ YES (DESTRUCTIVE) 🔥
    │   │
    │   ├─ deploy-l1-contracts ✅
    │   │   ├─ Depends on: [validate-inputs, approval, update-ten-apps-config]
    │   │   ├─ If image_build=no: ❌ BLOCKED (update-ten-apps-config skipped)
    │   │   ├─ Runs: go run ./testnet/launcher/l1contractdeployer/cmd
    │   │   └─ Outputs: Contract addresses, L1_START_HASH
    │   │
    │   ├─ argocd-delete-child-apps ✅
    │   │   ├─ Depends on: [validate-inputs, approval, deploy-l1-contracts]
    │   │   ├─ If deploy-l1-contracts blocked: ❌ BLOCKED
    │   │   ├─ Deletes: sequencer, validator-01, validator-02 (+ additional apps if specified)
    │   │   └─ Syncs parent app to recreate children
    │   │
    │   ├─ argocd-sync-apps ✅
    │   │   ├─ Depends on: [validate-inputs, argocd-delete-child-apps]
    │   │   ├─ Full app sync with all sync waves and hooks
    │   │   └─ Syncs: sequencer, validator-01, validator-02, gateway, tools
    │   │
    │   ├─ wait-argocd-healthy ✅
    │   │   ├─ Waits for all apps to be healthy
    │   │   └─ Timeout: sync_timeout (default 10m)
    │   │
    │   ├─ deploy-l2-contracts ✅
    │   │   ├─ Depends on: [validate-inputs, deploy-l1-contracts, wait-argocd-healthy]
    │   │   ├─ Parses L1 config from ten-apps
    │   │   ├─ Runs: go run ./testnet/launcher/l2contractdeployer/cmd
    │   │   └─ Grants sequencer permissions, sets challenge period
    │   │
    │   └─ post-deployment ✅
    │       └─ Success notification
    │
    └─ NO (NON-DESTRUCTIVE) 📦
        │
        ├─ argocd-delete-resources ✅
        │   ├─ Depends on: [validate-inputs, approval, argocd-sync-apps (must be skipped)]
        │   ├─ Works with image_build=no: ✅ YES
        │   ├─ Deletes per app:
        │   │   ├─ host deployment
        │   │   ├─ enclave statefulset
        │   │   └─ enclave02 statefulset (if exists)
        │   └─ Apps: sequencer, validator-01, validator-02
        │
        ├─ argocd-sync-apps ✅
        │   ├─ Depends on: [validate-inputs, argocd-delete-resources, argocd-delete-child-apps (must be skipped)]
        │   ├─ Works with image_build=no: ✅ YES
        │   ├─ For each app:
        │   │   ├─ Checks if OutOfSync
        │   │   ├─ If Synced: ⏭️ Skip
        │   │   ├─ If OutOfSync:
        │   │   │   ├─ Gets list of OutOfSync resources
        │   │   │   └─ Syncs ONLY those specific resources (not full app)
        │   │   └─ Uses --prune=false (no deletions)
        │   └─ Apps: sequencer, validator-01, validator-02, gateway, tools
        │
        ├─ wait-argocd-healthy ✅
        │   └─ Waits for apps to be healthy
        │
        └─ post-deployment ✅
            └─ Success notification
```

---

## All Permutations Matrix

### Legend
- ✅ Runs successfully
- ⏭️ Skipped (intentional)
- ❌ Blocked/Fails
- 🔒 Requires approval

---

### Permutation 1: DEV + NON-DESTRUCTIVE + IMAGE_BUILD=YES + HSM=NO
```
✅ validate-inputs
✅ build-images (no HSM)
✅ update-ten-apps-config (image tags only)
⏭️ approval (not required for dev)
✅ argocd-delete-resources
✅ argocd-sync-apps (only OutOfSync resources)
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS - Rolling update with new images, no downtime
```

### Permutation 2: DEV + NON-DESTRUCTIVE + IMAGE_BUILD=NO
```
✅ validate-inputs
⏭️ build-images
⏭️ update-ten-apps-config
⏭️ approval
✅ argocd-delete-resources
✅ argocd-sync-apps (uses existing images from ten-apps)
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS - Restart pods with existing images
```

### Permutation 3: DEV + DESTRUCTIVE + IMAGE_BUILD=YES + HSM=NO
```
✅ validate-inputs
✅ build-images (no HSM)
✅ update-ten-apps-config (image tags + node selectors if provided)
⏭️ approval (not required for dev)
✅ deploy-l1-contracts
✅ argocd-delete-child-apps
✅ argocd-sync-apps (full sync)
✅ wait-argocd-healthy
✅ deploy-l2-contracts
✅ post-deployment

Result: ✅ SUCCESS - Full redeployment with new images and contracts
```

### Permutation 4: DEV + DESTRUCTIVE + IMAGE_BUILD=NO
```
✅ validate-inputs
⏭️ build-images
⏭️ update-ten-apps-config
⏭️ approval
❌ deploy-l1-contracts (BLOCKED - update-ten-apps-config skipped)
❌ argocd-delete-child-apps (BLOCKED)
❌ argocd-sync-apps (BLOCKED)
❌ wait-argocd-healthy (BLOCKED)
❌ deploy-l2-contracts (BLOCKED)
❌ post-deployment

Result: ❌ FAILS - Destructive requires config update (image_build must be yes)
```

### Permutation 5: UAT + NON-DESTRUCTIVE + IMAGE_BUILD=YES + HSM=NO
```
✅ validate-inputs
✅ build-images (no HSM)
✅ update-ten-apps-config (image tags only)
⏭️ approval (not required for uat unless confirmation='confirm')
✅ argocd-delete-resources
✅ argocd-sync-apps (only OutOfSync resources)
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS - Rolling update with new images
```

### Permutation 6: UAT + NON-DESTRUCTIVE + IMAGE_BUILD=NO
```
✅ validate-inputs
⏭️ build-images
⏭️ update-ten-apps-config
⏭️ approval
✅ argocd-delete-resources
✅ argocd-sync-apps (uses existing images)
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS - Restart with existing images
```

### Permutation 7: UAT + DESTRUCTIVE + IMAGE_BUILD=YES + HSM=NO
```
✅ validate-inputs
✅ build-images (no HSM)
✅ update-ten-apps-config (image tags + node selectors)
⏭️ approval (not required unless confirmation='confirm')
✅ deploy-l1-contracts
✅ argocd-delete-child-apps
✅ argocd-sync-apps (full sync)
✅ wait-argocd-healthy
✅ deploy-l2-contracts
✅ post-deployment

Result: ✅ SUCCESS - Full redeployment
```

### Permutation 8: UAT + DESTRUCTIVE + IMAGE_BUILD=NO
```
✅ validate-inputs
⏭️ build-images
⏭️ update-ten-apps-config
⏭️ approval
❌ deploy-l1-contracts (BLOCKED)
❌ argocd-delete-child-apps (BLOCKED)
❌ argocd-sync-apps (BLOCKED)
❌ wait-argocd-healthy (BLOCKED)
❌ deploy-l2-contracts (BLOCKED)
❌ post-deployment

Result: ❌ FAILS - Destructive requires image_build=yes
```

### Permutation 9: UAT + NON-DESTRUCTIVE + IMAGE_BUILD=YES + CONFIRMATION='confirm'
```
✅ validate-inputs (NEEDS_APPROVAL=true)
✅ build-images
✅ update-ten-apps-config
🔒 approval (REQUIRED - shows summary, waits for human approval)
  ├─ APPROVED → Continue
  └─ REJECTED → revert-config-on-rejection
✅ argocd-delete-resources
✅ argocd-sync-apps
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS (if approved) - With approval gate
```

### Permutation 10: SEPOLIA + NON-DESTRUCTIVE + IMAGE_BUILD=YES + HSM=YES + CONFIRMATION='confirm'
```
✅ validate-inputs (NEEDS_APPROVAL=true - mandatory for sepolia)
✅ build-images (WITH HSM signing)
✅ update-ten-apps-config
🔒 approval (MANDATORY - shows summary, waits)
  ├─ APPROVED → Continue
  └─ REJECTED → revert-config-on-rejection
✅ argocd-delete-resources
✅ argocd-sync-apps (only OutOfSync)
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS (if approved) - Production safe, HSM-signed
```

### Permutation 11: SEPOLIA + NON-DESTRUCTIVE + IMAGE_BUILD=NO + CONFIRMATION='confirm'
```
✅ validate-inputs (NEEDS_APPROVAL=true)
⏭️ build-images
⏭️ update-ten-apps-config
⏭️ approval (skipped - no config to approve)
✅ argocd-delete-resources
✅ argocd-sync-apps (existing images)
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS - Restart with existing images (no approval needed)
```

### Permutation 12: SEPOLIA + DESTRUCTIVE + IMAGE_BUILD=YES + HSM=YES + CONFIRMATION='confirm'
```
✅ validate-inputs (NEEDS_APPROVAL=true)
✅ build-images (WITH HSM)
✅ update-ten-apps-config (tags + node selectors)
🔒 approval (MANDATORY - detailed summary)
  ├─ APPROVED → Continue
  └─ REJECTED → revert-config-on-rejection + STOP
✅ deploy-l1-contracts
✅ argocd-delete-child-apps
✅ argocd-sync-apps (full sync)
✅ wait-argocd-healthy
✅ deploy-l2-contracts
✅ post-deployment

Result: ✅ SUCCESS (if approved) - Full production deployment with HSM
```

### Permutation 13: SEPOLIA + DESTRUCTIVE + IMAGE_BUILD=NO + CONFIRMATION='confirm'
```
✅ validate-inputs
⏭️ build-images
⏭️ update-ten-apps-config
⏭️ approval
❌ deploy-l1-contracts (BLOCKED)
❌ argocd-delete-child-apps (BLOCKED)
❌ argocd-sync-apps (BLOCKED)
❌ wait-argocd-healthy (BLOCKED)
❌ deploy-l2-contracts (BLOCKED)
❌ post-deployment

Result: ❌ FAILS - Destructive requires image_build=yes
```

### Permutation 14: MAINNET + NON-DESTRUCTIVE + IMAGE_BUILD=YES + HSM=YES + CONFIRMATION='confirm'
```
✅ validate-inputs (NEEDS_APPROVAL=true)
✅ build-images (WITH HSM - mandatory for mainnet)
✅ update-ten-apps-config
🔒 approval (MANDATORY - extra warnings)
  ├─ APPROVED → Continue
  └─ REJECTED → revert-config-on-rejection
✅ argocd-delete-resources
✅ argocd-sync-apps (only OutOfSync)
✅ wait-argocd-healthy
✅ post-deployment

Result: ✅ SUCCESS (if approved) - Production safe
```

### Permutation 15: MAINNET + DESTRUCTIVE + IMAGE_BUILD=YES + HSM=YES + CONFIRMATION='confirm'
```
✅ validate-inputs (NEEDS_APPROVAL=true)
✅ build-images (WITH HSM)
✅ update-ten-apps-config
🔒 approval (MANDATORY - shows DOWNTIME WARNING)
  ├─ APPROVED → Continue
  └─ REJECTED → revert-config-on-rejection
✅ deploy-l1-contracts (checks gas price)
✅ argocd-delete-child-apps
✅ argocd-sync-apps (full sync)
✅ wait-argocd-healthy
✅ deploy-l2-contracts
✅ post-deployment

Result: ✅ SUCCESS (if approved) - Full mainnet redeployment
```

### Permutation 16: MAINNET + DESTRUCTIVE + IMAGE_BUILD=NO + CONFIRMATION='confirm'
```
✅ validate-inputs
⏭️ build-images
⏭️ update-ten-apps-config
⏭️ approval
❌ deploy-l1-contracts (BLOCKED)
...

Result: ❌ FAILS - Destructive requires image_build=yes
```

---

## Key Insights

### ✅ Valid Combinations:
1. **Non-destructive + image_build=yes** → Full update with new images
2. **Non-destructive + image_build=no** → Restart with existing images
3. **Destructive + image_build=yes** → Full redeployment with new images

### ❌ Invalid Combinations:
1. **Destructive + image_build=no** → ❌ FAILS (L1 deployment blocked)

### 🔒 Approval Required:
- **Always**: sepolia-testnet, mainnet
- **Optional**: dev-testnet, uat-testnet (with confirmation='confirm')
- **Shows**: Detailed summary of what will be deployed
- **Auto-revert**: If rejected, reverts ten-apps config

### 🔐 HSM Signing:
- **Recommended**: sepolia-testnet, mainnet
- **Optional**: dev-testnet, uat-testnet
- **Effect**: Only affects enclave image build

---

## Recommendations

### For Development (dev-testnet):
```bash
# Quick iteration with new images
deployment_strategy: non-destructive
image_build: yes
use_azure_hsm: no

# Quick restart without rebuilding
deployment_strategy: non-destructive
image_build: no

# Full redeploy with contracts
deployment_strategy: destructive
image_build: yes
use_azure_hsm: no
```

### For UAT (uat-testnet):
```bash
# Standard update
deployment_strategy: non-destructive
image_build: yes
use_azure_hsm: no

# Full redeploy (test before production)
deployment_strategy: destructive
image_build: yes
use_azure_hsm: yes  # Test HSM signing
confirmation: confirm  # Enable approval
```

### For Sepolia (sepolia-testnet):
```bash
# Production-like update
deployment_strategy: non-destructive
image_build: yes
use_azure_hsm: yes
confirmation: confirm  # MANDATORY

# Full redeploy (rare)
deployment_strategy: destructive
image_build: yes
use_azure_hsm: yes
confirmation: confirm  # MANDATORY
```

### For Mainnet:
```bash
# Standard update (minimal downtime)
deployment_strategy: non-destructive
image_build: yes
use_azure_hsm: yes  # MANDATORY
confirmation: confirm  # MANDATORY

# Emergency full redeploy (causes downtime!)
deployment_strategy: destructive
image_build: yes
use_azure_hsm: yes  # MANDATORY
confirmation: confirm  # MANDATORY
max_gas_gwei: 5  # Check gas prices
```

---

## Workflow File Location
`.github/workflows/manual-deploy-k8s-consolidated.yml`

