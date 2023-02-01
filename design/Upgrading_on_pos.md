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
The purpose of this mechanism is to delegate to the community the final decision over who is deceitful: the developer or the challenger.
The community can step in and put their stake on the side they think is not deceitful. The side who loses will lose the entire deposited stake.
This is a mechanism of last resort which will act as a deterrent for anyone to act in bad faith.

Todo: establish an incentive for the community to vote.

The most likely scenario is that Developer D submits a version proposal with Stake Sd. Researcher R  spots something and places a challenge with stake Sr where Sr << Sd.
Then talk to each other offline, and either the researcher withdraws the challenge in a few days without any penalty or the researcher receives some percentage of Sd as a bug bounty, and the code must be fixed or completely withdrawn. 
If they don't agree, then they appeal to the community which will discover who is deceitful by reviewing the claims.
If one of them is obviously deceitful, then it's free money for the community.
If someone wants to push a malicious version, they will have to fight against the entire community to push it. And, similarly, if someone wants to stop a valid upgrade.

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
1. The current version of the enclave starts with a "weak subjectivity checkpoint" as part of its Attestation, which is the same as any beacon chain client.
Note: this first assumption leads to an interesting recursivity which we'll analyze below.

2. The attackers are part of the original validator set but only control a low amount of the total stake. 
Note: an attacker controlling a super-majority of the total Ethereum stake will be able to break the mechanism. 

3. The attackers aim to convince the enclave to upgrade itself to a malicious version without public scrutiny.

4. The enclave can execute the proof-of-stake protocol and asses the finality of blocks.

5. During any upgrade, the current enclave will share some secrets with the future enclave. 

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

As soon as they control the super-majority, they wake the enclave up and start feeding malicious finalised blocks with the same result as above.

Note 1: Another attack can combine the above techniques to reduce the number of unfinalised blocks.

Note 2: Both these attacks will reveal past data. They cannot read future data if the upgrading changes the secrets.

### Mitigations

Given these attacks' "subjective" nature, we are unlikely to find a definitive "objective" solution.
Instead, we can aim to introduce heuristics to make it improbable or to increase the capital required to impractical amounts.

The response to suspicious activity must be a gradual deterioration of the service until the enclave stops working altogether.

Note: In this document, the service is the upgrade process, but it similarly applies to the revelation period.

The risk of choosing too strict heuristics with abrupt decisions is that the network can become nonfunctional.

#### Reliable source of time - preventing fast-forwarding

A reliable time source inside the enclave can easily prevent the inactivity leak attack because it requires a fast-forward of time until the attackers dominate the total stake.

A relatively easy way to introduce a decentralized, trustless, cryptographic source of time inside the enclave is to feed it the Bitcoin block headers.
The Bitcoin blocks can be used as an approximate but reliable time source by applying the proof-of-work heuristic mentioned above.
With this mechanism in place, the attacker can slow down time, but they can only fast-forward it by spending a lot of money. 

The process itself is simple. The enclave will run the light-chain bitcoin code, resulting in a database entry with the approximate time.
If the timestamp of the Ethereum chain diverges significantly from the Bitcoin chain, the enclave will give out some warnings and eventually stop.

This heuristic will likely prevent an entire class of attacks and act as an effective defence in depth. 

#### Reliable source of time - preventing setting the time in the past

The heuristic described above gives us a reliable decentralised way of preventing fast-forwarding time. We need something similar for rewinding time. 

The requirement is that the enclave must discover when it is living in the past. Such a mechanism can prevent an entire class of attacks. For example, some replay attacks rely on an outdated and vulnerable enclave, maybe running on vulnerable hardware, being able to load old data and leak it to the attacker. Or the Ethereum "long range attack" where an enclave loaded from a snapshot from months ago is being shown a valid non-canonical Ethereum chain.

The requirement is that the host must convince the enclave that the data that it receives is real-time data and not a replay. 

The insight behind this solution is that the attacker cannot replay data if it contains randomness generated securely before the data is fed. Luckily, SGX is able to generate secure randomness.
The challenge is to design a protocol that uses this insight in a secure and practical way.

A simple solution would be to publish a nonce generated inside an enclave to Bitcoin and then for the enclave to check the timestamp of the block when it is included. To make an attack extremely expensive, the enclave could wait for some confirmations. Since the nonce that was just generated was included with a lot of confidence in a live chain, the timestamp of those blocks must be current. Unfortunately, this solution is very expensive, which makes it impractical. It can be used to form a mental model of the desired result of an alternative. 

A cheap and practical solution will use the existing Obscuro nodes. 

Since it is not practical to validate every data that goes into the enclave, the mechanism can periodically validate a timestamp (or equivalent) which the enclave can use to compare against the data it receives. If the time difference between the data and the validated timestamp is too large, the enclave will eventually stop and exit.


##### High level solution

Every time an enclave receives data showing that a new Ethereum epoch started, the enclave will generate a payload containing a newly generated nonce and the checkpoint hash. The host is responsible with collecting signatures from the other Obscuro nodes over this payload. Once it collects enough (2/3rd), it returns the proof into the enclave. Assuming the other nodes are not colluding against it, the current enclave has a good confidence level that everyone else has the same last epoch, because otherwise they would not sign this payload. 
The rule is that an enclave receiving a request from another enclave will sign over it only if the latest checkpoint coincides.
If the enclave does not receive the confirmation by the end of the next epoch, it will send a warning and eventually, after a few more tries, just shut down.

