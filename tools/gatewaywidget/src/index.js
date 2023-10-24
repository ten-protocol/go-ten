
// Embedding Obscuro Gateway Widget Styles
(function() {
    const styles = `/* Obscuro Gateway Widget CSS Code */
#obscuro-button {
    border: none;
    cursor: pointer;
    outline: none;
    position: fixed;
    bottom: 20px;
    right: 20px;
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background-color: #1D1D1D;
    color: #fff;
    font-size: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0px 0px 7.5px 1px rgba(0,0,0,0.75);
    transition: transform 0.3s ease;
    z-index: 9999
}

#obscuro-button:hover {
    transform: scale(1.2);
}

/* Obscuro Gateway Panel */
#obscuro-panel {
    box-shadow: 0px 0px 10px 1px rgba(0,0,0,0.75);
    border-radius: 15px;
    border: 1px solid #ffffff9f; 
    background-color: #1D1D1D;
    position: fixed; 
    bottom: 100px;
    right: 20px;
    width: 290px;
    max-height: 600px;
    overflow-y: auto;
    z-index: 9998;
    padding: 20px;
    flex-direction: column;
    font-family: 'Onest', sans-serif;
}

#options {
    position: absolute;
    cursor: pointer;
    right: 20px;
    border-radius: 25%;
    font-size: 20px;
    width: 25px;
    height: 25px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #d01d38;
    color: #fff;
    border: 1px solid #d01d38;
    transition: background-color 0.3s ease;
}

#options:hover {
    background-color: #ff6262;
    border: 1px solid #ff6262;
}

#options-menu {
    position: absolute;
    top: 80px; /* Adjust this as needed */
    right: 15px;
    background-color: #1D1D1D;
    color: white;
    box-shadow: 0px 0px 5px 0.5px rgba(0,0,0,0.75);
    z-index: 10000; 
    width: 150px; /* Or whatever width you desire */
    cursor: pointer;
}

#options-menu ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

#options-menu li {
    padding: 10px;
    border-bottom: 1px solid #353535;
}

#options-menu li:last-child {
    border-bottom: none;
}

h3 {
    color: #fff;
}

hr {
    border: 0;
    border-top: 1px solid #353535; /* Dark line */
    margin-bottom: 15px;
}

button {
    margin: 5px;
    padding: 10px 15px;
    background-color: #00a8e6;
    color: #fff;
    border: 1px solid #00a8e6;
    border-radius: 7.5px;
    transition: background-color 0.3s ease;
    width: 135px;
    display: block; /* Makes the button behave like a block-level element */
    margin-left: auto; /* Centers the button */
    margin-right: auto;
    cursor: pointer;
}

button:hover {
    background-color: #555;
    border: 1px solid #00a8e6;
}

#accountsTable {
    width: 100%;
    margin-top: 20px;
    border-collapse: collapse; 
    border-spacing: 0; 
    border-radius: 10px; 
    overflow: hidden;
}

#accountsTable th, #accountsTable td {
    padding: 10px;
    color: #fff;
    background-color: #222222;
    text-align: center;
    border: 0.01em solid #353535;
}

#accountsTable th {
    background-color: #000000;
    color: white;
}

#status {
    margin-top: 10px;
    color: #ffae00;
    font-size: 15px;
    text-align: center;
}

.hidden {
    display: none;
}

.close-logo {
    font-size: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
}
`;
    
    const styleElement = document.createElement('style');
    styleElement.type = 'text/css';
    styleElement.innerHTML = styles;

    document.head.appendChild(styleElement);
})();



// Automatically embed the Obscuro Gateway Widget upon script execution
(function() {
    const widgetHTML = `<!-- Obscuro Gateway Widget HTML Code -->
    <button id="obscuro-button">
        <span class="open-logo">◠.</span>
        <span class="close-logo hidden">X</span>
    </button>

    <!-- Obscuro Gateway Panel -->
    <div id="obscuro-panel" class="hidden">
        <div>
            <!-- Options Menu Button -->
            <div id="options">⫶</div>

            <!-- Dropdown Menu Content -->
            <div id="options-menu" class="hidden">
                <ul>
                    <li><div id="revokeUserID">Revoke UserID</div></li>
                </ul>
            </div>
            <h3>Obscuro Gateway</h3>
            <hr>

            <!-- Join Obscuro Network -->
            <button id="join">JOIN OBSCURO</button>

            <!-- Status Area -->
            <div id="status"></div>

            <!-- Accounts Table -->
            <table id="accountsTable">
                <thead>
                    <tr>
                        <th>Account</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody id="tableBody"></tbody>
            </table>
        </div>
    </div>

    <!-- Add necessary styles for the hidden class -->
    <style>
        .hidden {
            display: none;
        }
    </style>`;
    
    // Create a new div for the widget and set its innerHTML
    const widgetContainer = document.createElement('div');
    widgetContainer.innerHTML = widgetHTML;

    // Append the widget container to the document body
    document.body.appendChild(widgetContainer);
})();


