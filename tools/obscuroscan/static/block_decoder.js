"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit"
const idFormDecryptTxBlob = "form-decrypt-tx-blob"
const idDecryptedTxs = "decryptedTxs";
const idEncryptedTxBlob = "encryptedTxBlob";
const pathDecryptTxBlob = "/decrypttxblob/";
const methodPost = "POST";

const initialize = () => {
    const decryptedTxsArea = document.getElementById(idDecryptedTxs);

    document.getElementById(idFormDecryptTxBlob).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();

        const encryptedTxBlob = document.getElementById(idEncryptedTxBlob).value
        const decryptTxBlobResp = await fetch(pathDecryptTxBlob, {
            method: methodPost,
            body: encryptedTxBlob
        });

        if (decryptTxBlobResp.ok) {
            const json = JSON.parse(await decryptTxBlobResp.text())
            decryptedTxsArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            decryptedTxsArea.innerText = "Failed to decrypt transaction blob. Cause: " + await decryptTxBlobResp.text()
        }
    })
}

window.addEventListener(eventDomLoaded, initialize);