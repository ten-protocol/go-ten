# Complete Rollback Procedures & Emergency Response

## Overview

This document provides comprehensive rollback procedures and emergency response protocols for the Ten Protocol's contract upgrade system. It addresses the critical scenario of malicious upgrades and provides step-by-step emergency procedures.

## How Major Protocols Handle Emergency Rollbacks

Unisqap, Compound, Aave, MakerDAO, and Curve all use the same approach.

**Universal Pattern:**
1. **Emergency Pause** - All protocols have immediate pause functionality
2. **No Emergency Rollback** - All rollbacks go through normal governance timelock
3. **24-48 Hour Delays** - Accepted even for emergency rollbacks
4. **Security Over Speed** - All protocols prioritize security over speed

### Why No Emergency Rollback?

**Security Design:**
- **Timelock delays** prevent rushed decisions
- **Governance process** ensures community consensus
- **No single point of failure** for rollbacks

**Attack Prevention:**
- **Prevents malicious rollbacks** by compromised governance
- **Maintains system integrity** even during emergencies
- **Protects against social engineering** attacks

**Community Trust:**
- **Transparent process** - all rollbacks are visible
- **Community consensus** - requires governance approval
- **Audit trail** - all actions are recorded

### Industry Standard Approach

**Two-Phase Emergency Response:**
1. **Phase 1**: Emergency pause (immediate) - Stop malicious code
2. **Phase 2**: Normal governance rollback (24-hour delay) - Restore security

## High Level Overview

### The Problem
If a hacker gains access to the multisig and schedules a malicious upgrade that could drain funds, here's what happens and how to respond:

### Emergency Response Process
1. **Detect Malicious Upgrade** - During 24-hour timelock delay
2. **Emergency Pause** - Immediately pause all contracts (0-1 hour)
3. **Investigate & Prepare** - Analyze threat and prepare rollback (1-24 hours)
4. **Normal Rollback** - Execute rollback through timelock (24+ hours)
5. **Verify & Resume** - Verify security and unpause contracts

### Key Commands
```bash
# Emergency pause all contracts
export NETWORK_CONFIG_ADDR="0x1111111111111111111111111111111111111111"
export EMERGENCY_PAUSER_ADDRESS="0x1234567890123456789012345678901234567890"
npx hardhat run scripts/emergency/001_emergency_pause.ts --network mainnet

# Rollback to secure implementation (24-hour delay)
export TIMELOCK_ADDRESS="0x9876543210987654321098765432109876543210"
npx hardhat run scripts/upgrade/001_upgrade_contracts.ts --network mainnet

# Unpause contracts after rollback
export EMERGENCY_ACTION="unpause"
npx hardhat run scripts/emergency/001_emergency_pause.ts --network mainnet
```

### Industry Standard Approach
- **Emergency Pause** - Immediate response (industry standard)
- **Normal Rollback** - 24-hour delay for security (industry standard)
- **No Emergency Rollback** - Cannot bypass timelock (industry standard)
- **Security Over Speed** - Prioritizes system integrity (industry standard)

---

## Detailed Implementation

## Critical Scenario: Malicious Upgrade That Could Drain Funds

### The Problem
If a hacker gains access to the multisig and schedules a malicious upgrade that could drain funds, here's what happens and how to respond:

### What Happens Currently (Without Enhanced Rollback)

#### Timeline of a Malicious Upgrade:
```
T+0h:   Hacker gains multisig access → schedules malicious upgrade
T+24h:  Malicious upgrade executes → funds can be drained
T+24h+: System compromised → potential fund loss
```

#### Current Limitations:
- No cancellation - Cannot stop upgrade once scheduled
- 24-hour delay - Gives attacker time to prepare
- No emergency pause - Contracts continue operating
- No immediate rollback - Must wait for execution

## Emergency Response Procedures

### Phase 1: Immediate Response (0-1 hour)

#### 1. Alert Team
- **Notify all stakeholders** immediately
- **Activate emergency procedures** - Follow emergency plan
- **Assess threat level** - Determine potential damage

#### 2. Emergency Pause (Immediate)
```bash
# Pause all contracts immediately to stop malicious code execution
export NETWORK_CONFIG_ADDR="0x1111111111111111111111111111111111111111"
export EMERGENCY_PAUSER_ADDRESS="0x1234567890123456789012345678901234567890"

# Pause all contracts
npx hardhat run scripts/emergency/001_emergency_pause.ts --network mainnet
```