const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idJoin = "join";
const idOptionsButton = "options";
const idOptionsMenu = "options-menu"
const idRevokeUserID = "revokeUserID";
const idStatus = "status";
const idAccountsTable = "accountsTable";
const idTableBody = "tableBody";
const obscuroGatewayVersion = "v1";
const obscuroGatewayAddress = "https://testnet.obscu.ro";
const pathJoin = obscuroGatewayAddress + "/" + obscuroGatewayVersion + "/join/";
const pathAuthenticate = obscuroGatewayAddress + "/" + obscuroGatewayVersion + "/authenticate/";
const pathQuery = obscuroGatewayAddress + "/" + obscuroGatewayVersion + "/query/";
const pathRevoke = obscuroGatewayAddress + "/" + obscuroGatewayVersion + "/revoke/";
const pathVersion = obscuroGatewayAddress + "/version/";
const obscuroChainIDDecimal = 443;
const methodPost = "post";
const methodGet = "get";
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};

const metamaskPersonalSign = "personal_sign";
const obscuroChainIDHex = "0x" + obscuroChainIDDecimal.toString(16);

function isValidUserIDFormat(value) {
    return typeof value === 'string' && value.length === 64;
}

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
        accountCell.textContent = account.slice(0, 22);
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

// Widget UI interactions
function toggleObscuroPanel() {
    const panel = document.getElementById('obscuro-panel');
    if (panel.classList.contains('hidden')) {
        panel.classList.remove('hidden');
    } else {
        panel.classList.add('hidden');
    }
}

document.getElementById("options").addEventListener("click", function() {
    const optionsMenu = document.getElementById("options-menu");
    if (optionsMenu.style.display === "none" || optionsMenu.style.display === "") {
        optionsMenu.style.display = "block";
    } else {
        optionsMenu.style.display = "none";
    }
});


document.getElementById('obscuro-button').addEventListener('click', function() {
    const panel = document.getElementById('obscuro-panel');
    const openLogo = document.querySelector('.open-logo');
    const closeLogo = document.querySelector('.close-logo');

    // Toggle the panel's visibility
    panel.classList.toggle('hidden');

    // Switch between the ◠. and X logos
    openLogo.classList.toggle('hidden');
    closeLogo.classList.toggle('hidden');
});


const initialize = async () => {
    const joinButton = document.getElementById(idJoin);
    const revokeUserIDButton = document.getElementById(idRevokeUserID);
    const statusArea = document.getElementById(idStatus);
    const optionsButton = document.getElementById(idOptionsButton);
    const optionsMenu = document.getElementById(idOptionsMenu);
    const accountsTable = document.getElementById(idAccountsTable)
    const tableBody = document.getElementById(idTableBody);
    // getUserID from the gateway with getStorageAt method
    let userID = await getUserID()

    function displayOnlyJoin() {
        joinButton.style.display = "block"
        revokeUserIDButton.style.display = "none"
        optionsButton.style.display = "none";
        optionsMenu.style.display = "none";
        accountsTable.style.display = "none"
    }

    async function displayConnectedAndJoinedSuccessfully() {
        joinButton.style.display = "none"
        revokeUserIDButton.style.display = "block"
        optionsButton.style.display = "flex";
        optionsMenu.style.display = "none";
        accountsTable.style.display = "block"
        await populateAccountsTable(document, tableBody, userID)
        statusArea.innerText = "Successfully connected";
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
        // check if we are on an obscuro chain
        if (await isObscuroChain()) {
            userID = await getUserID()
            if (!isValidUserIDFormat(userID)) {
                statusArea.innerText = "Please remove existing Obscuro network from metamask and start again."
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
            for (const account of accounts) {
                await authenticateAccountWithObscuroGateway(ethereum, account, userID)
                accountsTable.style.display = "block"
                await populateAccountsTable(document, tableBody, userID)
            }

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
        let result = await revokeUserID(userID);

        await populateAccountsTable(document, tableBody, userID)

        if (result) {
            displayOnlyJoin()
            statusArea.innerText = "Revoked UserID";
        } else {
            statusArea.innerText = "Revoking UserID failed";
        }
    })
}

window.addEventListener(eventDomLoaded, checkIfMetamaskIsLoaded);
