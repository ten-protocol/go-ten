# An objective source of truth for the TEN Enclave

A TEE-based network like TEN needs a reliable and objective source of truth to function correctly and protect the
data privacy in the face of complex attacks by the node operator.

## Problems

The lack of reliable source-of-truths mechanisms can lead to the following types of attacks:
- the revelation period can be gamed
- attacks against previous versions of the enclave that have known vulnerabilities can reveal historical data
- transparent upgrades cannot be implemented

Note: we will analyse each of these attacks below.

TEN is an Ethereum Layer 2, so it is natural to consider the Ethereum chain as the source of truth.
The transition to proof-of-stake is making the problem even more difficult because one of the tradeoffs made when
designing that protocol was to weaken the objectivity assumption to "Weak subjectivity".

### Revelation period

TEN has a built-in mechanism to reveal encryption keys for transactions after a configured amount of time.
The difficulty is how to measure time reliably inside an enclave.

When Ethereum was based on proof-of-work, we could rely on the number of Ethereum blocks that have been consumed by the enclave
as long as the "work" that went into calculating them was above a threshold.
This simple mechanism made an attack prohibitively expensive.

We will list below a couple of attacks that are possible against a proof-of-stake network that make achieving the guarantee more
challenging.

### Attacks against historical vulnerable versions

The threat is that a hacker learns about an SGX vulnerability after there is already a patch available and wants to use
compromised hardware on a snapshot of the network state to decrypt the historical data.

Note: TEN reveals data anyway, so the impact of such an attack in real life would be reduced, but we need to make our best effort to keep
the guarantee of the revelation period.


### Transparent TEN Upgrades

A TEE-based system like TEN must be upgradeable.
During an upgrade, the new software version must be able to pick up from the old version, which means that the previous
version must share some of its secrets with the new version.

If the upgrade mechanism is opaque and centralised, whoever controls it can unilaterally upgrade the software to a malicious
version and gain access to the secrets. The possibility of such an attack performed behind closed doors can severely
undermine the guarantees of a privacy-preserving solution.

Note: This is a completely different class of attacks from the ones against the typical transparent chains.

**Requirement: Every release must have a built-in *decentralised and transparent* mechanism for accepting future versions
of itself as valid and, more importantly, for rejecting malicious versions.**


#### Solution 1 - Upgrading based on a single signer

The most naive upgrading mechanism requires the same private key to sign the new version as the key used to sign the current one.
This straightforward method can be seen as the same developer being trusted with upgrading.
The problem with this approach is that, without additional protections, it gives the developer the power to retrieve the secrets unchallenged.
They can create a new malicious version, upgrade a node decoupled from the internet to this version, and access the secrets.

This mechanism is simple but opaque and centralised.

#### Solution 2 - Upgrading based on multiple signers

An improved option is a Multi-Sig approach, where a committee of independent parties must sign over each upgrade.
The developer proposes a new version and submits it for auditing to the members, who, upon careful analysis, will sign
over it if they think it is valid and protects their interests.

This could be feasible if the committee members have a good reputation and the right incentives. But unfortunately,
it is challenging to find such reputable parties, especially for a young network.
If there is suspicion of collusion between them, the perception of security can be permanently compromised.
Also, operationally, it makes upgrades more difficult since every committee member must review and sign the code during
a relatively short period.


#### Solution 3 - Upgrade lifecycle managed by a smart contract

Another option is to leverage the Ethereum Layer 1 for transparency and decentralisation by entrusting the task of managing
upgrades to a smart contract.
The TEE would accept a new software version as a valid upgrade if it went through a predetermined lifecycle on Ethereum enforced
by the smart contract.

At the end of the lifecycle, the smart contract emits an `UpgradeEvent`, which will be consumed by the enclave.

This decouples the process of whitelisting a new version from its result.

This means that the process can start out more centralised and eventually progress to the decentralised solution described below.

##### 3.1 Centralised Transparent Smart-contract managed upgrade

In phase 1, the smart contract can have a simple ``approve`` function, which will trigger the event when called by an Admin.

This mechanism is a transparent version of "Solution 1", as it prevents the Admin from unilaterally performing
an upgrade in a dark room without anyone knowing.

##### 3.2 Decentralised Smart-contract managed upgrade

In phase 2, we can decentralise by creating a more complex incentive-driven lifecycle.

1. A developer proposes a new version by publishing a GitHub tag (hash) and the Attestation of an enclave built from that code,
   together with a stake. Initially, the developer will likely be the TEN Foundation.

2. Anyone posting a stake can make a challenge to this version. The challenge period is predetermined.

