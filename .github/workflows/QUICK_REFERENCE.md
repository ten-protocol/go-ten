# Consolidated Deployment - Quick Reference

## TL;DR - Most Common Operations

### 1. Update Images on Dev (Non-Destructive)
```
Workflow: k8s Deploy Consolidated

Inputs:
  testnet_type: dev-testnet
  deployment_strategy: non-destructive
  image_build: yes
  image_tag: v1.5.8.1

Result: 20 min | New images + L1 deployed | ArgoCD apps synced | NO L2 reset
```

### 2. Full Dev Reset (Destructive)
```
Workflow: k8s Deploy Consolidated

Inputs:
  testnet_type: dev-testnet
  deployment_strategy: destructive
  image_build: yes
  image_tag: v1.5.8.1

Result: 37 min | New images + L1 + L2 deployed | Apps recreated | Full network reset
```

### 3. Update Dev Config Only (Non-Destructive, No Images)
```
Workflow: k8s Deploy Consolidated

Inputs:
  testnet_type: dev-testnet
  deployment_strategy: non-destructive
  image_build: no
  image_tag: (leave blank)

Result: 15 min | L1 deployed | ArgoCD apps synced with new config from ten-apps
```

### 4. Production Destructive (Requires Approval)
```
Workflow: k8s Deploy Consolidated

Inputs:
  testnet_type: mainnet
  deployment_strategy: destructive
  image_build: yes
  image_tag: v1.5.8.1
  confirmation: confirm

Result: (approval wait) + 37 min | Full reset | Apps recreated | L2 deployed
```

---

## Decision Tree

```
START
  │
  ├─ Want to update images/config only?
  │   ├─ YES → deployment_strategy: non-destructive
  │   │         (✅ ArgoCD apps synced, ❌ NO L2 reset)
  │   │
  │   └─ NO → deployment_strategy: destructive
  │            (✅ Full reset with L2, ❌ needs approval for prod)
  │
  ├─ Build new images?
  │   ├─ YES → image_build: yes, image_tag: v1.5.8.1
  │   └─ NO → image_build: no
  │
  ├─ Target production (sepolia/mainnet)?
  │   ├─ YES → confirmation: "confirm" (type exactly)
  │   └─ NO → confirmation: (leave blank)
  │
  ├─ Hit "Run workflow" ▶️
  │
  └─ Wait for completion ⏳
```

---

## What Gets Deleted (Destructive Only)

### Default (Always Deleted)
- ✅ `sequencer`
- ✅ `validator-01`
- ✅ `validator-02`
- ❌ Parent app (e.g., `ten-dev-apps`) - NOT deleted

### Optional (If Specified)
- ✅ `gateway` (if in `additional_apps_to_delete`)
- ✅ `tools` (if in `additional_apps_to_delete`)
- ✅ `postgres-client` (if in `additional_apps_to_delete`)
- ✅ Any other apps (if in `additional_apps_to_delete`)

### Examples by Environment

**Dev (default deletion):**
- ✅ `dev-sequencer`
- ✅ `dev-validator-01`
- ✅ `dev-validator-02`
- ❓ `dev-gateway` (only if in `additional_apps_to_delete: gateway`)
- ❓ `dev-tools` (only if in `additional_apps_to_delete: tools`)

### UAT
- ✅ `uat-sequencer`
- ✅ `uat-validator-01`
- ✅ `uat-validator-02`
- ✅ `uat-gateway`
- ✅ `uat-tools`
- ❌ `ten-uat-apps` (parent, NOT deleted)

### Sepolia
- ✅ `sepolia-sequencer`
- ✅ `sepolia-validator-01`
- ✅ `sepolia-validator-02`
- ✅ `sepolia-gateway`
- ✅ `sepolia-gateway-dexynth`
- ✅ `sepolia-gateway-pentest`
- ✅ `sepolia-tools`
- ❌ `ten-sepolia-apps` (parent, NOT deleted)

### Mainnet
- ✅ `mainnet-sequencer`
- ✅ `mainnet-validator-01`
- ✅ `mainnet-validator-02`
- ✅ `mainnet-gateway`
- ✅ `mainnet-postgres-client`
- ✅ `mainnet-tools`
- ❌ `ten-mainnet-apps` (parent, NOT deleted)

---

## Approval Requirements Matrix

|  | Non-Destructive | Destructive |
|---|---|---|
| **dev-testnet** | ❌ No | ❌ No |
| **uat-testnet** | ❌ No | ❌ No |
| **sepolia-testnet** | ❌ No | ✅ YES |
| **mainnet** | ❌ No | ✅ YES |

**For approval**: Set `confirmation: "confirm"` (type exactly)

---

## Expected Timeline

| Phase | Non-Destructive | Destructive |
|-------|---|---|
| Build images | 5 min | 5 min |
| Update ten-apps | 2 min | 2 min |
| Deploy L1 | 10 min | 10 min |
| Delete apps | ❌ - | 2 min |
| Sync ArgoCD | 3 min | 3 min |
| Wait healthy | ❌ - | 5 min |
| Deploy L2 | ❌ - | 10 min |
| **TOTAL** | **~20 min** | **~37 min** |

Add approval wait time for sepolia/mainnet destructive.

---

## What Happens Automatically

- ✅ Images are built and pushed to registry
- ✅ Image tags are updated in ten-apps YAML
- ✅ Changes are committed to ten-apps repo
- ✅ L1 contracts are deployed
- ✅ ArgoCD apps are synced
- ✅ Health checks are performed
- ✅ L2 contracts deployed (destructive only)
- ✅ Test repo notification sent (dev/uat only)

