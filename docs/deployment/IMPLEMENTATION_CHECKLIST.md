# Implementation Checklist

## Pre-Deployment Setup

### 1. GitHub Environment Configuration

**For each environment** (dev-testnet, uat-testnet, sepolia-testnet, mainnet):

- [ ] Go to Settings → Environments → Create environment (if not exists)
- [ ] Set environment name: `dev-testnet` / `uat-testnet` / `sepolia-testnet` / `mainnet`
- [ ] Add required secrets:
  - [ ] `ARGOCD_SERVER_NONPROD` (for dev/uat/sepolia): https://argocd-nonprod.internal
  - [ ] `ARGOCD_SERVER_PROD` (for mainnet): https://argocd-prod.internal
  - [ ] `ARGOCD_TOKEN`: Service account token
  - [ ] `DEPLOY_ACTIONS_PAT`: GitHub PAT for ten-apps repo access
  - [ ] `REGISTRY_PASSWORD`: Azure registry password
  - [ ] `AZURE_CREDENTIALS`: Azure auth JSON
  - [ ] `L1_HTTP_URL`: L1 RPC endpoint
  - [ ] `ACCOUNT_PK_WORKER`: Deployer private key
  - [ ] `ETHERSCAN_API_KEY`: For contract verification
  - [ ] `L2_DEPLOYER_KEY`: L2 deployer key
  - [ ] `GH_TOKEN`: GitHub token for ten-test dispatch

- [ ] Add required variables (if using):
  - [ ] `ACCOUNT_ADDR_NODE_0`: Initial sequencer address
  - [ ] `CHAIN_ID`: Network chain ID
  - [ ] `DOCKER_BUILD_TAG_*`: Image build tags (optional, auto-generated)
  - [ ] `L1_CHALLENGE_PERIOD`: Challenge period value
  - [ ] `FAUCET_INITIAL_FUNDS`: Faucet funding amount

- [ ] **For sepolia-testnet & mainnet ONLY**:
  - [ ] Enable "Required reviewers"
  - [ ] Add team members as reviewers
  - [ ] Set review timeout (e.g., 24 hours)

### 2. ArgoCD Service Account Setup

**For each ArgoCD instance** (nonprod and prod):

```bash
# 1. Create namespace (if not exists)
kubectl create namespace argocd

# 2. Create service account
kubectl create serviceaccount github-actions -n argocd

# 3. Grant permissions
kubectl create clusterrolebinding github-actions \
  --clusterrole=cluster-admin \
  --serviceaccount=argocd:github-actions

# 4. Get token (valid for 1 year)
TOKEN=$(kubectl -n argocd create token github-actions --duration=8760h)
echo $TOKEN

# 5. Copy token to GitHub Secrets
# Paste in ARGOCD_TOKEN for corresponding environment
```

- [ ] Token created for nonprod ArgoCD → Copy to `ARGOCD_TOKEN` in dev-testnet env
- [ ] Token created for nonprod ArgoCD → Copy to `ARGOCD_TOKEN` in uat-testnet env
- [ ] Token created for nonprod ArgoCD → Copy to `ARGOCD_TOKEN` in sepolia-testnet env
- [ ] Token created for prod ArgoCD → Copy to `ARGOCD_TOKEN` in mainnet env

### 3. Verify Ten-Apps Repository Access

- [ ] `DEPLOY_ACTIONS_PAT` has `repo` scope
- [ ] PAT has write access to ten-apps repo
- [ ] Test: `curl -H "Authorization: token $PAT" https://api.github.com/repos/ten-protocol/ten-apps`
- [ ] Response shows repo details (not 401)

### 4. Verify Docker Registry Access

- [ ] `REGISTRY_PASSWORD` is valid for `testnetobscuronet.azurecr.io`
- [ ] Test: `docker login -u testnetobscuronet -p $PASSWORD testnetobscuronet.azurecr.io`
- [ ] Can push images: `docker push testnetobscuronet.azurecr.io/obscuronet/test:latest`

### 5. Verify Azure Credentials

- [ ] `AZURE_CREDENTIALS` is valid JSON (service principal)
- [ ] Contains: `clientId`, `clientSecret`, `subscriptionId`, `tenantId`
- [ ] Test: `az login --service-principal -u $clientId -p $clientSecret --tenant $tenantId`

---

## Testing & Validation

### Phase 1: Dev Non-Destructive (No Risk)

