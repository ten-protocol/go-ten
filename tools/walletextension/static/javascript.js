const initialize = () => {
    const generateViewingKeyButton = document.getElementById('generateViewingKey');

    generateViewingKeyButton.addEventListener('click', async () => {
        const viewingPublicKeyResp = await fetch('/getviewingkey');
        const viewingKey = await viewingPublicKeyResp.text();
        console.log(viewingKey)

        const accounts = await ethereum.request({method: 'eth_requestAccounts'});
        const account = accounts[0]; // todo - allow use of other accounts?
        const signedViewingKey = await ethereum.request({
            method: 'personal_sign',
            // Without a prefix such as 'vk', personal_sign transforms the data for security reasons.
            params: ['vk' + viewingKey, account]
        });
        console.log(signedViewingKey)

        const signedViewingKeyJson = {
            "viewingKey": viewingKey,
            "signature": signedViewingKey
        }
        console.log(signedViewingKeyJson)

        await fetch(
            '/storeviewingkey/', {
                method: 'post',
                headers: {'Content-Type': 'application/json', 'Accept': 'application/json'},
                mode: 'cors',
                body: JSON.stringify(signedViewingKeyJson)
            }
        );

        console.log("success!")
    })
}

window.addEventListener('DOMContentLoaded', initialize);