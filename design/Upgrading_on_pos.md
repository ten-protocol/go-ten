# How the transition of Ethereum to "Proof of Stake" affects Obscuro?

There are subtle but meaningful differences between the proof-of-work and proof-of-stake protocols.
The most significant difference is relevant when you consider the upgradeability of a privacy L2. 

## Upgradability of Obscuro

A TEE-based system like Obscuro must be upgradeable.
The relevant requirement is that data is private for anyone, including the software developers and the infrastructure operators. 

During an upgrade, the new software version must be able to pick up from the old version, which means that the previous version must share some of its secrets with the new version.

### Requirement for the upgrading process 

We'll start with the requirement and then show that the process is flawed and can be attacked if these conditions are not met.
 
Requirement: Every release must have a built-in decentralized and transparent mechanism for accepting future versions of itself as valid and, more importantly, for rejecting malicious versions.

If this mechanism is opaque or centralized, whoever controls it can unilaterally upgrade the software to a malicious version and gain access to the secrets. The possibility of such an attack performed behind closed doors can severely undermine the guarantees of a privacy-preserving solution.

Note: This is a different class of attacks from the ones against the typical transparent chains.

### Upgrading based on a single signer

The most naive upgrading mechanism requires the same private key to sign the new version as the key used to sign the current one. 
This straightforward method can be seen as the same developer being trusted with upgrading.
The problem with this approach is that, without additional protections, it gives the developer the power to retrieve the secrets unchallenged.
They can create a new malicious version, upgrade a node decoupled from the internet to this version, and access the secrets.

This mechanism is simple but opaque and centralized.

### Upgrading based on multiple signers

An improved option is a Multi-Sig approach, where a committee of independent parties must sign over each upgrade. 
The developer proposes a new version and submits it for auditing to the members, who, upon careful analysis, will sign over it if they think it is valid and protects their interests.

This could be feasible if the committee members have a good reputation and the right incentives. But unfortunately, it is challenging to find such reputable parties, especially for a young network.
If there is suspicion of collusion between them, the perception of security can be permanently compromised. Also, operationally, it makes upgrades more difficult since every committee member must review and sign the code during a period.


### Upgrade lifecycle managed by a smart contract

Another option is to leverage the Ethereum Layer 1 for transparency and decentralization by entrusting the task of managing upgrades to a smart contract. The TEE would accept a new software version as a valid upgrade if it went through a predetermined lifecycle on Ethereum enforced by the smart contract.

#### The lifecycle of a Smart-contract managed upgrade

1. A developer proposes a new version by publishing a GitHub tag (hash) and the Attestation, together with a stake.
Note: Initially, the developer will likely be the Obscuro Foundation.

2. Anyone posting a stake can make a challenge to this version. The challenge period is predetermined.

3. The original developer first evaluates the challenge. If it is valid, then a reduced stake is paid out to the challenger as a bug bounty.

4. If the developer rejects the challenge, the challenger can request a vote from the community.  
Todo: establish an incentive for the community to vote.

This mechanism has several advantages:
- It is fully transparent. Everyone can inspect the code before it gets rolled out. And the history of all upgrades is recorded on a transparent, immutable chain.
- The code must be published and go through a transparent lifecycle before it receives any secrets. This removes the power of a developer or committee to gain a stealth advantage. 
- It incentivizes security researchers to find issues in the codebase.
- It is as decentralized as the smart contract, which can be easily audited.

## Weak subjectivity of Ethereum 2

Because of the many advantages, Obscuro decided to implement the transparent, smart-contract-driven upgrade process. 

The design needs two key ingredients:
1. The enclave must understand the outputs of the smart contract.
2. The upgrader must convince the enclave that the data it feeds is the actual Etherem canonical chain. 

The first part is straightforward because the enclave can interpret and authenticate events emitted on Ethereum against a chain. This functionality is already required for cross-chain messaging and authenticating rollups.

The second point is difficult to achieve in Ethereum 2 because of 'weak subjectivity'. We'll see how, in the context of TEEs, the "weak" turns into "strong" subjectivity.

The threat is that an attacker can bypass the transparency guarantees of the upgrading mechanism by producing a valid Ethereum fork, different from the canonical chain, that contains a malicious version of the enclave. 
The cryptoeconomic guarantees of Proof-of-Stake don't apply because the attacker will not broadcast the fork to the outside world. 

Note: a naive solution requires the validated versions to be signed by a committee, which brings us back to the previously discussed problems.

The threat against an enclave differs from the threat against an Ethereum node because an enclave has no trusted friends to feed it with a "Weak Subjectivity Checkpoint." 
The assumption is that Ethereum nodes are run by operators who are incentivized to run against a canonical fork.
In this attack, the operator controls the data flow into the enclave, including the clock, and has the incentive to run against a malicious fork.

The enclave needs a reliable mechanism to determine whether the chain it is being fed is canonical.

For proof-of-work chains, we can use a simple heuristic. The enclave can calculate the amount of "work" that went into mining the blocks it was fed. If the total work decreases below a threshold, it will reject future blocks.

