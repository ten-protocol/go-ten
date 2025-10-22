# Consolidated K8s Deployment Workflow

## 📋 Overview

This directory contains the **consolidated Kubernetes deployment workflow** that replaces the previous 3-step manual process.

**Old workflow:**
- `manual-deploy-k8s-testnet-before-nodes.yml` (L1 contracts)
- `manual-deploy-k8s-testnet-after-nodes.yml` (L2 contracts + ArgoCD)
- `build-release-images.yml` (image building)

**New workflow:**
- `manual-deploy-k8s-consolidated.yml` (everything in one)

---

## 🚀 Quick Start

### Run a Deployment

1. Go to GitHub → Actions → "k8s Deploy Consolidated"
2. Click "Run workflow"
3. Fill in the form:
   - `testnet_type`: dev-testnet, uat-testnet, sepolia-testnet, or mainnet
   - `deployment_strategy`: non-destructive or destructive
   - `image_build`: yes or no
   - `image_tag`: (leave blank for auto-detect)
   - Other fields: optional
4. Click "Run workflow" button
5. Monitor in the Actions tab

### Approval Required?
- Destructive + sepolia-testnet → Approval needed
- Destructive + mainnet → Approval needed
- Everything else → No approval

---

## 📚 Documentation Files

Read these in order:

### 1. **QUICK_REFERENCE.md** ⭐ START HERE
- Most common operations
- Decision tree
- Troubleshooting quick fixes
- Common patterns
- ~5 minute read

### 2. **DEPLOYMENT_GUIDE.md**
- Complete input parameter reference
- Detailed usage examples
- Approval gate details
- Environment configuration
- Health check info
- Troubleshooting checklist

### 3. **CONSOLIDATION_SUMMARY.md**
- What changed from old to new
- Key features
- Job execution flow
- Timeline estimates
- FAQ

### 4. **IMAGE_TAGGING_GUIDE.md**
- Image naming strategy
- Git tag vs commit hash
- How versions are determined
- Best practices
- Registry management

### 5. **NODE_SELECTOR_GUIDE.md** (Destructive Only)
- Node selector overview (SGX constraints)
- When node selectors apply
- How to specify target nodes
- Monitoring pod rescheduling
- Troubleshooting

### 6. **IMPLEMENTATION_CHECKLIST.md**
- Pre-deployment setup
- Testing phases
- Post-deployment verification
- Cleanup steps
- Sign-off template

---

## 🎯 Workflow Features

### Deployment Strategies

**Non-Destructive:**
- ✅ Build images (optional)
- ✅ Deploy L1 contracts
- ✅ Sync ArgoCD apps (incremental)
- ❌ NO L2 deployment
- ❌ Pods stay on current nodes
- **Use:** Config/image updates

**Destructive:**
- ✅ Build images (optional)
- ✅ Deploy L1 contracts
- ✅ Delete child ArgoCD apps
- ✅ Sync ArgoCD apps (fresh)
- ✅ Deploy L2 contracts
- ✅ Update node selectors (if provided)
- **Use:** Network reset, maintenance

### Automatic Features

- 🖼️ **Image building**: Optional, with git tag or commit hash detection
- 📝 **YAML updates**: Auto-update ten-apps with new image tags/node selectors
- 📤 **Git commits**: Auto-commit changes to ten-apps repo
- 🔒 **Approval gates**: Automatic for destructive prod deployments
- ⏳ **Health checks**: Wait for ArgoCD apps to reach Healthy state
- 🔔 **Notifications**: Trigger ten-test repository dispatch (dev/uat)

---

## 📊 Input Parameters

| Input | Type | Required | Default | Notes |
|-------|------|----------|---------|-------|
| `testnet_type` | choice | ✅ | - | dev, uat, sepolia, mainnet |
| `deployment_strategy` | choice | ✅ | - | non-destructive or destructive |
| `image_build` | choice | ✅ | - | yes or no |
| `image_tag` | string | ❌ | auto-detect | git tag or commit hash |
| `confirmation` | string | ❌ | - | Type "confirm" for prod destructive |
| `log_level` | number | ❌ | 3 | 1-5 (Error to Trace) |
| `max_gas_gwei` | string | ❌ | 1.5 | Gas price for L1 |
| `sync_timeout` | choice | ❌ | 10m | ArgoCD sync timeout |
| `sequencer_node_selector` | string | ❌ | - | DESTRUCTIVE ONLY |
| `validator_01_node_selector` | string | ❌ | - | DESTRUCTIVE ONLY |
| `validator_02_node_selector` | string | ❌ | - | DESTRUCTIVE ONLY |
| `gateway_node_selector` | string | ❌ | - | DESTRUCTIVE ONLY |

