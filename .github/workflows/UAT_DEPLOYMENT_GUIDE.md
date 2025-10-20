# UAT Deployment Guide - Focus on Destructive Testing

## 🎯 Overview

This guide focuses on **UAT destructive deployments** which will be your primary testing strategy.

### Key Points
- ✅ UAT is NOT production (no approval gate)
- ✅ Destructive is the main test scenario
- ✅ Full network reset each time
- ✅ Test before moving to production
- ✅ No downtime impact

---

## 🚀 Typical UAT Destructive Deployment

### Standard Setup
```
testnet_type: uat-testnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.8.0
additional_apps_to_delete: (blank - delete default 3 apps only)
confirmation: (leave blank - UAT doesn't need approval)
```

**What happens:**
1. Build images with tag `v1.5.8.0`
2. Update ten-apps YAML
3. Deploy L1 contracts
4. Delete: uat-sequencer, uat-validator-01, uat-validator-02
5. Keep: uat-gateway, uat-tools
6. Sync all apps
7. Wait for healthy state
8. Deploy L2 contracts
9. ⏱️ Total: ~50-60 minutes

---

## 📋 Common UAT Scenarios

### Scenario 1: Full Service Reset
**When:** Testing all components after major code change
```
testnet_type: uat-testnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.8.0
additional_apps_to_delete: gateway,tools

Result:
  ✅ All child apps deleted and recreated
  ✅ New L1 & L2 contracts deployed
  ✅ Full network reset
  ✅ Clean state for testing
```

### Scenario 2: Core Services Only
**When:** Testing sequencer/validators without resetting gateway/tools
```
testnet_type: uat-testnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.8.0
additional_apps_to_delete: (leave blank)

Result:
  ✅ sequencer, validator-01, validator-02 recreated
  ✅ gateway, tools stay running
  ✅ Faster deployment (~50 min)
  ✅ Gateway/tools can be reused
```

### Scenario 3: New Image Test
**When:** Testing new image build without changing infrastructure
```
testnet_type: uat-testnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.9.0
additional_apps_to_delete: (blank)

Result:
  ✅ New images built and pushed
  ✅ Apps recreated with new images
  ✅ Same infrastructure
  ✅ Tests new build
```

### Scenario 4: Reuse Previous Build
**When:** Re-testing with same images (already built)
```
testnet_type: uat-testnet
deployment_strategy: destructive
image_build: no
image_tag: v1.5.8.0
additional_apps_to_delete: (blank)

Result:
  ✅ No new builds
  ✅ Uses existing images
  ✅ Faster deployment (~35 min)
  ✅ Fresh L1 & L2 contracts
```

### Scenario 5: Config Update
**When:** Testing config changes without reset
```
testnet_type: uat-testnet
deployment_strategy: non-destructive
image_build: no
additional_apps_to_delete: N/A

Result:
  ✅ Apps stay running
  ✅ Configs updated
  ✅ No pod recreation
  ✅ Minimal downtime
  ✅ Fast deployment (~15 min)
```

---

## 📊 Deployment Timeline

### Destructive (Full Reset)
```
Time    Phase                        Duration
─────────────────────────────────────────────
0:00    Start                        -
0:05    Build images                 20-30 min
0:30    Update ten-apps              2 min
0:32    Deploy L1 contracts          10 min
0:42    Delete apps                  2 min
0:44    Sync ArgoCD apps             3 min
0:47    Wait for healthy state       5 min
0:52    Deploy L2 contracts          10 min
1:02    Post-deployment              1 min
─────────────────────────────────────────────
TOTAL: ~1 hour
```

### Destructive (No Build)
```
Time    Phase                        Duration
─────────────────────────────────────────────
0:00    Start                        -
0:01    Validate inputs              1 min
0:02    Deploy L1 contracts          10 min
0:12    Delete apps                  2 min
0:14    Sync ArgoCD apps             3 min
0:17    Wait for healthy state       5 min
0:22    Deploy L2 contracts          10 min
0:32    Post-deployment              1 min
─────────────────────────────────────────────
TOTAL: ~35 minutes
```

---

## ✅ Pre-Deployment Checklist

### Before Running
- [ ] Code merged to main branch
- [ ] Tests passing on PR
- [ ] Images ready to build OR tag specified
- [ ] Team notified of UAT deployment
- [ ] No other UAT deployments running