3. The original developer first evaluates the challenge. If it is valid, then a reduced stake is paid out to the challenger as a bug bounty.

4. If the developer rejects the challenge, the challenger can request a vote from the community.  
   The purpose of this mechanism is to delegate to the community the final decision over who is deceitful: the developer or the challenger.
   The community can step in and put their stake on the side they think is right. The side that loses will lose the entire deposited stake.
   This is a mechanism of last resort which will act as a deterrent against anyone acting in bad faith.

Todo: establish an incentive for the community to vote.

The most likely scenario is that Developer D submits a version proposal with Stake `Sd`. Researcher `R`  spots something
and places a challenge with stake `Sr` where `Sr << Sd`.
Then talk to each other offline, and either the researcher withdraws the challenge in a few days without any penalty or
the researcher receives some percentage of `Sd` as a bug bounty, and the code must be fixed or completely withdrawn.
If they don't agree, then they appeal to the community, which will discover who is deceitful by reviewing the claims.
If one of them is obviously deceitful, then it's free money for the community.
If someone wants to push a malicious version, they will have to fight against the entire community to push it. And similarly, if someone wants to stop a valid upgrade.

This mechanism has several advantages:
- It is fully transparent. Everyone can inspect the code before it gets rolled out. And the history of all upgrades is recorded on a transparent, immutable chain.
- The code must be published and go through a transparent lifecycle before it receives any secrets. This removes the power of a developer or committee to gain a stealth advantage.
- It incentivises security researchers to find issues in the codebase.
- It is as decentralised as the smart contract, which can be easily audited.

#### Transparent upgrades conclusion

TEN will implement the transparent, smart-contract-driven upgrade process, starting out with the centralised approach and then decentralising in the next phase.

The design needs two key ingredients:
1. The enclave must understand the outputs of the smart contract.
2. The node operator must convince the enclave that the data it feeds is the actual Ethereum canonical chain.

The first part is straightforward because the enclave can interpret and authenticate events emitted on Ethereum against a chain.
This functionality is already required for cross-chain messaging and authenticating rollups.

The second point is difficult to achieve in Ethereum 2 because of 'weak subjectivity'.
In the next section, we'll see how, in the context of TEEs, the "weak" turns into "strong" subjectivity.

## Weak subjectivity of Ethereum 2

"Weak subjectivity" is a tradeoff of the Ethereum 2 consensus design.
What it means is that if a new node joins Ethereum 2 after a few years, it won't be able to tell for sure if the chain
that is presented to it is the canonical chain.
In proof-of-work chains, a node could look at the total amount of work that went into a chain, which is an external,
"objective" means.
The solution is for a node joining the network after a while to start up using a "weak subjectivity checkpoint",
which it has to retrieve from a trusted friend or website.
The assumption is that Ethereum nodes are run by operators who are incentivised to run against a canonical fork, so it
makes sense for them to perform this extra simple step.

While weak subjectivity poses no real problem against the Ethereum network, it becomes a real problem for TEEs.

### Weak subjectivity and TEEs

The threat against an enclave differs from the threat against an Ethereum node because an enclave has no trusted friends
to feed it with a valid "Weak Subjectivity Checkpoint".
Its only link to the outside world is the node operator, who is the attacker trying to exploit this in order to extract secrets.
The operator controls the data flow into the enclave, including the clock.
Without additional mechanisms, there is no way for an enclave to tell if it wasn't fed any data for a month or for a second.

In particular, the threat against the upgrading mechanism is that an attacker can bypass the transparency guarantees by producing a valid
Ethereum fork, different from the canonical chain, contains a whitelisted but malicious version of the enclave.
The crypto-economic guarantees of Proof-of-Stake don't apply in this case because the attacker will not broadcast the fork to the outside world,
so breaking the protocol will not be discovered by anyone else.

A naive solution is to require new versions to be signed by a committee, which brings us back to an opaque solution.

This document will not go into the details of the Ethereum proof-of-stake protocol because it is very complex.
Instead, we'll present some possible practical attacks that illustrate the threat and propose heuristics designed to mitigate them.


### Weak subjectivity attacks against a TEE

Assumptions:

1. The current version of the enclave starts with a "weak subjectivity checkpoint" as part of its Attestation, the same as any beacon chain client.
   Note: this first assumption leads to an interesting recursivity which we'll analyse below. (TODO)

2. The attackers are part of the original Ethereum validator set but only control a low amount of the total stake.
   Note: an attacker controlling a super-majority of the total Ethereum stake will be able to cause bigger harm to the ecosystem.

3. The attackers aim to extract secrets from the enclave ahead of time.