This document will not go into the details of the Ethereum proof-of-stake protocol because it is very complex. Instead, we'll present some possible practical attacks that illustrate the threat and propose heuristics designed to mitigate them.

### Weak subjectivity attacks against a TEE

Assumptions:
The current version of the enclave starts with a "weak subjectivity checkpoint" as part of its Attestation, which is the same as any beacon chain client.
Note: this first assumption leads to an interesting recursivity which we'll analyze below.

The attackers are part of the original validator set but only control a low amount of the total stake. 
Note: an attacker controlling a super-majority of the total Ethereum stake will be able to break the mechanism. 

The attackers aim to convince the enclave to upgrade itself to a malicious version without public scrutiny.

The enclave can execute the proof-of-stake protocol and asses the finality of blocks.

During any upgrade, the current enclave will share some secrets with the future enclave. 

#### Inactivity leak attack

The attackers want to exploit the "inactivity leak". Given that they are already registered validators, the simplest attack is for these validators to attempt to create a chain where everyone else has suddenly dropped off.

They will produce blocks when the protocol decides it is their turn. And they will be the only block producers. Since the stake they control is very low, the checkpoints will not finalize. After four epochs, the "inactivity" rule will kick in, which will chip away at the stake of all the inactive validators (basically everyone except the attackers).
After many blocks, the attackers will control 2/3rd of the entire stake and will be able to finalize checkpoints.

Note that the blocks produced during this attack are not broadcast, so the validators will not get penalized in the actual canonical chain for signing conflicting attestations.

As a result of this attack, the attacker can feed finalized blocks into the enclave. These finalized blocks will contain upgrade transactions that appear valid to the enclave but violate the transparency assumption, which was the mechanism's purpose. 

Without additional heuristics, the enclave will receive a "Valid Upgrade" event by the "upgrade" smart contract and will assume that it can safely release secrets to that version.

#### Long-range attack

A variation of the previous attack is the typical "long-range" attack from the Ethereum literature.
Validators withdraw their stake and sell their keys on the black market, and the enclave attackers buy them.

To prepare, the attackers will keep an enclave dormant until they control a super-majority of the stake at the time when the enclave goes to sleep.

As soon as they control the super-majority, they wake the enclave up and start feeding malicious finalized blocks with the same result as above.

Note 1: Another attack can combine the above techniques to reduce the number of unfinalized blocks.

Note 2: Both these attacks will reveal past data. They cannot read future data if the upgrading changes the secrets.

### Mitigations

Given these attacks' "subjective" nature, we are unlikely to find a definitive "objective" solution.
Instead, we can aim to introduce heuristics to make it improbable or to increase the capital required to impractical amounts.
The enclave will stop working if the heuristics start looking suspicious, thus preventing the attackers from achieving their goal. 
The risk of choosing too strict heuristics is that the network can become nonfunctional.

#### Reliable source of time - preventing fast-forwarding

A reliable time source inside the enclave can easily prevent the inactivity leak attack because it requires a fast-forward of time until the attackers dominate the total stake.

A relatively easy way to introduce a decentralized, trustless, cryptographic source of time inside the enclave is to feed it the Bitcoin block headers.
The Bitcoin blocks can be used as an approximate but reliable time source by applying the proof-of-work heuristic mentioned above.
With this mechanism in place, the attacker can slow down time, but they can only fast-forward it by spending a lot of money. 

The process itself is simple. The enclave will run the light-chain bitcoin code, resulting in a database entry with the approximate time.
If the timestamp of the Ethereum chain diverges significantly from the Bitcoin chain, the enclave will give out some warnings and eventually stop.

This heuristic will likely prevent an entire class of attacks and act as an effective defense in depth. 

#### Reliable source of time - preventing setting the time in the past

During normal operation, all Obscuro nodes will gossip payloads produced by the enclave containing a generated nonce, the last nonces received from peers, the current timestamp, and the head of the Ethereum chain.

Todo: more details about the payload.

The enclave will generate a new nonce and expect it to be included in the payloads of the other validators after it is fed a number of Ethereum blocks.

This mechanism makes the long-range attack more difficult because the attacker must also control a majority of the Obscuro nodes at the time of the fork.

Todo: more details about this protocol.  

#### Mandatory periodical re-attestation (every couple of months)

When a month has passed since the last version, a developer must submit a new Attestation (that can be the same code if there is no improvement) but will contain a "weak subjectivity checkpoint."

Once installed, the enclave generates a new shared secret.

This mechanism works because, in the pos protocol, validators can only withdraw with a limited rate. 
If the re-attestation period is well chosen, the enclave will always have a relatively up-to-date validator set.

The re-attestation period has to be enforced by the enclave to happen after every N epochs. It will function as a key rotation mechanism as well.

How does this mechanism protect against the attacks mentioned above?

Todo: Good question. It feels recursive


#### Proof-of-stake heuristics

We can add some pos-specific heuristics like:
The total amount of stake must be larger than an amount that is unlikely to be reached but still large enough not to be practical. 
The number of unfinalized epochs should never be larger than a couple of weeks. 



