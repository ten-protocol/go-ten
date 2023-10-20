const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idJoin = "join";
const idRevokeUserID = "revokeUserID";
const idStatus = "status";
const idAccountsTable = "accountsTable";
const idTableBody = "tableBody";
const idInformation = "information";
const idInformation2 = "information2";
const idBegin = "begin-box";
const idSpinner = "spinner";
const obscuroGatewayVersion = "v1";
const pathJoin = obscuroGatewayVersion + "/join/";
const pathAuthenticate = obscuroGatewayVersion + "/authenticate/";
const pathQuery = obscuroGatewayVersion + "/query/";
const pathRevoke = obscuroGatewayVersion + "/revoke/";
const pathVersion = "/version/";
const obscuroChainIDDecimal = 443;
const methodPost = "post";
const methodGet = "get";
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};

const metamaskPersonalSign = "personal_sign";
const obscuroChainIDHex = "0x" + obscuroChainIDDecimal.toString(16); // Convert to hexadecimal and prefix with '0x'

function isValidUserIDFormat(value) {
    return typeof value === 'string' && value.length === 64;
}

let obscuroGatewayAddress = window.location.protocol + "//" + window.location.host;

let provider = null;

async function fetchAndDisplayVersion() {
    try {
        const versionResp = await fetch(
            pathVersion, {
                method: methodGet,
                headers: jsonHeaders,
            }
        );
        if (!versionResp.ok) {
            throw new Error("Failed to fetch the version");
        }

        let response = await versionResp.text();

        const versionDiv = document.getElementById("versionDisplay");
        versionDiv.textContent = "Version: " + response;
    } catch (error) {
        console.error("Error fetching the version:", error);
    }
}

function getNetworkName(gatewayAddress) {
    switch(gatewayAddress) {
        case 'https://uat-testnet.obscu.ro':
            return 'Obscuro UAT-Testnet';
        case 'https://dev-testnet.obscu.ro':
            return 'Obscuro Dev-Testnet';
        default:
            return 'Obscuro Testnet';
    }
}

function getRPCFromUrl(gatewayAddress) {
    // get the correct RPC endpoint for each network
    switch(gatewayAddress) {
        // case 'https://testnet.obscu.ro':
        //     return 'https://rpc.sepolia-testnet.obscu.ro'
        case 'https://sepolia-testnet.obscu.ro':
            return 'https://rpc.sepolia-testnet.obscu.ro'
        case 'https://uat-testnet.obscu.ro':
            return 'https://rpc.uat-testnet.obscu.ro';
        case 'https://dev-testnet.obscu.ro':
            return 'https://rpc.dev-testnet.obscu.ro';
        default:
            return gatewayAddress;
    }
}

