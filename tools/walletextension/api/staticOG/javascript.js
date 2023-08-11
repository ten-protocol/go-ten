const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idJoin = "join";
const idAddAccount = "addAccount";
const idStatus = "status";
const pathJoin = "/join/";
const pathAuthenticate = "/authenticate/";
const obscuroChainIDDecimal = 777
const methodPost = "post";
const methodGet = "get"
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};
const metamaskRequestAccounts = "eth_requestAccounts";
const metamaskPersonalSign = "personal_sign";

function isValidUserIDFormat(value) {
    return typeof value === 'string' && value.length === 64;
}

async function addNetworkToMetaMask(ethereum, userID, chainIDDecimal) {
    // add network to MetaMask
    let chainIdHex = "0x" + chainIDDecimal.toString(); // Convert to hexadecimal and prefix with '0x'

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
        return false
    }
    return true
}

async function authenticateAccountWithObscuroGateway(ethereum, account, userID) {
    const textToSign = "Register " + userID + " for " + account;
    const signature = await ethereum.request({
        method: metamaskPersonalSign,
        params: [textToSign, account]
    }).catch(_ => { return -1 })
    if (signature === -1) {
        return "Signing failed"
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
    return await authenticateResp.text()
}

const initialize = () => {
    const joinButton = document.getElementById(idJoin);
    const addAccountButton = document.getElementById(idAddAccount);
    const statusArea = document.getElementById(idStatus);

    // get ObscuroGatewayUserID from local storage
    let userID = localStorage.getItem("ObscuroGatewayUserID")

    // check if userID exists and has correct type and length (is valid) and display either
    // option to join or to add new account to existing user
    if (isValidUserIDFormat(userID)) {
        statusArea.innerText = "Your userID is: " + userID
        joinButton.style.display = "none"
        addAccountButton.style.display = "block"
    } else {
        joinButton.style.display = "block"
        addAccountButton.style.display = "none"
    }

    let ethereum = window.ethereum;
    if (!ethereum) {
        joinButton.style.display = "none"
        addAccountButton.style.display = "none"
        statusArea.innerText = "Please install MetaMask to use Obscuro Gateway"
    }


    joinButton.addEventListener(eventClick, async () => {
        // join Obscuro Gateway
        const joinResp = await fetch(
            pathJoin, {
                method: methodGet,
                headers: jsonHeaders,
            }
        );
        if (!joinResp.ok) {
            statusArea.innerText = "Failed to join. \nError: " + joinResp
            return
        }

        // save userID to the localStorage and hide button that enables users to join
        userID = await joinResp.text();
        localStorage.setItem("ObscuroGatewayUserID", userID);
        joinButton.style.display = "none"

        // add Obscuro network to Metamask
        let networkAdded = await addNetworkToMetaMask(ethereum, userID, obscuroChainIDDecimal)
        if (!networkAdded) {
            statusArea.innerText = "Failed to add network"
            return
        }

        // get accounts from metamask
        const accounts = await ethereum.request({method: metamaskRequestAccounts});
        if (accounts.length === 0) {
            statusArea.innerText = "No MetaMask accounts found."
            return
        }
        let authenticateAccountStatus = await authenticateAccountWithObscuroGateway(ethereum, accounts[0], userID)

        statusArea.innerText += "\n Authentication status: " + authenticateAccountStatus

        // show users an option to add another account
        addAccountButton.style.display = "block"
    })

    addAccountButton.addEventListener(eventClick, async () => {
        // check if we have userID and it is correct length
        if (!isValidUserIDFormat(userID)){
            statusArea.innerText = "\n Please join Obscuro network first"
            joinButton.style.display = "block"
            addAccountButton.style.display = "none"
        }

        // Get account and prompt user to sign joining with selected account
        const accounts = await ethereum.request({method: metamaskRequestAccounts});
        if (accounts.length === 0) {
            statusArea.innerText = "No MetaMask accounts found."
            return
        }
        let authenticateAccountStatus = await authenticateAccountWithObscuroGateway(ethereum, accounts[0], userID)
        statusArea.innerText += "\n Authentication status: " + authenticateAccountStatus
    })

}

window.addEventListener(eventDomLoaded, initialize);