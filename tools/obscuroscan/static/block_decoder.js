"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit"
const idFormDecryptRollup = "form-decrypt-rollup"
const idDecryptedRollup = "decryptedRollup";
const idEncryptedRollup = "encryptedRollup";
const pathDecryptRollup = "/decryptrollup/";
const methodPost = "POST";

const initialize = () => {
    const decryptedRollupArea = document.getElementById(idDecryptedRollup);

    document.getElementById(idFormDecryptRollup).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();

        const encryptedRollup = document.getElementById(idEncryptedRollup).value
        const decryptRollupResp = await fetch(pathDecryptRollup, {
            method: methodPost,
            body: encryptedRollup
        });

        if (decryptRollupResp.ok) {
            const json = JSON.parse(await decryptRollupResp.text())
            decryptedRollupArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            decryptedRollupArea.innerText = "Failed to decrypt rollup. Cause: " + await decryptRollupResp.text()
        }
    })
}

window.addEventListener(eventDomLoaded, initialize);