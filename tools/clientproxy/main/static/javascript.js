const initialize = () => {
    const generateViewingKeyButton = document.getElementById('generateViewingKey');
    const viewingKey = document.getElementById('viewingKey')
    const signViewingKeyButton = document.getElementById('signViewingKey');
    const signedViewingKey = document.getElementById('signedViewingKey')

    // todo - disable button after click
    generateViewingKeyButton.addEventListener('click', async () => {
        const viewingPublicKeyResp = await fetch('/generateViewingKeyPair');
        const viewingPublicKey = await viewingPublicKeyResp.text()
        viewingKey.innerText = viewingPublicKey
    })

    // todo - disable button after click
    signViewingKeyButton.addEventListener('click', async () => {
        const accounts = await ethereum.request({method: 'eth_requestAccounts'});
        const account = accounts[0]
        const signedMessage = await ethereum.request({
            method: 'personal_sign',
            params: [viewingKey.innerText, account]
        });
        signedViewingKey.innerText = `account ${account} signed viewing key: ${signedMessage}`
    });

    // todo - additional button to submit signed viewing key
}

window.addEventListener('DOMContentLoaded', initialize);