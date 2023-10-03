const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idJoin = "join";
const idAddAccount = "addAccount";
const idAddAllAccounts = "addAllAccounts";
const idRevokeUserID = "revokeUserID";
const idStatus = "status";
const idUserID = "userID";
const obscuroGatewayVersion = "v1"
const pathJoin = obscuroGatewayVersion + "/join/";
const pathAuthenticate = obscuroGatewayVersion + "/authenticate/";
const pathQuery = obscuroGatewayVersion + "/query/";
const pathRevoke = obscuroGatewayVersion + "/revoke/";
const obscuroChainIDDecimal = 443;
const methodPost = "post";
const methodGet = "get";
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};
const metamaskRequestAccounts = "eth_requestAccounts";
const metamaskPersonalSign = "personal_sign";

function isValidUserIDFormat(value) {
    return typeof value === 'string' && value.length === 64;
}

let obscuroGatewayAddress = window.location.protocol + "//" + window.location.host;


async function addNetworkToMetaMask(ethereum, userID, chainIDDecimal) {
    // add network to MetaMask
    let chainIdHex = "0x" + chainIDDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'

    try {
        await ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [
                {
                    chainId: chainIdHex,
                    chainName: 'Obscuro Testnet',
                    nativeCurrency: {
                        name: 'Sepolia Ether',
                        symbol: 'ETH',
                        decimals: 18
                    },
                    rpcUrls: [obscuroGatewayAddress+"/"+obscuroGatewayVersion+'/?u='+userID],
                    blockExplorerUrls: ['https://testnet.obscuroscan.io'],
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
    const isAuthenticated = await accountIsAuthenticated(account, userID)

    if (isAuthenticated) {
        return "Account is already authenticated"
    }

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

async function accountIsAuthenticated(account, userID) {
    const queryAccountUserID = pathQuery+"?u="+userID+"&a="+account
    const isAuthenticatedResponse = await fetch(
        queryAccountUserID, {
            method: methodGet,
            headers: jsonHeaders,
        }
    );
    let response = await isAuthenticatedResponse.text();
    let jsonResponseObject = JSON.parse(response);
    return jsonResponseObject.status
}

async function revokeUserID(userID) {
    const queryAccountUserID = pathRevoke+"?u="+userID
    const revokeResponse = await fetch(
        queryAccountUserID, {
            method: methodGet,
            headers: jsonHeaders,
        }
    );
    return revokeResponse.ok
}

const initialize = () => {
    const joinButton = document.getElementById(idJoin);
    const addAccountButton = document.getElementById(idAddAccount);
    const addAllAccountsButton = document.getElementById(idAddAllAccounts);
    const revokeUserIDButton = document.getElementById(idRevokeUserID);
    const statusArea = document.getElementById(idStatus);
    const userIDArea = document.getElementById(idUserID);

    // get ObscuroGatewayUserID from local storage
    let userID = localStorage.getItem("ObscuroGatewayUserID")

    // check if userID exists and has correct type and length (is valid) and display either
    // option to join or to add new account to existing user
    if (isValidUserIDFormat(userID)) {
        userIDArea.innerText = "Your userID is: " + userID
        joinButton.style.display = "none"
        addAccountButton.style.display = "block"
        addAllAccountsButton.style.display = "block"
        revokeUserIDButton.style.display = "block"
    } else {
        joinButton.style.display = "block"
        addAccountButton.style.display = "none"
        revokeUserIDButton.style.display = "none"
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
        statusArea.innerText = "Successfully joined Obscuro Gateway";
        // show users an option to add another account and revoke userID
        addAccountButton.style.display = "block"
        addAllAccountsButton.style.display = "block"
        revokeUserIDButton.style.display = "block"
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
        statusArea.innerText = "\n Authentication status: " + authenticateAccountStatus
    })

    addAllAccountsButton.addEventListener(eventClick, async () => {
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

        for (const account of accounts) {
            let authenticateAccountStatus = await authenticateAccountWithObscuroGateway(ethereum, account, userID)
            statusArea.innerText += "\n Authentication status: " + authenticateAccountStatus + " for account: " + account;
        }
    })

    revokeUserIDButton.addEventListener(eventClick, async () => {
        let result = await revokeUserID(userID);
        if (result) {
            localStorage.removeItem("ObscuroGatewayUserID")
            joinButton.style.display = "block";
            revokeUserIDButton.style.display = "none";
            addAllAccountsButton.style.display = "none";
            userIDArea.innerText = "";
            statusArea.innerText = "Revoking UserID successful. Please remove current network from Metamask."
            addAccountButton.style.display = "none";
        }else{
            statusArea.innerText = "Revoking UserID failed";
        }
    })

}

window.addEventListener(eventDomLoaded, initialize);