---

## 🏗️ Job Flow

```
validate-inputs
    ↓
[approval] (if destructive + prod)
    ↓
build-images (if image_build: yes)
    ↓
update-ten-apps-config (update YAML + commit)
    ↓
deploy-l1-contracts
    ↓
[argocd-delete-child-apps] (if destructive only)
    ↓
argocd-sync-apps
    ↓
[wait-argocd-healthy] (if destructive)
    ↓
[deploy-l2-contracts] (if destructive only)
    ↓
post-deployment
```

---

## ⏱️ Timing

| Phase | Time |
|-------|------|
| Validate | 1 min |
| Build images | 20-30 min |
| Update ten-apps | 2 min |
| Deploy L1 | 10 min |
| Delete apps | 2 min |
| Sync ArgoCD | 3 min |
| Wait healthy | 5 min |
| Deploy L2 | 10 min |
| **Total (destructive)** | **~50-60 min** |
| **Total (non-destructive)** | **~35-45 min** |

---

## ✅ Approval Process

**When approval is triggered:**
1. Workflow pauses at `approval` job
2. GitHub sends notification to configured reviewers
3. Reviewer goes to Actions → Workflow run → Review deployments
4. Reviewer clicks "Approve and deploy"
5. Workflow continues automatically

**Who needs to approve:**
- Set in GitHub Environment settings (Settings → Environments → [env] → Required reviewers)

**Environments requiring approval:**
- `sepolia-testnet` (destructive only)
- `mainnet` (destructive only)

---

## 🔍 Monitoring Deployments

### In GitHub Actions
- Go to Actions → "k8s Deploy Consolidated"
- Click running workflow
- Expand each job to see logs
- Search logs for ✅/❌ symbols

### Key Log Markers
- ✅ Job succeeded
- ❌ Job failed
- 📦 Image operation
- 📝 YAML update
- 🔄 ArgoCD sync
- ⏳ Waiting for condition
- 🔒 Approval required

### Post-Deployment Verification
```bash
# Check ArgoCD apps
argocd app get dev-sequencer --refresh

# Check deployed pods
kubectl get pods -n dev -o wide

# Check image tags in ten-apps
cd ten-apps
git log --oneline -5
git show HEAD
```

---

## 🐛 Troubleshooting

### Workflow Fails at `build-images`
- Check Docker build logs
- Verify Azure credentials valid
- Check disk space on runner
- See BUILD_ERRORS in logs

### Workflow Fails at `argocd-sync-apps`
- Check ArgoCD server is reachable
- Verify ARGOCD_TOKEN is valid
- Check app health: `argocd app get <app-name>`
- Increase `sync_timeout` if needed

### Workflow Fails at `deploy-l2-contracts`
- Only runs for destructive deployments
- Check node is healthy: `kubectl get nodes`
- Check contract deployer logs
- Review L2 deployment error

### Ten-apps Not Updated
- Check DEPLOY_ACTIONS_PAT has write access
- Verify git commit in ten-apps repo
- Check workflow logs for "Updated and pushed"

### Approval Not Appearing
- Verify environment has "Required reviewers"
- Check you selected "destructive" strategy
- Confirm targeting sepolia or mainnet
- Verify "confirmation: confirm" field filled

---

## 🚨 Common Mistakes

❌ **Leave approval fields blank for production destructive**
- Always type "confirm" in confirmation field

❌ **Mix commit hashes and git tags**
- Be consistent with versioning

❌ **Run destructive without testing non-destructive first**
- Always test on dev first

❌ **Forget to push git tags**
- `git push origin v1.5.8.0` is required

❌ **Expect L2 deployment on non-destructive**
- L2 only runs for destructive

---

## 📖 Example Workflows

### Example 1: Update Images on Dev
```
testnet_type: dev-testnet
deployment_strategy: non-destructive
image_build: yes
image_tag: v1.5.8.0

Result: Images built → ten-apps updated → L1 deployed → apps synced (NO L2 reset)
```

### Example 2: Full Dev Reset
```
testnet_type: dev-testnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.8.0

Result: Images built → apps deleted → recreated → L2 deployed (FULL RESET)
```

### Example 3: UAT Deploy (No Build)
```
testnet_type: uat-testnet
deployment_strategy: non-destructive
image_build: no
image_tag: v1.5.8.0

Result: Use existing image → L1 deployed → apps synced
```