**What this does:**
- Immediately pauses all contracts
- Stops malicious code execution
- Prevents fund drainage - No transactions can be processed
- Bypasses timelock - No delay for emergency response

### Phase 2: Investigation & Preparation (1-24 hours)

#### 3. Investigate the Threat
- **Analyze upgrade code** - Determine what malicious code does
- **Assess potential damage** - Calculate potential fund loss
- **Prepare counter-measures** - Develop secure rollback

#### 4. Prepare Secure Rollback
- **Deploy secure implementation** - Clean version without malicious code
- **Test thoroughly** - Ensure rollback works correctly
- **Prepare execution** - Ready to rollback immediately

### Phase 3: Rollback Execution (24+ hours)

#### 5. Normal Timelock Rollback (No Emergency Bypass)
```bash
# Rollback to secure implementation using normal timelock process
export NETWORK_CONFIG_ADDR="0x1111111111111111111111111111111111111111"
export EMERGENCY_PAUSER_ADDRESS="0x1234567890123456789012345678901234567890"
export TIMELOCK_ADDRESS="0x9876543210987654321098765432109876543210"

# Rollback all contracts (uses normal timelock - 24-hour delay)
npx hardhat run scripts/upgrade/001_upgrade_contracts.ts --network mainnet
```

**What this does:**
- Deploys secure implementation - New clean version
- Schedules rollback through timelock - Normal process
- Waits for 24-hour delay - Security enforced
- Executes rollback after delay - System restored
- **Secure:** Uses normal timelock process, no bypass possible

**Important:** Emergency rollback is **not possible** because the timelock enforces a minimum delay (24 hours) that cannot be bypassed. This is a **security feature**.

#### 6. Verify Rollback
```bash
# Verify rollback was successful
export NETWORK_CONFIG_ADDR="0x1111111111111111111111111111111111111111"

# Verify all rollbacks
export EMERGENCY_ACTION="verify"
npx hardhat run scripts/emergency/002_emergency_rollback.ts --network mainnet
```

#### 7. Unpause Contracts
```bash
# Unpause all contracts after rollback
export NETWORK_CONFIG_ADDR="0x1111111111111111111111111111111111111111"
export EMERGENCY_PAUSER_ADDRESS="0x1234567890123456789012345678901234567890"

# Unpause contracts
export EMERGENCY_ACTION="unpause"
npx hardhat run scripts/emergency/001_emergency_pause.ts --network mainnet
```

## Enhanced Security Features

### Emergency Pause Functionality


#### Integrate into Existing Upgradeable Contracts
```solidity
// Add to existing upgradeable contracts
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";

contract UpgradeableContract is Pausable, Initializable, AccessControlUpgradeable {
    bytes32 public constant EMERGENCY_PAUSER_ROLE = keccak256("EMERGENCY_PAUSER_ROLE");
    
    function initialize(address _emergencyPauser) public initializer {
        __AccessControl_init();
        __Pausable_init();
        
        // Grant emergency pauser role
        _grantRole(EMERGENCY_PAUSER_ROLE, _emergencyPauser);
        
        // Grant DEFAULT_ADMIN_ROLE to the same address for role management
        _grantRole(DEFAULT_ADMIN_ROLE, _emergencyPauser);
    }
    
    function emergencyPause() external onlyRole(EMERGENCY_PAUSER_ROLE) {
        _pause();
        emit EmergencyPaused(msg.sender);
    }
    
    function emergencyUnpause() external onlyRole(EMERGENCY_PAUSER_ROLE) {
        _unpause();
        emit EmergencyUnpaused(msg.sender);
    }
    
    function grantEmergencyPauserRole(address newPauser) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newPauser != address(0), "Invalid pauser address");
        _grantRole(EMERGENCY_PAUSER_ROLE, newPauser);
    }
    
    function revokeEmergencyPauserRole(address pauser) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(pauser != address(0), "Invalid pauser address");
        _revokeRole(EMERGENCY_PAUSER_ROLE, pauser);
    }
    
    modifier whenNotPaused() override {
        require(!paused(), "Contract is paused");
        _;
    }
}
```