4. The enclave can execute the proof-of-stake protocol and asses the finality of blocks.

#### Inactivity leak attack

The attackers want to exploit the "inactivity leak". Given that they are already registered validators, the simplest
the attack is for these validators to attempt to create a parallel chain where everyone else has suddenly dropped off.

They will produce blocks when the protocol decides it is their turn, and they will be the only block producers left.
Since the stake they control is very low, the checkpoints will not finalise. After four epochs, the "inactivity" rule will
kick in, which will chip away at the stake of all the inactive validators (basically everyone except the attackers).
After many blocks, the attackers will control 2/3rd of the entire stake and will be able to finalise checkpoints.

Note that the blocks produced during this attack are not broadcast, so the validators will not get penalised in the actual
canonical chain for signing conflicting attestations.

As a result of this attack, the attacker can eventually feed finalised blocks into the enclave. These finalised blocks can
contain upgrade transactions that appear valid to the enclave but violate the transparency assumption.

Without additional heuristics, the enclave will receive a valid Upgrade event from the smart contract and will assume
that it can safely release secrets to that version.

#### Long-range attack

A variation of the previous attack is the typical "long-range" attack from the Ethereum literature.
Validators withdraw their stake and sell their keys on the black market, and the enclave attackers buy them.

To prepare, the attackers will keep an enclave dormant until they control a super-majority of the stake at a time when
the enclave goes to sleep.

As soon as they control the super-majority, they wake the enclave up and start feeding malicious finalised blocks.

Note 1: Another attack can combine the above techniques to reduce the number of unfinalised blocks.

## Proposed solutions

Given these attacks' "subjective" nature, it is challenging to find an "objective" solution.
Instead, we can aim to introduce heuristics to make the attack more difficult or to increase the capital required to impractical amounts.

The response of the enclave to suspicious activity must be a gradual deterioration of the service until the enclave stops working altogether.
The risk of choosing too strict heuristics with abrupt decisions is that the network can become nonfunctional.


### Preventing time fast-forwarding or feeding data with future timestamps

The inactivity leak attack requires that the time inside the enclave will be ahead of the real-time because it takes many
blocks for the attackers to dominate the total stake.

A relatively easy way to introduce a decentralised, trustless, cryptographic source of time inside the enclave is to feed
it the Bitcoin block headers.
Bitcoin can be used as an approximate time source by applying the proof-of-work threshold heuristic mentioned above.

If this mechanism is implemented in the enclave, the attacker can slow down time, but they can only fast-forward it by spending a lot of money.
Basically, the attacker needs to produce blocks faster than the entire Bitcoin network for a significant period of time.

The process itself is simple. The enclave will run the light-chain bitcoin code and will extract the timestamp from the
blocks that are being fed. If the timestamp of the Ethereum chain diverges significantly from the Bitcoin chain,
the enclave will give out some warnings and eventually stop.

This heuristic will prevent an entire class of attacks and act as an effective defence in depth against others.

### Preventing feeding past data or starting up an outdated enclave

The heuristic described above gives us a reliable decentralised way of preventing fast-forwarding time, but an even larger
threat is posed by the ability of an attacker to replay old data.

For example, the class of attacks like the Ethereum "long range attack" where an enclave loaded from a snapshot
from months ago is being presented as a valid but non-canonical Ethereum chain that contains a malicious upgrade.

**The high-level requirement for this mechanism is that the enclave must be convinced that the data that it receives
is real-time data and not a replay, and an enclave won't load any data in memory without cryptographic proof
that it is up-to-date.**

*The insight behind the solution is that the attacker cannot replay data if that data contains randomness generated securely
just before the data is fed.*

The challenge is to design a protocol that uses this insight in a secure and practical way.
The protocol has to be interactive because it necessarily involves two steps:
1. Generating the randomness inside the TEE.
2. Somehow, including the randomness inside the live data.

A simple solution is to generate the nonce, then publish it to Bitcoin, and then feed the Bitcoin block containing that nonce back into the enclave.
If the timestamp of this block is ahead of the current time, then the enclave knows it is being attacked.
To make this attack extremely expensive, the enclave could wait for some confirmations on top of that block.
This approach leverages the objectivity of a pow network, but unfortunately, it is expensive, which makes it impractical.

We'll use it as a reference for the mental model of the desired result.

The protocol in essence uses a liveliness nonce and is an adapted challenge-response mechanism that ensures the freshness
of the data fed into the enclave.

Note that the naive solution described above does not validate every piece of data that goes into the enclave.
It is a mechanism that periodically validates a timestamp (or equivalent) that the enclave can use to compare against the data it receives.
If the time difference between the data and the validated timestamp is too large, the enclave will eventually stop and exit.

