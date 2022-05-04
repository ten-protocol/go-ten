const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idGenerateViewingKey = "generateViewingKey";
const idStatus = "status";
const pathGetViewingKey = "/getviewingkey/";
const pathStoreViewingKey = "/storeviewingkey/";
const methodPost = "post";
const jsonHeaders = {
    "Accept": "application/json",
    "Content-Type": "application/json"
};
const metamaskRequestAccounts = "eth_requestAccounts";
const metamaskPersonalSign = "personal_sign";
const personalSignPrefix = "vk";

const initialize = () => {
    const generateViewingKeyButton = document.getElementById(idGenerateViewingKey);
    const statusArea = document.getElementById(idStatus);

    generateViewingKeyButton.addEventListener(eventClick, async () => {
        const viewingPublicKeyResp = await fetch(pathGetViewingKey); // todo - handle failure of request
        const viewingKey = await viewingPublicKeyResp.text();

        const accounts = await ethereum.request({method: metamaskRequestAccounts}); // todo - handle failure of request
        const account = accounts[0]; // todo - allow use of other accounts?
        const signedBytes = await ethereum.request({
            method: metamaskPersonalSign,
            // Without a prefix such as 'vk', personal_sign transforms the data for security reasons.
            params: [personalSignPrefix + viewingKey, account]
        }); // todo - handle failure of request

        const signedViewingKeyJson = {
            "viewingKey": viewingKey,
            "signedBytes": signedBytes
        }

        const resp = await fetch(
            pathStoreViewingKey, {
                method: methodPost,
                headers: jsonHeaders,
                body: JSON.stringify(signedViewingKeyJson)
            }
        ); // todo - handle failure of request

        if (resp.status >= 200 && resp.status < 300) {
            statusArea.innerText = `Account: ${account}\nViewing key: ${viewingKey}\nSigned bytes: ${signedBytes}`
        }
    })
}

window.addEventListener(eventDomLoaded, initialize);