### Inputs Verified
- [ ] `testnet_type: uat-testnet`
- [ ] `deployment_strategy: destructive` (for full test)
- [ ] `image_build: yes` (if new images) or `no` (if reusing)
- [ ] `image_tag: v1.5.X.X` (if specified)
- [ ] `additional_apps_to_delete: (blank or custom)`

---

## 🔍 Monitoring During Deployment

### GitHub Actions UI
```
Actions → k8s Deploy Consolidated → Running workflow

Check:
  ✅ Job status (should show ▶️ running)
  ✅ Logs for each job
  ✅ Search for ❌ errors
  ✅ Look for ✅ completion messages
```

### Real-Time Monitoring
```bash
# Watch pod creation
kubectl get pods -n uat -w

# Watch app status
argocd app list | grep uat-

# Check app health
argocd app get uat-sequencer --refresh

# View app logs
argocd app logs uat-sequencer --tail 100
```

### Key Milestones to Watch
1. **Build complete** → "✅ Pushed image"
2. **L1 deployed** → "NETWORK_CONFIG_ADDR=0x..."
3. **Apps deleted** → "✅ Child app deletion complete"
4. **Apps syncing** → "🔄 Syncing ArgoCD applications"
5. **Apps healthy** → "✅ Health check complete"
6. **L2 deployed** → "L2 deployer logs"
7. **Done** → "✅ Deployment workflow completed"

---

## 📝 Post-Deployment Verification

### Immediate Checks (5 min after done)
```bash
# Check pods are running
kubectl get pods -n uat -o wide
# Should show: sequencer, validator-01, validator-02 + others

# Check app health
argocd app get uat-sequencer
# Should show: Healthy, Synced

# Check contracts deployed
cd ten-apps
git log --oneline -1
# Should show: "chore: update image tags..."
```

### Network Verification (10 min after done)
```bash
# Connect to sequencer
curl http://uat-sequencer.ten.xyz:80/rpc

# Check block height
cast rpc eth_blockNumber --rpc-url http://uat-sequencer.ten.xyz:80

# Submit test transaction
# Use your test scripts
```

### L2 Contract Verification
```bash
# Check L2 contracts deployed
cd ten-apps
grep -A 10 "l1Config:" nonprod-argocd-config/apps/envs/uat/valuesFile/l1-values.yaml

# Verify in L2
# Run L2 tests against uat-validator-01.ten.xyz
```

---

## 🚨 Troubleshooting UAT Issues

### Issue: Build Failed
```
Error: "Docker build failed"

Solution:
  1. Check images with: az acr repository show-tags --name testnetobscuronet --repository obscuronet/enclave
  2. Verify Docker login: docker login testnetobscuronet.azurecr.io
  3. Check disk space: df -h
  4. Re-run with same image_tag or new image_tag
```

### Issue: L1 Deployment Failed
```
Error: "L1 contract deployer exited with code 1"

Solution:
  1. Check logs: kubectl logs -n uat <pod-name>
  2. Verify gas price: Check max_gas_gwei parameter
  3. Check L1 RPC: curl $L1_HTTP_URL
  4. Re-run with non-destructive first to debug
```

### Issue: ArgoCD Sync Timeout
```
Error: "Failed to sync uat-sequencer (timeout)"

Solution:
  1. Check app status: argocd app get uat-sequencer
  2. Increase timeout: Set sync_timeout to 15m or 20m
  3. Check resources: kubectl describe node
  4. Re-run deployment
```

### Issue: Apps Not Healthy
```
Error: "Apps did not reach healthy state"

Solution:
  1. Check pod status: kubectl get pods -n uat
  2. Describe pod: kubectl describe pod -n uat <pod-name>
  3. Check logs: kubectl logs -n uat <pod-name>
  4. Wait longer: Sometimes pods take time to start
```

### Issue: L2 Deployment Failed
```
Error: "L2 contract deployer exited with code 1"

Solution:
  1. Verify L1 contracts deployed: Check L1 values in YAML
  2. Check L2 RPC: curl http://uat-validator-01.ten.xyz:80
  3. Check L2 logs: kubectl logs -n uat uat-l2-deployer
  4. Re-run just L2: Use non-destructive to skip to L2
```

---

## 📈 Testing Strategy for UAT

### Day 1: Fresh Deployment
```
1. Run destructive with new images
2. Verify all pods running
3. Verify L1 & L2 contracts deployed
4. Run smoke tests
5. Note any issues
```

### Day 2: Resync Test
```
1. Run non-destructive to verify sync works
2. Update some config in ten-apps
3. Verify ArgoCD picks up changes
4. Run tests again
```