---

## What You Need to Provide

- ✅ `testnet_type` - Pick one: dev, uat, sepolia, mainnet
- ✅ `deployment_strategy` - Pick: non-destructive or destructive
- ✅ `image_build` - Pick: yes or no
- ❓ `image_tag` - Optional (auto-detect if blank)
- ❓ `confirmation` - Required if destructive + sepolia/mainnet
- ❓ `log_level` - Optional (default: 3)
- ❓ `max_gas_gwei` - Optional (default: 1.5)
- ❓ `sync_timeout` - Optional (default: 10m)
- ❓ `sequencer_node_selector` - Optional (leave blank to keep current)
- ❓ `validator_01_node_selector` - Optional (leave blank to keep current)
- ❓ `validator_02_node_selector` - Optional (leave blank to keep current)
- ❓ `gateway_node_selector` - Optional (leave blank to keep current)
- ❓ `additional_apps_to_delete` - Optional, destructive only (comma-separated, e.g., "gateway,tools")

---

## Troubleshooting - Quick Fixes

| Problem | Solution |
|---------|----------|
| "Confirmation field must say 'confirm'" | Type exactly: `"confirm"` in confirmation field |
| "App not found or already deleted" | Normal during destructive, safe to continue |
| "Failed to sync <app> (timeout)" | Increase sync_timeout to 15m or 20m |
| "App did not reach healthy state" | Check K8s pods: `kubectl get pods -n <env>` |
| "L2 contracts not deployed" | Only runs for destructive. Use destructive strategy. |
| "Image tag not updated in ten-apps" | Check git log in ten-apps repo for commit |

---

## Rollback

If something goes wrong:

```bash
# Option 1: Config-only rollback (non-destructive)
# Re-run workflow with:
deployment_strategy: non-destructive
image_build: no

# Option 2: Full reset (destructive)
# Re-run with fresh image tag
# This deletes and recreates everything

# Option 3: Manual rollback
cd ten-apps
git revert <bad-commit-hash>
git push
# Then run non-destructive to pick up changes
```

---

## Common Patterns

### Pattern 1: Weekly Image Update
```
Monday 9am:
  - Build new images from main branch
  - Non-destructive sync
  - Verify in dev overnight

Tuesday 9am:
  - If good, repeat on uat
```

### Pattern 2: Major Release
```
1. Build images locally, tag v1.5.8.0
2. Run destructive on dev with v1.5.8.0
3. Test for 2 hours
4. Run destructive on uat with v1.5.8.0
5. Approval from team lead
6. Run destructive on mainnet with v1.5.8.0
```

### Pattern 3: Config-Only Fix
```
1. Update value in ten-apps YAML (e.g., gas price)
2. Run non-destructive on dev
3. Verify config change applied
4. Repeat on uat/mainnet as needed
```

### Pattern 4: Destructive with Additional App Deletion
```
1. Run destructive deployment
2. Fill in additional_apps_to_delete (optional)
3. Examples:
   - Leave blank: only delete sequencer, validator-01, validator-02
   - "gateway": delete those 3 + gateway
   - "gateway,tools": delete those 3 + gateway + tools
   - "gateway,tools,postgres-client": delete those 3 + all others
4. Default apps (seq, val-01, val-02) ALWAYS deleted if destructive
5. Parent app NEVER deleted
```

### Pattern 5: Move Nodes to Different Hardware (Destructive)
```
1. List available nodes: kubectl get nodes -L kubernetes.io/hostname
2. Run destructive deployment with:
   - sequencer_node_selector: <node-hostname>
   - validator_01_node_selector: <node-hostname>
   - validator_02_node_selector: <node-hostname>
   - gateway_node_selector: <node-hostname> (optional)
3. Default apps deleted and recreated on new nodes
4. Changes committed to ten-apps automatically
```

---

## GitHub Actions UI

Navigate to:
```
GitHub → Actions → "[M] k8s Deploy Consolidated" → Run workflow
```

Fill in the form and click **"Run workflow"** button.

Workflow appears in the list instantly. Click it to watch live logs.

---

## Important Notes

⚠️ **Non-destructive:**
- Does NOT reset networks
- Does NOT redeploy L2 contracts
- Safe for config/image updates

⚠️ **Destructive:**
- Deletes all child apps
- Redeploys everything from scratch
- Requires approval for production
- Takes ~37 minutes

⚠️ **Production (mainnet/sepolia):**
- Requires `confirmation: "confirm"`
- Requires human approval
- Coordinate with team before running

---

## Gotchas

1. **Non-destructive does NOT deploy L2**
   - Even if you run it, L2 won't be deployed
   - Use destructive for full reset

2. **Image tag is automatic if left blank**
   - Finds latest git tag
   - Specify explicitly for reproducibility

3. **Parent apps are never deleted**
   - Only child apps (sequencer, validator, gateway, etc.)
   - Parent apps (ten-dev-apps, etc.) stay intact

4. **Confirmation must be exact: `"confirm"`**
   - Not "Confirm" or "confirmed"
   - Exactly: `confirm`

5. **L2 deployment only runs after ArgoCD is healthy**
   - Workflow waits up to 5 minutes
   - Fails if apps don't reach Healthy state

---

## Next Steps

1. Read full `DEPLOYMENT_GUIDE.md`
2. Test on dev-testnet (non-destructive first)
3. Test on dev-testnet (destructive)
4. Verify ArgoCD apps synced correctly
5. Report any issues to DevOps team

---

**Questions?** Check DEPLOYMENT_GUIDE.md or CONSOLIDATION_SUMMARY.md for detailed info.