- [ ] Run workflow: `k8s Deploy Consolidated`
- [ ] Inputs:
  - `testnet_type`: dev-testnet
  - `deployment_strategy`: non-destructive
  - `image_build`: no
  - `confirmation`: (leave blank)

- [ ] Verify workflow completes:
  - [ ] All jobs completed successfully
  - [ ] No approval gate appeared (expected)
  - [ ] L1 contracts deployed (check logs)
  - [ ] ArgoCD apps synced (check logs)

- [ ] Verify in ten-apps repo:
  - [ ] No git commits (expected, no image build)

- [ ] Verify in ArgoCD:
  ```bash
  argocd app get dev-sequencer --refresh
  # status.health.status should be "Healthy"
  # status.sync.status should be "Synced"
  ```

### Phase 2: Dev Destructive with Images

- [ ] Run workflow: `k8s Deploy Consolidated`
- [ ] Inputs:
  - `testnet_type`: dev-testnet
  - `deployment_strategy`: destructive
  - `image_build`: yes
  - `image_tag`: v1.5.8.0
  - `confirmation`: (leave blank)

- [ ] Verify workflow completes:
  - [ ] No approval gate (dev is auto-approved)
  - [ ] Images built and pushed
  - [ ] L1 contracts deployed
  - [ ] Child apps deleted (check logs for each)
  - [ ] Child apps recreated and synced
  - [ ] Apps reached Healthy state
  - [ ] L2 contracts deployed

- [ ] Verify in ten-apps repo:
  - [ ] Git commit with message "chore: update image tags to v1.5.8.0 for dev"
  - [ ] YAML files updated with new tag

- [ ] Verify in ArgoCD:
  ```bash
  argocd app get dev-sequencer
  # Should show new image tag: v1.5.8.0
  argocd app get dev-validator-01
  argocd app get dev-validator-02
  argocd app get dev-gateway
  argocd app get dev-tools
  # All should be Healthy + Synced
  ```

- [ ] Verify network is functional:
  - [ ] Connect to dev network
  - [ ] Send test transaction
  - [ ] Verify in block explorer

### Phase 3: UAT Non-Destructive

- [ ] Run workflow: `k8s Deploy Consolidated`
- [ ] Inputs:
  - `testnet_type`: uat-testnet
  - `deployment_strategy`: non-destructive
  - `image_build`: no
  - `confirmation`: (leave blank)

- [ ] Verify workflow completes
- [ ] Verify ArgoCD apps synced

### Phase 4: Sepolia Destructive (Approval Testing)

- [ ] Run workflow: `k8s Deploy Consolidated`
- [ ] Inputs:
  - `testnet_type`: sepolia-testnet
  - `deployment_strategy`: destructive
  - `image_build`: yes
  - `image_tag`: v1.5.8.0
  - `confirmation`: confirm ← **REQUIRED**

- [ ] Verify approval gate appeared:
  - [ ] GitHub Actions shows pending approval
  - [ ] Approval notification sent (if configured)
  - [ ] Shows "Waiting for approval from: <reviewer>"

- [ ] Approve workflow:
  - [ ] As configured reviewer, go to Actions → Workflow run → Review deployments
  - [ ] Click "Approve and deploy"

- [ ] Verify workflow continues after approval
- [ ] Verify sepolia apps updated

### Phase 5: Mainnet Destructive (FINAL CHECK - Production)

- [ ] **COORDINATE WITH TEAM BEFORE RUNNING**
- [ ] Schedule maintenance window
- [ ] Notify users of planned downtime
- [ ] Have rollback plan ready

- [ ] Run workflow: `k8s Deploy Consolidated`
- [ ] Inputs:
  - `testnet_type`: mainnet
  - `deployment_strategy`: destructive
  - `image_build`: yes
  - `image_tag`: v1.5.8.0
  - `confirmation`: confirm ← **REQUIRED**

- [ ] Verify approval gate appeared and approve

- [ ] Monitor deployment:
  - [ ] Apps are recreating
  - [ ] No errors in logs
  - [ ] Contracts deployed successfully
  - [ ] Network is operational

---

## Post-Deployment Verification

### After Each Deployment

- [ ] GitHub Actions workflow shows ✅ (all green)
- [ ] No failed jobs
- [ ] Artifact uploads succeeded
- [ ] ArgoCD shows all apps as "Healthy" + "Synced"

```bash
# Verify all apps healthy
for app in sequencer validator-01 validator-02 gateway tools; do
  argocd app get $app --refresh | grep -E "Health|Sync"
done
```

