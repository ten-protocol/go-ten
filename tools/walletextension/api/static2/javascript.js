const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idJoin = "join";
const idAddNetwork = "addNetwork"
const idStatus = "status";
const idAuthenticate = "authenticate";
const pathJoin = "/join/";
const pathAuthenticate = "/authenticate/";
const methodPost = "post";
const methodGet = "get"
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};
const metamaskRequestAccounts = "eth_requestAccounts";
const metamaskPersonalSign = "personal_sign";

const initialize = () => {
    const joinButton = document.getElementById(idJoin);
    const addNetworkButton = document.getElementById(idAddNetwork)
    const authenticateButton = document.getElementById(idAuthenticate)
    const statusArea = document.getElementById(idStatus);

    let userID;

    joinButton.addEventListener(eventClick, async () => {
        const joinResp = await fetch(
            pathJoin, {
                method: methodGet,
                headers: jsonHeaders,
            }
        );

        if (!joinResp.ok) {
            statusArea.innerText = "Failed to join."
            return
        }

        userID = await joinResp.text();
        statusArea.innerText = "Your userID is: " + userID

        console.log("User ID is: ")
        console.log(userID)
        addNetworkButton.style.display = "block"

    })

    addNetworkButton.addEventListener(eventClick, async () => {
        let ethereum = window.ethereum;

        let chainIdDecimal = 777;
        let chainIdHex = "0x" + chainIdDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'

        if (ethereum) {
            try {
                await ethereum.request({
                    method: 'wallet_addEthereumChain',
                    params: [
                        {
                            chainId: chainIdHex,
                            chainName: 'Obscuro Testnet',
                            nativeCurrency: {
                                name: 'Obscuro',
                                symbol: 'OBX',
                                decimals: 18
                            },
                            rpcUrls: ['http://127.0.0.1:3000/?u='+userID],
                            blockExplorerUrls: null,
                        },
                    ],
                });
            } catch (error) {
                console.error(error);
            }
            authenticateButton.style.display = "block"
        } else {
            alert('MetaMask is not installed. Please install MetaMask and try again.');
        }
    });

    authenticateButton.addEventListener(eventClick, async () => {
        if (typeof ethereum === "undefined") {
            statusArea.innerText = "`ethereum` object is not available. Please install and enable MetaMask."
            return
        }

        const accounts = await ethereum.request({method: metamaskRequestAccounts});
        if (accounts.length === 0) {
            statusArea.innerText = "No MetaMask accounts found."
            return
        }
        // Accounts is "An array of a single, hexadecimal Ethereum address string.", so we grab the single entry at index zero.
        const account = accounts[0];

        const textToSign = "Register " + userID + " for " + account;
        const signature = await ethereum.request({
            method: metamaskPersonalSign,
            // Without a prefix such as 'vk', personal_sign transforms the data for security reasons.
            params: [textToSign, account]
        }).catch(_ => { return -1 })
        if (signature === -1) {
            statusArea.innerText = "Failed to sign viewing key."
            return
        }

        const authenticateUserURL = pathAuthenticate+"?u="+userID
        const authenticateFields = {"signature": signature, "message": textToSign}
        const authenticateResp = await fetch(
            authenticateUserURL, {
                method: methodPost,
                headers: jsonHeaders,
                body: JSON.stringify(authenticateFields)
            }
        );
        let authenticateStatus = await authenticateResp.text();
        statusArea.innerText += "\n Authentication status: " + authenticateStatus
    });

}

window.addEventListener(eventDomLoaded, initialize);