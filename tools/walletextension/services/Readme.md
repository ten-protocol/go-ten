# Implement session keys - guide for app developers 

If the user selects "no-click UX" (or something) when the game starts, do the following:

1) call eth_getStorageAt with the address 0x0000000000000000000000000000000000000003 (the other parameters don't matter). This will return the address of the session key.
2) Create a transaction that transfers some eth to this address. (Maybe you ask the user to decide how many moves they want to prepay or something). The user has to sign this in their wallet. Then submit the tx.
3) Once the receipt is received you call eth_getStorageAt with 0x0000000000000000000000000000000000000004 . This means that you tell the gateway to activate the session key.
4) All the moves made by the user now can be sent with eth_sendRawTransaction or eth_sendTransaction unsigned. They will be signed by the gateway with the session key.
5) When the game is finished create a tx that moves the funds back from the SK to the main address. This will get singed with the SK by the gateeway
6) Call: eth_getStorageAt with 0x0000000000000000000000000000000000000005 - this deactivates the key.