- [ ] Network is functional
  - [ ] Can query node RPC
  - [ ] Can see blocks being produced
  - [ ] Can submit transactions

- [ ] Contract addresses correct
  - [ ] Compare with ten-apps values files
  - [ ] Verify on block explorer

### Rollback Test (Optional but Recommended)

- [ ] After successful destructive on dev
- [ ] Run non-destructive to ensure it works
- [ ] Run destructive with previous image tag
- [ ] Verify apps recreated with old images
- [ ] Confirms rollback procedure works

---

## Documentation & Handoff

### Documentation

- [ ] Team has read `QUICK_REFERENCE.md`
- [ ] Team has read `DEPLOYMENT_GUIDE.md`
- [ ] Team has read `CONSOLIDATION_SUMMARY.md`
- [ ] Team has access to this checklist

### Runbooks Updated

- [ ] Update oncall runbook to reference new workflow
- [ ] Remove references to old 3-step manual process
- [ ] Document approval gate process

### Team Training

- [ ] Team walkthrough of new workflow
- [ ] Demo of inputs and approval process
- [ ] Q&A session

### Access Control

- [ ] Verified reviewers list for prod environments
- [ ] Confirmed approval notification settings
- [ ] Tested approval process works

---

## Cleanup (After All Tests Pass)

### Archive Old Workflows

```bash
cd .github/workflows

# Rename (keep as reference for 2 weeks)
mv manual-deploy-k8s-testnet-before-nodes.yml manual-deploy-k8s-testnet-before-nodes.yml.archive
mv manual-deploy-k8s-testnet-after-nodes.yml manual-deploy-k8s-testnet-after-nodes.yml.archive
mv build-release-images.yml build-release-images.yml.archive
```

- [ ] Old workflow files archived (renamed with `.archive`)
- [ ] Commit change with message: "archive: old deployment workflows"

### Final Verification

- [ ] New workflow runs successfully
- [ ] All team members can see it in GitHub Actions
- [ ] Old workflows no longer appear in regular lists

---

## Troubleshooting During Testing

### Issue: "Confirmation field must say 'confirm'"
- [ ] Verify you typed exactly: `confirm` (not "Confirm", not "confirmed")
- [ ] Check for extra spaces

### Issue: Workflow not showing in Actions
- [ ] Verify file is in `.github/workflows/` directory
- [ ] Verify filename ends with `.yml`
- [ ] Commit and push to repo
- [ ] Refresh GitHub Actions page

### Issue: "App not found" during destructive
- [ ] This is normal if app was already deleted
- [ ] Workflow continues safely
- [ ] Check logs show "App X not found or already deleted"

### Issue: ArgoCD sync times out
- [ ] Increase `sync_timeout` to 15m or 20m
- [ ] Check if K8s cluster is responsive
- [ ] Check for pod errors: `kubectl get pods -n dev`

### Issue: Approval not appearing
- [ ] Verify environment has "Required reviewers" enabled
- [ ] Verify you selected "destructive" strategy
- [ ] Verify you're targeting sepolia or mainnet
- [ ] Check GitHub environment settings

### Issue: Images not pushed to registry
- [ ] Verify `REGISTRY_PASSWORD` is correct
- [ ] Verify docker login works manually
- [ ] Check Azure subscription limits not exceeded

---

## Success Criteria

✅ **All of these must be true**:

1. **Workflow executes end-to-end** without failures
2. **Images are built** and tagged correctly in registry
3. **Ten-apps is updated** with new image tags (git commit)
4. **L1 contracts deployed** successfully
5. **ArgoCD apps synced** to "Healthy" + "Synced" state
6. **L2 contracts deployed** (for destructive only)
7. **Approval gate works** (for prod destructive)
8. **Network is operational** after deployment
9. **All logs are clear** (no error messages)
10. **Team is trained** on new workflow

---

## Sign-Off

- [ ] QA Lead: Tested and verified workflow _________________ Date: _______
- [ ] DevOps Lead: Verified infrastructure setup _________________ Date: _______
- [ ] Team Lead: Reviewed documentation _________________ Date: _______
- [ ] Deployment Ready ✅

---

## After Implementation

### Week 1: Monitor
- [ ] Daily check that deployments work
- [ ] Monitor for any issues
- [ ] Collect team feedback

### Week 2: Archive
- [ ] Delete archived old workflow files
- [ ] Update all documentation
- [ ] Announce consolidation complete

### Week 3+: Maintain
- [ ] Use consolidated workflow exclusively
- [ ] Archive this checklist
- [ ] Update runbooks as needed
