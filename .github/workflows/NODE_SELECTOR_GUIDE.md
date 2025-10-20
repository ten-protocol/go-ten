# Node Selector Management Guide (SGX Nodes)

## Overview

The consolidated deployment workflow supports **optional node selector updates** for **DESTRUCTIVE deployments only**. This allows you to:

- ✅ Ensure pods recreate on specific SGX-enabled nodes
- ✅ Pin sequencer/validators/gateway to correct hardware
- ✅ Rebuild and reschedule pods to assigned nodes
- ✅ Changes auto-committed to ten-apps repo
- ✅ Only specify nodes for pods you want to move
- ⚠️ **NOTE**: Non-destructive deployments ignore node selectors (pods stay where they are)

---

## Node Selector Inputs

Available in the deployment workflow UI:

| Input | Description | Example |
|-------|-------------|---------|
| `sequencer_node_selector` | Target node for sequencer pods | `aks-sgxpool01-61714098-vmss000004` |
| `validator_01_node_selector` | Target node for validator-01 pods | `aks-sgxpool01-61714098-vmss000005` |
| `validator_02_node_selector` | Target node for validator-02 pods | `aks-sgxpool02-25481487-vmss000004` |
| `gateway_node_selector` | Target node for gateway pods | `aks-pool-25481487-vmss000001` |

**All node selector inputs are optional:**
- Leave blank to keep current assignment
- Only update the nodes you want to move
- Empty fields are ignored (no changes made)

---

## How It Works

### Flow

```
1. User fills in deployment form:
   - deployment_strategy: non-destructive (recommended)
   - image_build: no (if only moving nodes)
   - sequencer_node_selector: aks-sgxpool01-61714098-vmss000006
   - (leave other node selectors blank)

2. Workflow executes:
   - Downloads ten-apps repo
   - Updates ONLY the sequencer values file
   - Sets: enclave.nodeSelector, enclave02.nodeSelector, host.nodeSelector
   - Leaves validator & gateway unchanged

3. Git commit:
   - Message: "chore: update node selectors for dev"
   - Commits to ten-apps repo

4. ArgoCD sync:
   - Detects config changes
   - Terminates old sequencer pods
   - Schedules new sequencer pods on target node
   - Validators & gateway unaffected
```

### YAML Mutation

**Before:**
```yaml
enclave:
  nodeSelector:
    kubernetes.io/hostname: aks-sgxpool01-61714098-vmss000004

host:
  nodeSelector:
    kubernetes.io/hostname: aks-sgxpool01-61714098-vmss000004
```

**After (if `sequencer_node_selector: aks-sgxpool01-61714098-vmss000006`):**
```yaml
enclave:
  nodeSelector:
    kubernetes.io/hostname: aks-sgxpool01-61714098-vmss000006

host:
  nodeSelector:
    kubernetes.io/hostname: aks-sgxpool01-61714098-vmss000006
```

---

## Finding Node Hostnames

### List Available Nodes

```bash
# Show all nodes with hostname labels
kubectl get nodes -L kubernetes.io/hostname

# Example output:
NAME                                STATUS   ROLES    AGE     HOSTNAME
aks-sgxpool01-61714098-vmss000004   Ready    agent    200d    aks-sgxpool01-61714098-vmss000004
aks-sgxpool01-61714098-vmss000005   Ready    agent    200d    aks-sgxpool01-61714098-vmss000005
aks-sgxpool02-25481487-vmss000004   Ready    agent    180d    aks-sgxpool02-25481487-vmss000004
aks-pool-25481487-vmss000001        Ready    agent    90d     aks-pool-25481487-vmss000001
```

### Check Current Node Assignments

```bash
# For dev environment
kubectl get pods -n dev -o wide | grep -E "sequencer|validator|gateway"

# Example output:
dev-sequencer-0            Running   aks-sgxpool01-61714098-vmss000004
dev-validator-01-0         Running   aks-sgxpool01-61714098-vmss000005
dev-validator-02-0         Running   aks-sgxpool02-25481487-vmss000004
dev-gateway-0              Running   aks-pool-25481487-vmss000001
```

### Check Node Resources

```bash
# View node capacity
kubectl describe node aks-sgxpool01-61714098-vmss000004

# View current pods on node
kubectl get pods --all-namespaces --field-selector spec.nodeName=aks-sgxpool01-61714098-vmss000004
```

---

## Common Scenarios

### Scenario 1: Move Sequencer to New Hardware

**Situation:** New SGX node added to cluster, need to move sequencer there

**Steps:**
1. List available nodes:
   ```bash
   kubectl get nodes -L kubernetes.io/hostname
   ```

2. Identify new node hostname: `aks-sgxpool01-61714098-vmss000006`

3. Run deployment:
   - `testnet_type`: dev-testnet
   - `deployment_strategy`: non-destructive
   - `image_build`: no
   - `sequencer_node_selector`: aks-sgxpool01-61714098-vmss000006
   - (leave other node selectors blank)

