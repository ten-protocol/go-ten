
## Handling catastrophic events

TEN is facing more risk than a transparent network during unforeseen situations.
With traditional software, if a bug prevents all nodes from starting up, the developers can quickly fix
the bug, release the version, and the network will proceed. This works because there are no visibility restrictions on the existing data.
As long as the data is not corrupted beyond recovery, there is always a path forward.

The problem with TEE software is that upgrading must be very restricted for privacy reasons.
If the rules are too strict, they could combine with a software bug and leave the network completely bricked.
For example, during an upgrade, the old version will stop working at block 1000, but the new version cannot start either
because it crashes.
Or, it is possible that all nodes spontaneously crash during a normal operation before having the chance to hand over
the secret to an upgrade that fixes the bug.

One way to mitigate this is diversity in software. This would ensure that not all implementations are hit by the same problem.
Reaching software diversity will take a period, though.

Another option to mitigate the problem is to relax the security constraints at the cost of risking compromising privacy, at least for some parties.
One possibility is to allow the original application developer to unilaterally upgrade to a new version.

The preferred option is to create a "Safe mode", a very simplified code path inside the enclave whose only task is to hand out
the current Master Secret to an approved new version.
This will be started in command line only, receive the upgrade proof and attestation, and output the secrets.

### Responsibilities:

1. Interpret command line parameters.
1. Verify an attestation.
1. Verify the upgrade proof.
1. Unseal the master seed.
1. Encrypt the master seed with the key from attestation.

All steps except the third are relatively straightforward. We can use established, well-tested libraries, and very little
custom logic.

### Verifying the upgrade proof

The decentralised mechanism described in the "Weak subjectivity" document requires the enclave code to understand the consensus
protocol of the beacon chain. That is complex code which can go wrong in unforeseen ways.

Ideas:
We need an alternative simpler mechanism for this exceptional situation.
Ideally, a mechanism that relies on verifying digital signatures and Merkle Trees.

To reduce the scope of abuse, we propose that, during phase 1, only the sequencer can enter "Safe mode".

Upgrade proof signed by an ad-hoc "upgrade oracle" composed of most of the TEN node operators at the time
of each version. The public keys of all these nodes will be included in the image.


### Backup Key

Given there are still risks even with a "Safe Mode", the safest way is to start the network with "training wheels" on.
This means there should be as little code as possible to minimise the risks.

On a high level, the solution is to encrypt the master seed of each major version using "Threshold encryption".
The participants in this group encryption scheme must be chosen from the initial node operators.

Each participant will publish a public key to the management contract and will receive back their encrypted share of the master seed
split up using Shamir's Secret Sharing algorithm.

It is assumed that each participant has a secure process in place using HSMs where they generate and guard the key material
used for this scheme.

In phase 1, the first enclave for each new major version will create this backup before it starts producing rollups.

In case something goes wrong catastrophically, the participants in this group will collaborate and restore the functionality.