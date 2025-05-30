# Rogue P2P Node Operator

The main risk for this is DoS to other nodes and we currently have unauthenticated TCP connections which could allow MITM attacks. 

## Message flooding 

The code accepts unlimited connections and processes messages in goroutines `handleConnections()` and `handle()`. A rogue node could establish multiple connections and flood the network with messages. There's no rate limiting 
This could lead to resource exhaustion (memory, CPU) on the target node

## Oversized messages

The code checks message types `_minMsgType` to `_maxMsgType` but doesn't enforce a limit on the size of message contents.  A rogue node could send extremely large messages to consume memory. In `handle()` we start deserializing the message contents without checking. 

## Spam batch request from Sequencer

A rogue node could repeatedly request batches using `RequestBatchesFromSequencer()`, forcing the sequencer to process these requests. This could cause a DoS on the sequencer. 

## Health check manipulation
The `verifyHealth()` function relies on message timestamps from peerTracker
A rogue node could send periodic small messages to appear healthy while actually being malicious. This could prevent the health check from detecting other types of attacks


## Implement missing batch signature check 

In processing of `msgTypeBatches` we don't check the signature on the batch which opens the window for malicious batch data to get in during a MITM attack on the p2p connection. 

## Authenticate TCP connection 

A man in the middle attack on the sequencer connection would allow someone to distribute a malicious transaction message due to the encryption not allowing message verification. 


## Recommendations

* Implement rate limiting for 
    * Connections per peer
    * Messages per peer
    * Batch requests 
* Message size limits for any message type
* Validation of batch contents/ signature before processing 
* Maximum number of concurrent connections 
* Authenticate the TCP connections 
