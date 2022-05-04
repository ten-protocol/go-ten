const initialize = () => {
    const generateViewingKeyButton = document.getElementById('generateViewingKey');
    const statusArea = document.getElementById('status');

    generateViewingKeyButton.addEventListener('click', async () => {
        const viewingPublicKeyResp = await fetch('/getviewingkey'); // todo - handle failure of request
        const viewingKey = await viewingPublicKeyResp.text();

        const accounts = await ethereum.request({method: 'eth_requestAccounts'}); // todo - handle failure of request
        const account = accounts[0]; // todo - allow use of other accounts?
        const signedBytes = await ethereum.request({
            method: 'personal_sign',
            // Without a prefix such as 'vk', personal_sign transforms the data for security reasons.
            params: ['vk' + viewingKey, account]
        }); // todo - handle failure of request

        const signedViewingKeyJson = {
            "viewingKey": viewingKey,
            "signedBytes": signedBytes
        }

        const resp = await fetch(
            '/storeviewingkey/', {
                method: 'post',
                headers: {'Content-Type': 'application/json', 'Accept': 'application/json'},
                mode: 'cors',
                body: JSON.stringify(signedViewingKeyJson)
            }
        ); // todo - handle failure of request

        if (resp.status >= 200 && resp.status < 300) {
            statusArea.innerText = `Account: ${account}\nViewing key: ${viewingKey}\nSigned bytes: ${signedBytes}`
        }
    })
}

window.addEventListener('DOMContentLoaded', initialize);