async function addNetworkToMetaMask(ethereum, userID, chainIDDecimal) {
    // add network to MetaMask
    try {
        await ethereum.request({
            method: 'wallet_addEthereumChain',
            params: [
                {
                    chainId: obscuroChainIDHex,
                    chainName: getNetworkName(obscuroGatewayAddress),
                    nativeCurrency: {
                        name: 'Sepolia Ether',
                        symbol: 'ETH',
                        decimals: 18
                    },
                    rpcUrls: [getRPCFromUrl(obscuroGatewayAddress)+"/"+obscuroGatewayVersion+'/?u='+userID],
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
        if (await isObscuroChain()) {
            return await provider.send('eth_getStorageAt', ["getUserID", getRandomIntAsString(0, 1000), null])
        } else {
            return null
        }
    }catch (e) {
        console.log(e)
        return null;
    }
}

async function connectAccounts() {
    try {
        return await window.ethereum.request({ method: 'eth_requestAccounts' });
    } catch (error) {
        // TODO: Display warning to user to allow it and refresh page...
        console.error('User denied account access:', error);
        return null;
    }
}

async function isMetamaskConnected() {
    let accounts;
    try {
        accounts = await provider.listAccounts()
        return accounts.length > 0;

    } catch (error) {
        console.log("Unable to get accounts")
    }
    return false
}

// Check if Metamask is available on mobile or as a plugin in browser
// (https://docs.metamask.io/wallet/how-to/integrate-with-mobile/)
function checkIfMetamaskIsLoaded() {
    if (window.ethereum) {
        handleEthereum();
    } else {
        const statusArea = document.getElementById(idStatus);
        const table = document.getElementById("accountsTable");
        table.style.display = "none";
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

async function isObscuroChain() {
    let currentChain = await ethereum.request({ method: 'eth_chainId' });
    return currentChain === obscuroChainIDHex
}

async function switchToObscuroNetwork() {
    try {
        await ethereum.request({
            method: 'wallet_switchEthereumChain',
            params: [{ chainId: obscuroChainIDHex }],
        });
        return 0
    } catch (switchError) {
        return switchError.code
    }
    return -1
}

const initialize = async () => {
    const joinButton = document.getElementById(idJoin);
    const revokeUserIDButton = document.getElementById(idRevokeUserID);
    const statusArea = document.getElementById(idStatus);
    const informationArea = document.getElementById(idInformation);
    const informationArea2 = document.getElementById(idInformation2);
    const beginBox = document.getElementById(idBegin);
    const spinner = document.getElementById(idSpinner);

    const accountsTable = document.getElementById(idAccountsTable)
    const tableBody = document.getElementById(idTableBody);

    // getUserID from the gateway with getStorageAt method
    let userID = await getUserID()

    function displayOnlyJoin() {
        joinButton.style.display = "block"
        revokeUserIDButton.style.display = "none"
        accountsTable.style.display = "none"
        informationArea.style.display = "block"
        informationArea2.style.display = "none"

        beginBox.style.visibility = "visible";
        spinner.style.visibility = "hidden";
    }

    async function displayConnectedAndJoinedSuccessfully() {
        joinButton.style.display = "none"
        informationArea.style.display = "none"
        informationArea2.style.display = "block"
        revokeUserIDButton.style.display = "block"
        accountsTable.style.display = "block"
        
        await populateAccountsTable(document, tableBody, userID)
    }

    async function displayCorrectScreenBasedOnMetamaskAndUserID() {
        // check if we are on Obscuro Chain
        if(await isObscuroChain()){
            // check if we have valid userID in rpcURL
            if (isValidUserIDFormat(userID)) {
                return await displayConnectedAndJoinedSuccessfully()
            }
        }
        return displayOnlyJoin()
    }

    // load the current version
    await fetchAndDisplayVersion();

    await displayCorrectScreenBasedOnMetamaskAndUserID()

    joinButton.addEventListener(eventClick, async () => {
        // clean up any previous errors
        statusArea.innerText = ""
        // check if we are on an obscuro chain
        if (await isObscuroChain()) {
            userID = await getUserID()
            if (!isValidUserIDFormat(userID)) {
                statusArea.innerText = "Existing Obscuro network detected in MetaMask. Please remove before hitting begin"
            }
        } else {
            // we are not on an Obscuro network - try to switch
            let switched = await switchToObscuroNetwork();
            // error 4902 means that the chain does not exist
            if (switched === 4902 || !isValidUserIDFormat(await getUserID())) {
                // join the network
                const joinResp = await fetch(
                    pathJoin, {
                        method: methodGet,
                        headers: jsonHeaders,
                    });
                if (!joinResp.ok) {
                    console.log("Error joining Obscuro Gateway")
                    statusArea.innerText = "Error joining Obscuro Gateway. Please try again later."
                    return
                }
                userID = await joinResp.text();

                // add Obscuro network
                await addNetworkToMetaMask(window.ethereum, userID)
            }

            // we have to check if user has accounts connected with metamask - and promt to connect if not
            if (!await isMetamaskConnected()) {
                await connectAccounts();
            }

            // connect all accounts
            // Get an accounts and prompt user to sign joining with a selected account
            const accounts = await provider.listAccounts();
            if (accounts.length === 0) {
                statusArea.innerText = "No MetaMask accounts found."
                return
            }

            userID = await getUserID();
            beginBox.style.visibility = "hidden";
            spinner.style.visibility = "visible";
            for (const account of accounts) {
                await authenticateAccountWithObscuroGateway(ethereum, account, userID)
                accountsTable.style.display = "block"
                await populateAccountsTable(document, tableBody, userID)
            }
            beginBox.style.visibility = "visible";
            spinner.style.visibility = "hidden";

            // if accounts change we want to give user chance to add them to Obscuro
            window.ethereum.on('accountsChanged', async function (accounts) {
                if (isValidUserIDFormat(await getUserID())) {
                    userID = await getUserID();
                    for (const account of accounts) {
                        await authenticateAccountWithObscuroGateway(ethereum, account, userID)
                        accountsTable.style.display = "block"
                        await populateAccountsTable(document, tableBody, userID)
                    }
                }
            });

            await displayConnectedAndJoinedSuccessfully()
        }
    })

    revokeUserIDButton.addEventListener(eventClick, async () => {
        beginBox.style.visibility = "hidden";
        spinner.style.visibility = "visible";
        let result = await revokeUserID(userID);

        await populateAccountsTable(document, tableBody, userID)

        if (result) {
            displayOnlyJoin()
        } else {
            statusArea.innerText = "Revoking UserID failed";
        }
    })
    beginBox.style.visibility = "visible";
    spinner.style.visibility = "hidden";
}

window.addEventListener(eventDomLoaded, checkIfMetamaskIsLoaded);