An attacker feeding old data into an enclave will not be able to produce the confirmations because the peers would not sign over a payload that points to an unknown or a past epoch. After a couple of epochs, the enclave will shut down, thus disabling the attack.

The most difficult problem with this approach is to establish the validity of the confirmation signatures themselves.
In the hypothetical "long range" scenario where the attacker spins up an enclave on a snapshot of a database from 6 months ago, the Obscuro nodes that this enclave considers active might no longer be.
If these nodes are decommissioned, then they will not respond, so the attack will fail quickly. 

Let's assume the attacker now controls these servers. To perform the attack, they must all be in synch, and respond to each other.
To do this, they must be started. If we apply the same mechanism during the startup of an enclave, then none of these malicious enclaves will be able to respond to the others to unlock them, since the only unlocked enclaves will be the ones operating normally.

Todo: We need a mechanism to prevent a validator from being considered active if it doesn't participate in the confirming protocol. This is to prevent an attacker from keeping a large number of validator dormant. (Is this a problem?)

/* Ignore this commented section
--

To achieve it, the enclave will periodically generate a nonce, and will expect that nonce to be included in payloads signed by its peers in a short time. 
Due to the availability of a cryptographic random number generator, the attacker will not be able to replay old data.

During a replay attack, the attacker will have to wait until the present to be able to return that nonce back into the enclave, which will be suspicious if the time difference between when the nonce was generated, and when it was returned is large enough. By adding more information, besides the nonce in the payload, the enclave can make a better decision.
With this mechanism in place, a successful attacker must control a significant number of peers (Obscuro nodes) to trick an enclave into a replay attack. 


All validators must sign with an epoch delay over a payload containing nonces from everyone else and the epoch.
MTree nonces + BLS signature.

something included in the rollup to pay incentive????

1. after every (finalised) epoch each enclave gossips a signed nonce. 
2. every enclave accumulates all these nonces.
3. at the start of the next epoch, each enclave - bls sig over a Mtree of all the nonces, accompanied by the nonces 
4. each enclave aggregates all these sigs over the same payload
5. each enclave checks that there are enough that included its own nonce - and thus get confirmation that they are running in the same timeline as the rest
6. all this gossiped stuff gets transformed to a tx, and contribute to liveness check


##### Mechanism

This mechanism will apply both during startup and during normal operation.

It consists of multiple steps.

todo - a single payload containing both the request and the response


###### Step 1: Requesting the latest attestation and time confirmation from the network

The "tick" of time, from the point of view of an enclave, is the Ethereum slots and blocks.
At the start of every epoch during live operation or at the end of the "catchup phase", the enclave will enter "Step 1" 

Step 1 is a function running inside an enclave that generates a signed payload containing a nonce and a list of N random Obscuro nodes.

The list of Obscuro nodes is available in the management contract.

This signed payload is the output of this step.
```
{ 
	node_id
	date 
	nonce
	current attestation
}
```



###### Step 2: Response from the network

The host component broadcasts this payload to the Obscuro network, which will respond with a signed response:

```
{
	nonce
	L1 block height & hash
	L2 height & hash
	current network attestation
	time
}
```

The role of the nonce is to make sure that what is presented to the enclave is the latest data, and not some replay.

The enclave considers a nonce as confirmed if at least 2/3rd of all the Obscuro nodes included this nonce in their own payload.


---
The reader might notice a contradiction here. This mechanism is designed to prevent the "long range" attack where the attacker has purchased the keys of Ethereum validators and now must play a different version of the history on an enclave. With this mechanism in place, to be able to use the enclave, the attacker must control a significant number of the Obscuro nodes at the time when the fork was created.  
There is a chicken and egg problem at play, because each of those nodes will wait for confirmation from the others. There are no other nodes that are live that can respond satisfactory and unlock them.



###### 3. Verification
Once the host has gathered signatures from enough enclaves, it will feed them into the enclave.

it will check that:
- ```requested.nonce == response.nonce``` 
- that the ```payload.attestation == enclave.attestation``` 
- AND that the required number of selected nodes have signed.

Note: that "attestation" could be a security version number to allow rolling upgrades.

If the checks pass then the enclave proceeds and loads secrets. 

If they don't match, the enclave just exits (and deletes all data?).


##### Attacks

####### I. The random number generator is vulnerable, which allows a replay attack. (the enclave being fed a previous payload). 

####### II. The attacker can start an enclave, and then keep the node dormant (without restarting it) for a long interval until a vulnerability is found.

To mitigate this, an enclave that is not being fed data for a while must exit. The difficulty is to measure time in the enclave.

The attacker can keep the CPU itself frozen using "sgx step" techniques? - this is a costly attack, as it completely blocks a server.


During normal operation, all Obscuro nodes will gossip payloads produced by the enclave containing a generated nonce, the last nonces received from peers, the current timestamp, and the head of the Ethereum chain.

Todo: more details about the payload.

The enclave will generate a new nonce and expect it to be included in the payloads of the other validators after it is fed a number of Ethereum blocks.

This mechanism makes the long-range attack more difficult because the attacker must also control a majority of the Obscuro nodes at the time of the fork.

Todo: more details about this protocol.  

*/

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



