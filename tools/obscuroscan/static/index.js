"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit";
const idNumRollups = "numRollups";
const idFormGetRollup = "form-get-rollup";
const idRollupID = "rollupID";
const idBlock = "block";
const idRollup = "rollup";
const idDecryptedTxs = "decryptedTxs";
const pathNumRollups = "/numrollups/";
const pathBlock = "/block/";
const pathRollup = "/rollup/";
const pathDecryptTxBlob = "/decrypttxblob/";
const methodPost = "POST";
const jsonKeyHeader = "Header";
const jsonKeyL1Proof = "L1Proof";
const jsonKeyEncryptedTxBlob = "EncryptedTxBlob";

const initialize = () => {
    const numRollupsField = document.getElementById(idNumRollups);
    const blockArea = document.getElementById(idBlock);
    const rollupArea = document.getElementById(idRollup);
    const decryptedTxsArea = document.getElementById(idDecryptedTxs);

    setInterval(async () => {
        const numRollupsResp = await fetch(pathNumRollups);

        if (numRollupsResp.ok) {
            numRollupsField.innerText = "Total rollups: " + await numRollupsResp.text();
        } else {
            numRollupsField.innerText = "Failed to fetch number of rollups. Cause: " + await numRollupsResp.text()
        }
    }, 1000);

    document.getElementById(idFormGetRollup).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();

        const rollupNumber = document.getElementById(idRollupID).value;
        const rollupResp = await fetch(pathRollup, {
            method: methodPost,
            body: rollupNumber
        });

        if (!rollupResp.ok) {
            rollupArea.innerText = "Failed to fetch rollup. Cause: " + await rollupResp.text()
            blockArea.innerText = "Failed to fetch block. Cause: Could not retrieve rollup."
            return
        }

        const rollupJSON = JSON.parse(await rollupResp.text())
        rollupArea.innerText = JSON.stringify(rollupJSON, null, "\t");

        const blockResp = await fetch(pathBlock, {
            method: methodPost,
            body: rollupJSON[jsonKeyHeader][jsonKeyL1Proof]
        });

        if (blockResp.ok) {
            const json = JSON.parse(await blockResp.text())
            blockArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            blockArea.innerText = "Failed to fetch block. Cause: " + await blockResp.text()
        }

        const encryptedTxBlob = rollupJSON[jsonKeyEncryptedTxBlob]
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
    });
}

window.addEventListener(eventDomLoaded, initialize);