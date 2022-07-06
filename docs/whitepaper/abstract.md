# Abstract
_The full whitepaper is available to view [here](https://whitepaper.obscu.ro)._

We present Obscuro, a decentralised Ethereum Layer 2 Rollup protocol designed to achieve data confidentiality and prevent [Maximal Extractable Value (MEV)](https://ethereum.org/en/developers/docs/mev/) by leveraging hardware-based [Trusted Execution Environments (TEE)](https://en.wikipedia.org/wiki/Trusted_execution_environment).

The design of Obscuro ensures that hardware manufacturers do not have to be trusted for the safety of the ledger. If one manufacturer turns malicious or there is a breach in the TEE technology, the protocol falls back to the behaviour of a public blockchain that preserves the ledger's integrity but makes the transactions public. This situation will lead to a partial liveness failure because the withdrawal function will be suspended.

The design also focuses on preserving privacy for the limited period when it matters most, which removes the need for a privacy technique that is robust against all adversaries in perpetuity.

Obscuro sits in what we believe is a sweet spot between the existing [rollup-based L2 offerings](https://ethereum.org/en/developers/docs/scaling/layer-2-rollups/) Optimistic and ZK Rollups. The use of [confidential computing techniques](https://www.intel.co.uk/content/www/uk/en/security/confidential-computing.html) coupled with economic incentives allows Obscuro to retain the performance and programming model simplicity of Optimistic rollups, and on top of that attain confidentiality, short withdrawal periods, and address MEV.