To prevent an outdated enclave from starting up, the nonce generated inside it must be included in a Bitcoin block, but also
the latest attestation constraints must be published to Bitcoin. And at startup, the enclave would have to walk back the chain
until it finds the latest Attestation, and compare that to itself. If the current code fails the constraints, the enclave
will shut down before loading any secrets in memory.

In the next sections, we'll describe some more complex solutions.

Note: These are still unreviewed draft proposals at this stage.

#### Solution 1 - Batched publishing to Bitcoin

The naive solution described above is inefficient because every single TEN node has to periodically publish a transaction to Bitcoin.
One immediate improvement is for the TEN nodes to join forces and generate shared randomness (`SR`) during a cycle,
and then someone publishes this `SR` to Bitcoin.
This achieves the same result, as it convinces every enclave that participated in this cycle that its timestamp is valid.

A reasonable cycle should be as long as possible to minimise the costs and as short as possible to reduce the chance of time-related attacks.

An attacker is able to skew time by a maximum of 2 cycles. (todo - explanation)

The lifecycle of an upgrade, the main vulnerability, is measured in days, which means an attack will require waiting for at least one day.

Note: time inside the enclave is measured based on the timestamp of the Ethereum blocks fed into it.

For the rest of the proposal, we'll assume that the cycle duration is 12 hours.
Todo: This reasoning requires more rigour.

If the entire network has to submit a Bitcoin transaction every 12 hours, it implies a cost of a couple of USD per day for the entire network,
which is reasonable.

##### Protocol

Draft:

Every N epochs, each enclave will generate a payload containing a random nonce. This payload will be gossiped around as a normal p2p transaction.
After 10-20 epochs (to give all enclaves a chance to submit a nonce), all the nonces are added to a Merkle Tree.

Q1: why would someone aggregate these nonces? Find incentives

Q2: what exactly is published, by whom, why? Incentives


#### Solution 2 - Nonce validated by the TEN network

Another potential solution is to rely on the TEN network itself to validate the randomness.
This is more tricky to achieve because of the subjectivity of the TEN network participants.

##### Protocol

Draft:

Every N Ethereum epoch started, the enclave will generate a payload containing a newly generated nonce, the last checkpoint hash,
and the current Attestation.
The host is responsible for collecting signatures from the other TEN nodes over this payload.
The rule is that an enclave receiving a request from another enclave will sign over it only if the latest checkpoint coincides,
and the Attestation is valid.

Once it collects signatures from a majority of its peers, it returns the proof to the enclave.
Assuming the other nodes are not colluding against it, the current enclave has a good confidence level that everyone else
agrees with this payload because otherwise, they would not sign over it.
The signature happens inside the enclave, so unless the enclave can be impersonated, the signature cannot be forged.

If the enclave does not receive the confirmation by the end of the next epoch, it will send a warning and eventually,
after a few more tries, it just shuts down.

An attacker feeding old data into an enclave will not be able to produce the confirmations because the peers would not
sign over a payload that points to an unknown or a past epoch. After a couple of epochs, the enclave will shut down, thus disabling the attack.

The most difficult problem with this approach is establishing the validity of the confirmation signatures themselves.
In the hypothetical "long-range" scenario where the attacker spins up an enclave on a snapshot of a database from 6 months ago,
the TEN nodes that this enclave considers active might no longer be.
If these nodes are decommissioned, then they will not respond, so the attack will fail quickly.

The problem comes when we assume the attacker now controls these servers. To perform the attack, they must all be in synch and respond to each other.
To do this, they must be started. If we apply the same mechanism during the startup of an enclave, then none of these
malicious enclaves will be able to respond to the others to unlock them since the only unlocked enclaves will be the ones operating normally.

Todo: We need a mechanism to prevent a validator from being considered active if it doesn't participate in the confirming protocol.
This is to prevent an attacker from keeping a large number of validators dormant. (Is this a problem?)

Todo-define the protocol properly, including payloads, etc.

todo - catchup vs live

weak subjectivity checkpoint - signatures over it  

Todo classify all the possible cases:
- start a new node from new version
-start a node after a short period
- start a node after a long period
  ...m




### Proof-of-stake heuristics

Another type of defence is to create some proof-of-stake specific heuristics like:

- The total amount of stake must be larger than a preset threshold amount that is unlikely to be reached but still large enough not to be practical.
- The number of unfinalised epochs should never be larger than a couple of weeks.

Such heuristics should be designed to be highly unlikely to happen in practice, but likely during an attack.