### Benefits:
- Immediate response - Can pause contracts instantly
- No timelock delay - Bypasses 24-hour delay
- Prevents damage - Stops malicious code execution
- Simple implementation - Easy to add to contracts
- **Proper access control** - Role-based permissions instead of single address
- **Flexible role management** - Can add/remove emergency pausers as needed

### Trade-offs:
- Centralized control - Single point of failure
- Trust requirement - Must trust emergency pauser
- Potential abuse - Could be used maliciously

## Implementation Steps

### 1. Add Emergency Pause to Existing Contracts

```solidity
contract CrossChain is ICrossChain, Initializable, UnrenouncableOwnable2Step, ReentrancyGuardUpgradeable, PausableUpgradeable, AccessControlUpgradeable {
    bytes32 public constant EMERGENCY_PAUSER_ROLE = keccak256("EMERGENCY_PAUSER_ROLE");
    
    // ... existing mappings and variables ...
    
    function initialize(address owner, address _messageBus, address _emergencyPauser) public initializer {
        require(_messageBus != address(0), "Invalid message bus address");
        require(owner != address(0), "Owner cannot be 0x0");
        require(_emergencyPauser != address(0), "Emergency pauser cannot be 0x0");
        
        __UnrenouncableOwnable2Step_init(owner);
        __ReentrancyGuard_init();
        __Pausable_init();
        __AccessControl_init();
        
        merkleMessageBus = MerkleTreeMessageBus.IMerkleTreeMessageBus(_messageBus);
        messageBus = MessageBus.IMessageBus(_messageBus);
        
        // Grant emergency pauser role
        _grantRole(EMERGENCY_PAUSER_ROLE, _emergencyPauser);
        _grantRole(DEFAULT_ADMIN_ROLE, _emergencyPauser);
    }
    
    function emergencyPause() external onlyRole(EMERGENCY_PAUSER_ROLE) {
        _pause();
        emit EmergencyPaused(msg.sender);
    }
    
    function emergencyUnpause() external onlyRole(EMERGENCY_PAUSER_ROLE) {
        _unpause();
        emit EmergencyUnpaused(msg.sender);
    }
    
    function grantEmergencyPauserRole(address newPauser) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(newPauser != address(0), "Invalid pauser address");
        _grantRole(EMERGENCY_PAUSER_ROLE, newPauser);
    }
    
    function revokeEmergencyPauserRole(address pauser) external onlyRole(DEFAULT_ADMIN_ROLE) {
        require(pauser != address(0), "Invalid pauser address");
        _revokeRole(EMERGENCY_PAUSER_ROLE, pauser);
    }
    
    // Add pause modifier to critical functions
    modifier whenNotPaused() override {
        require(!paused(), "Contract is paused");
        _;
    }
    
    // Example: Add pause check to critical functions
    function someCriticalFunction() external whenNotPaused {
        // ... existing logic ...
    }
}
```

**Key points:**
- **Use AccessControlUpgradeable** - Leverage OpenZeppelin's role-based access control
- **Define EMERGENCY_PAUSER_ROLE** - Specific role for emergency pause functionality
- **Proper initialization** - Initialize all parent contracts (AccessControl, Pausable)
- **Role management functions** - Allow adding/removing emergency pausers
- **Override whenNotPaused** - Ensure pause functionality works correctly
- **Add pause checks** - Apply `whenNotPaused` modifier to critical functions

### 2. Create Emergency Response Script
```typescript
// scripts/emergency/001_emergency_pause.ts
async function emergencyPause() {
    const [emergencyPauser] = await ethers.getSigners();
    
    // Pause all contracts
    const contracts = [
        "CrossChain",
        "NetworkEnclaveRegistry", 
        "DataAvailabilityRegistry"
    ];
    
    for (const contractName of contracts) {
        const contract = await ethers.getContractAt(contractName, address);
        await contract.emergencyPause();
        console.log(`${contractName} paused`);
    }
}
```

### 3. Create Rollback Script
```typescript
// scripts/emergency/002_emergency_rollback.ts
async function emergencyRollback() {
    // Deploy secure implementation
    const secureImpl = await deploySecureImplementation();
    
    // Upgrade all contracts immediately
    const contracts = [
        "CrossChain",
        "NetworkEnclaveRegistry",
        "DataAvailabilityRegistry"
    ];
    
    for (const contractName of contracts) {
        await upgradeContract(contractName, secureImpl);
        console.log(`${contractName} rolled back`);
    }
}
```