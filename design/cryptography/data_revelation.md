# Data Revelation

## Scope

The data revelation ensures that Obscuro transactions encrypted with different keys, can be revealed independently at the right time.
 
The mechanism needs a reliable way to measure the time that cannot be easily attacked.

## Requirements
* Keys must be released only after X amount of time has passed
* Clock mechanism must determine how much time has passed

## Design

### Keys must be released only after X amount of time has passed

A database stores the rollup, decrypt time, key tuple.
When a new block arrives the enclave stores the current time and all keys before that time are now available.
Validators fetch the key from the enclave and store the keys locally for redundancy.

### Clock mechanism must determine how much time has passed

Given the predictable block creation rate the Revelation Mechanism Clock is based on the number of blocks that have been issued.

In ethereum blocks have a rate of ~12seconds. This is the standard tick of the revelation clock.

## Attacks

### Time fast-forward

This attack is described as someone fast forwarding the clock and having access to the revelation key before the time it would be expected.
The attack is mitigated because the keys are never released from the validators.