### Example 4: Mainnet Destructive (Approval)
```
testnet_type: mainnet
deployment_strategy: destructive
image_build: yes
image_tag: v1.5.8.0
confirmation: confirm

Result: ⏸️ APPROVAL REQUIRED → After approval → Full deployment
```

---

## 🔗 Related Workflows

**Other deployment workflows:**
- `manual-deploy-obscuro-gateway.yml` - Deploy gateway only
- `manual-deploy-testnet-l1.yml` - Deploy L1 only
- `manual-decommission-testnet-*.yml` - Clean up environments

**Other CI/CD workflows:**
- `build-pr.yml` - Build on PR
- `on-merged-pr.yml` - Build on merge
- `build-release-images.yml` - Release image building

---

## 📞 Support

### Self-Help First
1. Read `QUICK_REFERENCE.md`
2. Check workflow logs
3. Review troubleshooting section
4. Search GitHub issues

### Getting Help
- **Issues**: GitHub Issues in this repo
- **Questions**: Team Slack channel
- **Production incidents**: Page on-call

### Providing Feedback
- Report bugs: GitHub Issues
- Feature requests: Discussion
- Documentation feedback: Pull requests

---

## 🔄 Workflow File

**Main workflow file:**
- `manual-deploy-k8s-consolidated.yml` (~850 lines)

**Structure:**
```
Workflow name: [M] k8s Deploy Consolidated

Inputs:
  - testnet_type (required)
  - deployment_strategy (required)
  - image_build (required)
  - image_tag (optional)
  - confirmation (optional)
  - log_level (optional)
  - max_gas_gwei (optional)
  - sync_timeout (optional)
  - sequencer_node_selector (optional)
  - validator_01_node_selector (optional)
  - validator_02_node_selector (optional)
  - gateway_node_selector (optional)

Jobs:
  - validate-inputs
  - approval (conditional)
  - build-images (conditional)
  - update-ten-apps-config (conditional)
  - deploy-l1-contracts
  - argocd-delete-child-apps (conditional)
  - argocd-sync-apps
  - wait-argocd-healthy (conditional)
  - deploy-l2-contracts (conditional)
  - post-deployment
```

---

## ✨ Key Improvements Over Old Workflow

| Feature | Old | New |
|---------|-----|-----|
| **Manual steps** | 2-3 ArgoCD UI clicks | 0 (fully automated) |
| **Single workflow** | 3 separate workflows | 1 unified |
| **Approval gates** | Manual | Automatic |
| **Image tagging** | Manual + `build-release-images.yml` | Integrated |
| **Audit trail** | Split across 3 runs | Single run |
| **Rollback** | Complex | Simple re-run |
| **Config updates** | Manual + sync | Auto-commit + sync |
| **Health checks** | Manual verification | Automated polling |
| **Node selectors** | Manual in ArgoCD | UI inputs |

---

## 📝 Migration from Old Workflows

**Timeline:**
- Week 1: Test on dev
- Week 2: Test on uat
- Week 3: Cutover to production
- Week 4: Archive old workflows

**Steps:**
1. Read all documentation
2. Test non-destructive on dev
3. Test destructive on dev
4. Run through approval process on uat
5. Archive old workflow files
6. Update team runbooks
7. Announce to team

---

## 📌 Important Notes

⚠️ **Destructive deployments:**
- Use only when you need full network reset
- Always test on dev first
- Requires approval for production
- Plan for 50-60 minute downtime

⚠️ **Node selectors (destructive only):**
- SGX node constraints
- Only specify nodes you want to move
- Leave blank to keep current assignment
- Changes auto-committed to ten-apps

⚠️ **Image versioning:**
- Use git tags for releases
- Use commit hash for dev builds
- Consistent versioning across all images
- Auto-detected from git if not specified

⚠️ **Approval gates:**
- Destructive + sepolia → Approval required
- Destructive + mainnet → Approval required
- Must type "confirm" in confirmation field
- Approval must be from configured reviewer

---

## 🎓 Learning Path

**First time?**
1. Read `QUICK_REFERENCE.md` (5 min)
2. Watch workflow run on dev (30 min)
3. Read `DEPLOYMENT_GUIDE.md` (15 min)

**Deploying to uat?**
1. Read `CONSOLIDATION_SUMMARY.md` (10 min)
2. Run non-destructive on uat (30 min)

**Deploying to production?**
1. Read `IMPLEMENTATION_CHECKLIST.md` (20 min)
2. Coordinate with team (30 min)
3. Run approved deployment (60 min)
4. Monitor and verify (30 min)

---

**Happy deploying! 🚀**

*For questions or issues, see the Support section above.*