### Day 3: Full Reset Test
```
1. Run destructive again with additional apps
2. Delete gateway + tools too
3. Verify full network recovery
4. Run comprehensive tests
```

### Day 4: Performance Test
```
1. Run same version again (no build)
2. Measure deployment time
3. Measure sync time
4. Measure pod startup time
5. Document for baseline
```

---

## 🎯 UAT Approval-Free Advantages

✅ **No approval needed** - Run anytime
✅ **Full destructive available** - True testing
✅ **Rapid iteration** - Quick feedback loop
✅ **Safe failures** - No production impact
✅ **Testing flexibility** - Try any combination

---

## 📊 Expected Metrics

### Timing
- Build images: 20-30 min
- Deploy L1: 10 min
- Deploy L2: 10 min
- App sync: 3 min
- Total: 50-60 min (with build), 35 min (no build)

### Resource Usage
- Disk: ~2-3 GB per build
- Memory: ~4-6 GB during sync
- CPU: ~2-4 cores during build

### Network State
- Blocks: Starting from 0 after L1 reset
- L2 blocks: Should produce blocks quickly
- Gas: L1_CHALLENGE_PERIOD = 1 week (uat values)

---

## 🔗 Related Commands

### Useful kubectl Commands
```bash
# Get pods
kubectl get pods -n uat

# Describe pod
kubectl describe pod -n uat <pod-name>

# Check logs
kubectl logs -n uat <pod-name> -c <container>

# Port forward
kubectl port-forward -n uat <pod-name> 8080:80

# Delete pod (to restart)
kubectl delete pod -n uat <pod-name>

# Get node info
kubectl get nodes -L kubernetes.io/hostname
```

### Useful ArgoCD Commands
```bash
# List apps
argocd app list | grep uat

# Get app status
argocd app get uat-sequencer

# Refresh app
argocd app get uat-sequencer --refresh

# View logs
argocd app logs uat-sequencer --tail 50

# Manual sync
argocd app sync uat-sequencer --wait

# Delete app
argocd app delete uat-sequencer --yes
```

### Useful GitHub Actions Commands
```bash
# Watch workflow
gh run watch <run-id>

# View logs
gh run view <run-id>

# List recent runs
gh run list --workflow=manual-deploy-k8s-consolidated.yml
```

---

## 💾 Saved States

### Post-Successful Deployment
```bash
# Save git log
cd ten-apps
git log --oneline -5 > /tmp/uat_deployment.log

# Save pod status
kubectl get pods -n uat -o wide > /tmp/uat_pods.txt

# Save app status
argocd app list | grep uat > /tmp/uat_apps.txt

# Reference for comparison
```

---

## 🎓 Learning Resources

**For first-time UAT tester:**
1. Read this guide (15 min)
2. Run non-destructive on dev first (30 min)
3. Run destructive on dev (60 min)
4. Run destructive on UAT (60 min)
5. Debug issues together (varies)

**Total first-time investment:** ~3 hours

---

## 📞 Quick Help

### Who to Ask
- **Deployment issues**: DevOps team
- **Test failures**: QA team
- **Contract issues**: Smart contract team
- **Network issues**: Network team

### Where to Check
- **Workflow logs**: GitHub Actions
- **Pod logs**: `kubectl logs`
- **App status**: `argocd app get`
- **Recent git**: `cd ten-apps && git log`

---

## 🎉 Success Criteria

✅ **Deployment Complete** when:
- GitHub workflow shows ✅ (all green)
- All pods running: `kubectl get pods -n uat`
- ArgoCD shows Healthy: `argocd app get uat-sequencer`
- L1 contracts deployed: Check YAML
- L2 contracts deployed: Check logs
- Network responding: curl/cast working

🎯 **Ready for next phase** when:
- 3+ destructive deployments succeeded
- All pods stable for 30+ min
- L1 & L2 contracts verified
- Network tests passing
- Ready to move to sepolia

---

## 🚀 Next Steps After UAT Success

1. ✅ Document any issues found
2. ✅ Update configuration if needed
3. ✅ Prepare for sepolia deployment
4. ✅ Brief team on findings
5. ✅ Schedule sepolia test with approval

---

**Good luck with UAT testing! 🎯**

*For general deployment info, see DEPLOYMENT_GUIDE.md*
*For troubleshooting, see QUICK_REFERENCE.md*
*For image versioning, see IMAGE_TAGGING_GUIDE.md*
