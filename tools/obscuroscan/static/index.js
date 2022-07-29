"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit";
const typeClick = "click";
const idNumRollups = "numRollups";
const idFormGetRollup = "form-get-rollup";
const idRollupID = "rollupID";
const idResult = "result";
const idBlock = "block";
const idRollup = "rollup";
const idDecryptedTxs = "decryptedTxs";
const idRollupOne = "rollupOne";
const idRollupTwo = "rollupTwo";
const idRollupThree = "rollupThree";
const idRollupFour = "rollupFour";
const idRollupFive = "rollupFive";
const pathNumRollups = "/numrollups/";
const pathLatestRollups = "/latestrollups/";
const pathBlock = "/block/";
const pathRollup = "/rollup/";
const pathDecryptTxBlob = "/decrypttxblob/";
const methodPost = "POST";
const jsonKeyHeader = "Header";
const jsonKeyL1Proof = "L1Proof";
const jsonKeyEncryptedTxBlob = "EncryptedTxBlob";

// Updates the displayed stats.
async function updateStats(numRollupsField) {
    const numRollupsResp = await fetch(pathNumRollups);

    if (numRollupsResp.ok) {
        numRollupsField.innerText = "Total rollups: " + await numRollupsResp.text();
    } else {
        numRollupsField.innerText = "Failed to fetch number of rollups. Cause: " + await numRollupsResp.text()
    }
}

// Updates the list of latest rollups.
async function updateLatestRollups(rollupOneField, rollupTwoField, rollupThreeField, rollupFourField, rollupFiveField) {
    const latestRollupsResp = await fetch(pathLatestRollups);

    if (latestRollupsResp.ok) {
        const latestRollupsJSON = JSON.parse(await latestRollupsResp.text());
        rollupOneField.innerText = latestRollupsJSON[0];
        rollupTwoField.innerText = latestRollupsJSON[1];
        rollupThreeField.innerText = latestRollupsJSON[2];
        rollupFourField.innerText = latestRollupsJSON[3];
        rollupFiveField.innerText = latestRollupsJSON[4];
    } else {
        rollupOneField.innerText = "Failed to fetch latest rollups. Cause: " + await latestRollupsResp.text();
    }
}

// Displays the rollup that has been clicked on.
async function displayClickedRollup(event) {
    event.preventDefault();
    const rollupNumber = document.getElementById(event.target.id).innerText;
    await displayRollup(rollupNumber);
}

// Displays the rollup with the given identifier (either a rollup number or a transaction hash that it contains).
async function displayRollup(rollupID) {
    const resultPane = document.getElementById(idResult);
    const blockArea = document.getElementById(idBlock);
    const rollupArea = document.getElementById(idRollup);
    const decryptedTxsArea = document.getElementById(idDecryptedTxs);

    const rollupResp = await fetch(pathRollup, {
        method: methodPost,
        body: rollupID
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

    resultPane.scrollIntoView();
}

const initialize = () => {
    const numRollupsField = document.getElementById(idNumRollups);

    const rollupOneField = document.getElementById(idRollupOne);
    const rollupTwoField = document.getElementById(idRollupTwo);
    const rollupThreeField = document.getElementById(idRollupThree);
    const rollupFourField = document.getElementById(idRollupFour);
    const rollupFiveField = document.getElementById(idRollupFive);

    setInterval(async () => {
        await updateStats(numRollupsField);
        await updateLatestRollups(rollupOneField, rollupTwoField, rollupThreeField, rollupFourField, rollupFiveField);
    }, 1000);

    document.getElementById(idFormGetRollup).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();
        const rollupNumber = document.getElementById(idRollupID).value;
        await displayRollup(rollupNumber);
    });

    document.getElementById(idRollupOne).addEventListener(typeClick, displayClickedRollup);
    document.getElementById(idRollupTwo).addEventListener(typeClick, displayClickedRollup);
    document.getElementById(idRollupThree).addEventListener(typeClick, displayClickedRollup);
    document.getElementById(idRollupFour).addEventListener(typeClick, displayClickedRollup);
    document.getElementById(idRollupFive).addEventListener(typeClick, displayClickedRollup);
}

window.addEventListener(eventDomLoaded, initialize);