4. Monitor pod movement:
   ```bash
   kubectl get pods -n dev -o wide -w | grep sequencer
   ```

5. Verify in GitHub:
   - Check ten-apps commit message
   - YAML updated with new hostname

---

### Scenario 2: Move All Nodes (Maintenance)

**Situation:** Entire pool going down for maintenance, need to move all pods to backup pool

**Steps:**
1. List backup nodes:
   ```bash
   kubectl get nodes -L kubernetes.io/hostname | grep backup-pool
   ```

2. Identify backup nodes:
   - Sequencer: `aks-backup-pool-123-vmss000001`
   - Validator-01: `aks-backup-pool-123-vmss000002`
   - Validator-02: `aks-backup-pool-123-vmss000003`
   - Gateway: `aks-backup-pool-123-vmss000004`

3. Run deployment:
   - `testnet_type`: dev-testnet
   - `deployment_strategy`: non-destructive
   - `image_build`: no
   - `sequencer_node_selector`: aks-backup-pool-123-vmss000001
   - `validator_01_node_selector`: aks-backup-pool-123-vmss000002
   - `validator_02_node_selector`: aks-backup-pool-123-vmss000003
   - `gateway_node_selector`: aks-backup-pool-123-vmss000004

4. Monitor:
   ```bash
   kubectl get pods -n dev -o wide -w
   ```

---

### Scenario 3: Move Only Validators (Sequencer Stays)

**Situation:** Sequencer is happy where it is, but validators need to move

**Steps:**
1. Run deployment:
   - `testnet_type`: uat-testnet
   - `deployment_strategy`: non-destructive
   - `image_build`: no
   - `sequencer_node_selector`: (leave blank - no change)
   - `validator_01_node_selector`: aks-new-pool-456-vmss000005
   - `validator_02_node_selector`: aks-new-pool-456-vmss000006
   - `gateway_node_selector`: (leave blank - no change)

2. Result:
   - ✅ Validators move to new nodes
   - ✅ Sequencer stays on current node
   - ✅ Gateway stays on current node

---

### Scenario 4: Combined Image Update + Node Move

**Situation:** New image AND need to move pods to different hardware

**Steps:**
1. Run deployment:
   - `testnet_type`: dev-testnet
   - `deployment_strategy`: non-destructive
   - `image_build`: yes
   - `image_tag`: v1.5.8.1
   - `sequencer_node_selector`: aks-new-pool-789-vmss000001
   - (other node selectors blank or filled as needed)

2. Result:
   - ✅ New images built and pushed
   - ✅ Image tags updated in ten-apps
   - ✅ Node selectors updated in ten-apps
   - ✅ ArgoCD pulls new images on new nodes
   - ✅ Single git commit: "chore: update image tags to v1.5.8.1 and node selectors for dev"

---

## What Gets Updated

### Updated in ten-apps YAML

When you specify a node selector:

```yaml
# For sequencer/validator-01/validator-02:
enclave:
  nodeSelector:
    kubernetes.io/hostname: <NEW_HOSTNAME>

enclave02:
  nodeSelector:
    kubernetes.io/hostname: <NEW_HOSTNAME>

host:
  nodeSelector:
    kubernetes.io/hostname: <NEW_HOSTNAME>

# For gateway:
gateway:
  nodeSelector:
    kubernetes.io/hostname: <NEW_HOSTNAME>
```

### Not Updated

- ✅ Pod requests/limits (unchanged)
- ✅ Image tags (unchanged unless you build new images)
- ✅ Other configurations (unchanged)
- ✅ Other pods (unchanged)

---

## Commit Message Examples

**Images only:**
```
chore: update image tags to v1.5.8.1 for dev
```

**Node selectors only:**
```
chore: update node selectors for dev
```

**Both:**
```
chore: update image tags to v1.5.8.1 and node selectors for dev
```

**No changes:**
```
ℹ️  No changes to commit
```

---

## Pod Rescheduling Timeline

When you move a pod to a new node:

```
T+0:00   - Deployment workflow runs
          - ten-apps YAML updated
          - Git commit pushed

T+0:30   - ArgoCD detects change (within sync interval)
          - Updates application state

T+1:00   - Pod termination begins
          - Grace period: 30s
          - Pod terminates on old node

T+1:30   - Kubernetes schedules pod on new node
          - Pulls image (if not cached)
          - Starts container

T+2:00   - Pod ready on new node
          - Service traffic redirected

Total: ~2 minutes for complete pod migration
```

---

## Monitoring Node Moves

### Watch Pod Migration

```bash
# Real-time pod status (use -w for watch mode)
kubectl get pods -n dev -o wide -w | grep sequencer

# Show pod events
kubectl describe pod -n dev dev-sequencer-0

# Check pod logs (after rescheduling)
kubectl logs -n dev dev-sequencer-0 -c host --tail 50
```

