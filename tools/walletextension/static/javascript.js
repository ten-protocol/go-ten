const initialize = () => {
    const viewingKey = document.getElementById('viewingKey');
    const generateViewingKeyButton = document.getElementById('generateViewingKey');
    const signedViewingKey = document.getElementById('signedViewingKey');
    const signViewingKeyButton = document.getElementById('signViewingKey');
    const storeViewingKeyButton = document.getElementById('storeViewingKey');

    // todo - disable button after click
    generateViewingKeyButton.addEventListener('click', async () => {
        const viewingPublicKeyResp = await fetch('/getviewingkey');
        viewingKey.innerText = await viewingPublicKeyResp.text();
    })

    // todo - disable button after click
    signViewingKeyButton.addEventListener('click', async () => {
        const accounts = await ethereum.request({method: 'eth_requestAccounts'});
        const account = accounts[0];
        signedViewingKey.innerText = await ethereum.request({
            method: 'personal_sign',
            params: [viewingKey.innerText, account]
        });
    });

    // todo - disable button after click
    storeViewingKeyButton.addEventListener('click', async () => {
        const signedViewingKeyJson = {
            "viewingKey": viewingKey.innerText,
            "signature": signedViewingKey.innerText
        }
        console.log(JSON.stringify(signedViewingKeyJson))
        await fetch(
            '/storeviewingkey/', {
                method: 'post',
                headers: {'Content-Type': 'application/json', 'Accept': 'application/json'},
                mode: 'cors',
                body: JSON.stringify(signedViewingKeyJson)
            }
        );
    })
}

window.addEventListener('DOMContentLoaded', initialize);