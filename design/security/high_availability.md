# High availability

The TEN Sequencer must have HA capabilities. The reasoning is covered in the "Fast finality" design docs.

The requirement is that the service must continue even if the enclave of the sequencer crashes and is corrupted.

One option is to require a governance event to whitelist a new sequencer enclave in case the current one becomes unavailable.
This is not a good solution, because it would lead to significant downtime.

The preferred option is to mark multiple enclaves as sequencers, and make it the responsibility of the operator to switch
between them to ensure that service goes on uninterrupted.

## MEV Protection 

The Sequencer is the only TEN node capable of producing Light Batches, and thus the only node that can in theory
attempt to extract value via MEV. If the operator has multiple enclaves available operating in "Sequencer" mode it can both 
provide an uninterrupted service and extract some value at the same time.

Introducing startup delays does not help too much in this case, because the operator could hold key transactions for longer.


An alternative solution is to introduce transparency into the lifecycle events of the sequencer enclaves, such that the TEN network 
can assess the likelihood of bad behaviour.  

Lifecycle events:
- enclave starting up
- enclave become active
- enclave becoming passive

On a high level, all enclave lifecycle events of all the sequencer whitelisted enclaves must be published to the network.

### Protocol

Each Batch will have three elements:

1. The header
2. The transaction payload
3. The protocol payload

The protocol payload will not be included in the rollup published to the data availability layer. 

Note: the protocol payload can be used for other protocol specific messages, like current attestations.

```golang
type RestartEvent struct {
	enclaveId                   EnclaveId
	lastRestartTime             Timestamp
	batchHeadAtRestart          Hash
	currentBatchHead            Hash
	signature                   Signature
}

type ProtocolPayload struct {
    events []LifecycleEvent	
	
}
```

Todo - each event points to the previous event to prevent the operator from transitioning from active to passive without declaring.
?

Periodically (todo - how often? Every batch?), the sequencer host will request a `RestartEvent` from each of the whitelisted enclaves it controls.
It will add these events to the `ProtocolPayload`, and broadcast them to the TEN network together with the Batch.

Upon restart, each enclave records the required data as a variable, and will return that variable in the right struct each time it is being asked. 
This proof cannot be forged without a significant bug in the software or impersonation of the enclave. 

These events must be fed into the enclave producing batches, who will include one for itself.

If the operator does not include the events for the other whitelisted enclaves, it's a sign of foul play.

In case one of the whitelisted enclaves is corrupted, and cannot be restarted, the sequencer operator must publish a request to the
management contract to remove it from the whitelist, and to be replaced with another.

This means, there should not be long period of time with radio silence from any of the whitelisted enclaves.