### Verify Node Assignment

```bash
# Check current node for pod
kubectl get pod -n dev dev-sequencer-0 -o jsonpath='{.spec.nodeName}'

# Output: aks-sgxpool01-61714098-vmss000006
```

### Check ArgoCD Status

```bash
# Check if ArgoCD picked up changes
argocd app get dev-sequencer

# Status should show:
# Sync: Synced
# Health: Healthy
```

---

## Troubleshooting

### Issue: Pod Not Rescheduling to Target Node

**Symptoms:**
- Pod stays on old node after deployment
- `kubectl get pods -o wide` shows old node name

**Solutions:**
```bash
# 1. Check if node exists
kubectl get nodes -L kubernetes.io/hostname | grep <target-hostname>

# 2. Check node capacity
kubectl describe node <target-hostname>
# Look for: Allocatable, Conditions

# 3. Check for node selectors in ten-apps
argocd app get dev-sequencer -o yaml | grep -A 10 nodeSelector

# 4. Manually trigger sync
argocd app sync dev-sequencer --wait

# 5. Force pod recreation (last resort)
kubectl delete pod -n dev dev-sequencer-0
# New pod will start per updated configuration
```

### Issue: Target Node Not Available

**Symptoms:**
- Pod tries to reschedule but stays Pending
- `kubectl describe pod` shows "0/N nodes are available"

**Solutions:**
```bash
# 1. List healthy nodes
kubectl get nodes -o wide

# 2. Check node conditions
kubectl describe node <node-name> | grep -A 10 Conditions

# 3. Choose different node
# Re-run deployment with available node

# 4. Check resource requirements
kubectl describe node <target-node> | grep Allocated
```

### Issue: Wrong Node Selected

**Symptoms:**
- Pod on unexpected node
- Git commit shows wrong hostname

**Solutions:**
```bash
# 1. Run deployment again with correct hostname
deployment_strategy: non-destructive
image_build: no
sequencer_node_selector: <CORRECT_HOSTNAME>

# 2. Workflow will update YAML to correct value
# 3. ArgoCD will reschedule pod again
```

---

## Best Practices

### ✅ Do

- Use non-destructive deployments for node moves
- Verify node capacity before moving pods
- Update one pod type at a time (sequencer → validators → gateway)
- Monitor pods for 5 minutes after rescheduling
- Keep audit trail: each move commits to git

### ❌ Don't

- Move all pods simultaneously (unless necessary)
- Specify wrong node hostnames (copy-paste carefully!)
- Move pods to maintenance nodes without prior notice
- Forget to verify ten-apps commit was pushed
- Change node selector in ArgoCD UI manually (use workflow instead)

---

## Git Verification

After deployment, verify changes in ten-apps:

```bash
cd ten-apps

# Show latest commit
git log --oneline -1

# Show changed files
git show --name-only

# Show actual changes
git show HEAD
```

Expected output:
```
chore: update node selectors for dev

nonprod-argocd-config/apps/envs/dev/valuesFile/values-dev-sequencer.yaml
nonprod-argocd-config/apps/envs/dev/valuesFile/values-dev-validator-01.yaml
nonprod-argocd-config/apps/envs/dev/valuesFile/values-dev-validator-02.yaml
```

---

## Network Impact

Moving pods has **minimal network impact**:

- ✅ No downtime (pods migrate gracefully)
- ✅ Service endpoints updated automatically
- ✅ Existing connections drain gracefully
- ✅ New connections route to new pod location
- ✅ DNS resolves to stable service, not pod directly

---

## Performance After Move

First minutes after pod migration:

- ⏳ Initial sync period: ~30-60 seconds
- ⏳ State rehydration: ~2-5 minutes
- ⏳ Full performance: ~5-10 minutes

During this time:
- Pod is operational but catching up state
- May see increased latency temporarily
- Block production continues (no reorg)
- Safe to monitor, not safe to heavily load

---

## Rollback

If node move caused issues:

### Option 1: Move Back to Original Node

```bash
# Identify original node
kubectl get pods -n dev --field-selector metadata.name=dev-sequencer-0 -o jsonpath='{.spec.affinity}'

# Run deployment with original node selector
sequencer_node_selector: <ORIGINAL_HOSTNAME>
```

### Option 2: Revert Git Commit

```bash
cd ten-apps
git revert <bad-commit-hash>
git push

# Then sync ArgoCD
argocd app sync dev-sequencer --wait
```

---

## Support

For issues or questions:
1. Check pod logs: `kubectl logs -n <env> <pod-name>`
2. Check node status: `kubectl describe node <hostname>`
3. Check ArgoCD app: `argocd app get <app-name>`
4. Review git history: `cd ten-apps && git log`
5. Contact DevOps team with pod names and error messages
