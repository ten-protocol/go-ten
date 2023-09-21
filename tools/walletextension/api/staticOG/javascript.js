const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idJoin = "join";
const idAddAccount = "addAccount";
const idAddAllAccounts = "addAllAccounts";
const idRevokeUserID = "revokeUserID";
const idStatus = "status";
const obscuroGatewayVersion = "v1"
const pathJoin = obscuroGatewayVersion + "/join/";
const pathAuthenticate = obscuroGatewayVersion + "/authenticate/";
const pathQuery = obscuroGatewayVersion + "/query/";
const pathRevoke = obscuroGatewayVersion + "/revoke/";
const obscuroChainIDDecimal = 777;
const methodPost = "post";
const methodGet = "get";
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};

const metamaskPersonalSign = "personal_sign";

function isValidUserIDFormat(value) {
    return typeof value === 'string' && value.length === 64;
}

let obscuroGatewayAddress = window.location.protocol + "//" + window.location.host;

let provider = null;


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
                        name: 'Obscuro',
                        symbol: 'OBX',
                        decimals: 18
                    },
                    rpcUrls: [obscuroGatewayAddress+"/"+obscuroGatewayVersion+'/?u='+userID],
                    blockExplorerUrls: null
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

    const textToSign = "Register " + userID + " for " + account.toLowerCase();
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

function getRandomIntAsString(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    const randomInt = Math.floor(Math.random() * (max - min + 1)) + min;
    return randomInt.toString();
}


async function getUserID() {
    try {
        return await provider.send('eth_getStorageAt', ["getUserID", getRandomIntAsString(0, 1000), null])
    }catch (e) {
        console.log(e)
        return null;
    }
}

async function connectAccount() {
    try {
        return await window.ethereum.request({ method: 'eth_requestAccounts' });
    } catch (error) {
        // TODO: Display warning to user to allow it and refresh page...
        console.error('User denied account access:', error);
        return null;
    }
}

// Check if Metamask is available on mobile or as a plugin in browser
// (https://docs.metamask.io/wallet/how-to/integrate-with-mobile/)
function checkIfMetamaskIsLoaded() {
    if (window.ethereum) {
        handleEthereum();
    } else {
        // TODO: Refactor and change the way we hide and display items on our webpage
        document.getElementById(idJoin).style.display = "none";
        document.getElementById(idAddAccount).style.display = "none";
        document.getElementById(idAddAllAccounts).style.display = "none";
        document.getElementById(idRevokeUserID).style.display = "none";
        const statusArea = document.getElementById(idStatus);
        statusArea.innerText = 'Connecting to Metamask...';
        window.addEventListener('ethereum#initialized', handleEthereum, {
            once: true,
        });

        // If the event is not dispatched by the end of the timeout,
        // the user probably doesn't have MetaMask installed.
        setTimeout(handleEthereum, 3000); // 3 seconds
    }
}

function handleEthereum() {
    const { ethereum } = window;
    if (ethereum && ethereum.isMetaMask) {
        provider = new ethers.providers.Web3Provider(window.ethereum);
        initialize()
    } else {
        const statusArea = document.getElementById(idStatus);
        statusArea.innerText = 'Please install MetaMask to use Obscuro Gateway.';
    }
}

async function populateAccountsTable(document, tableBody, userID) {
    tableBody.innerHTML = '';
    const accounts = await provider.listAccounts();
    for (const account of accounts) {
        const row = document.createElement('tr');

        const accountCell = document.createElement('td');
        accountCell.textContent = account;
        row.appendChild(accountCell);

        const statusCell = document.createElement('td');

        statusCell.textContent = await accountIsAuthenticated(account, userID);  // Status is empty for now
        row.appendChild(statusCell);

        tableBody.appendChild(row);
    }
}

const initialize = async () => {
    const joinButton = document.getElementById(idJoin);
    const addAccountButton = document.getElementById(idAddAccount);
    const addAllAccountsButton = document.getElementById(idAddAllAccounts);
    const revokeUserIDButton = document.getElementById(idRevokeUserID);
    const statusArea = document.getElementById(idStatus);

    const accountsTable = document.getElementById('accountsTable')
    const tableBody = document.getElementById('tableBody');
    // getUserID from the gateway with getStorageAt method
    let userID = await getUserID()

    // check if userID exists and has a correct type and length (is valid) and display either
    // option to join or to add a new account to existing user
    if (isValidUserIDFormat(userID)) {
        joinButton.style.display = "none"
        addAccountButton.style.display = "block"
        addAllAccountsButton.style.display = "block"
        revokeUserIDButton.style.display = "block"
        accountsTable.style.display = "block"
        await populateAccountsTable(document, tableBody, userID)
    } else {
        joinButton.style.display = "block"
        addAccountButton.style.display = "none"
        revokeUserIDButton.style.display = "none"
        accountsTable.style.display = "none"
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
        accountsTable.style.display = "block"
        await populateAccountsTable(document, tableBody, userID)
    })

    addAccountButton.addEventListener(eventClick, async () => {
        // check if we have userID and it is the correct length
        if (!isValidUserIDFormat(userID)) {
            statusArea.innerText = "\n Please join Obscuro network first"
            joinButton.style.display = "block"
            addAccountButton.style.display = "none"
        }

        await connectAccount()

        // Get an account and prompt user to sign joining with a selected account
        const account = await provider.getSigner().getAddress();
        if (account.length === 0) {
            statusArea.innerText = "No MetaMask accounts found."
            return
        }
        let authenticateAccountStatus = await authenticateAccountWithObscuroGateway(ethereum, account, userID)
        //statusArea.innerText = "\n Authentication status: " + authenticateAccountStatus
        accountsTable.style.display = "block"
        await populateAccountsTable(document, tableBody, userID)
    })

    addAllAccountsButton.addEventListener(eventClick, async () => {
        // check if we have userID and it is the correct length
        if (!isValidUserIDFormat(userID)) {
            statusArea.innerText = "\n Please join Obscuro network first"
            joinButton.style.display = "block"
            addAccountButton.style.display = "none"
        }

        await connectAccount()

        // Get an account and prompt user to sign joining with selected account
        const accounts = await provider.listAccounts();
        if (accounts.length === 0) {
            statusArea.innerText = "No MetaMask accounts found."
            return
        }

        for (const account of accounts) {
            let authenticateAccountStatus = await authenticateAccountWithObscuroGateway(ethereum, account, userID)
            accountsTable.style.display = "block"
            await populateAccountsTable(document, tableBody, userID)
        }
    })

    revokeUserIDButton.addEventListener(eventClick, async () => {
        let result = await revokeUserID(userID);

        await populateAccountsTable(document, tableBody, userID)

        if (result) {
            joinButton.style.display = "block";
            revokeUserIDButton.style.display = "none";
            addAllAccountsButton.style.display = "none";
            statusArea.innerText = "Revoking UserID successful. Please remove current network from Metamask."
            addAccountButton.style.display = "none";
            accountsTable.style.display = "none"
        } else {
            statusArea.innerText = "Revoking UserID failed";
        }
    })

}

window.addEventListener(eventDomLoaded, checkIfMetamaskIsLoaded);