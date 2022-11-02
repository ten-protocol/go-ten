"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit";
const typeClick = "click";
const methodPost = "POST";
const jsonKeyHeader = "Header";
const jsonKeyL1Proof = "L1Proof";
const jsonKeyEncryptedTxBlob = "EncryptedTxBlob";

const idNumRollups = "numRollups";
const idNumTxs = "numTxs";
const idRollupTime = "rollupTime";
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
const idTxOne = "txOne";
const idTxTwo = "txTwo";
const idTxThree = "txThree";
const idTxFour = "txFour";
const idTxFive = "txFive";

const pathNumRollups = "/api/numrollups/";
const pathNumTxs = "/api/numtxs/";
const pathRollupTime = "/api/rolluptime/";
const pathLatestRollups = "/api/latestrollups/";
const pathLatestTxs = "/api/latesttxs/";
const pathBlock = "/api/block/";
const pathRollup = "/api/rollup/";
const pathDecryptTxBlob = "/api/decrypttxblob/";

// Updates the displayed stats.
async function updateStats() {
    const numRollupsField = document.getElementById(idNumRollups);
    const numTransactionsField = document.getElementById(idNumTxs);
    const rollupTimeField = document.getElementById(idRollupTime);

    const numRollupsResp = await fetch(pathNumRollups);
    if (numRollupsResp.ok) {
        numRollupsField.innerText = "Total rollups: " + await numRollupsResp.text();
    } else {
        numRollupsField.innerText = "Failed to fetch number of rollups.";
    }

    const numTransactionsResp = await fetch(pathNumTxs);
    if (numTransactionsResp.ok) {
        numTransactionsField.innerText = "Total transactions: " + await numTransactionsResp.text();
    } else {
        numTransactionsField.innerText = "Failed to fetch number of transactions.";
    }

    const rollupTimeResp = await fetch(pathRollupTime);
    if (rollupTimeResp.ok) {
        rollupTimeField.innerText = `Avg. rollup time: ${await rollupTimeResp.text()} secs`;
    } else {
        rollupTimeField.innerText = "Failed to fetch average rollup time.";
    }
}

// Updates the list of latest rollups.
async function updateLatestRollups() {
    const rollupOneField = document.getElementById(idRollupOne);
    const rollupTwoField = document.getElementById(idRollupTwo);
    const rollupThreeField = document.getElementById(idRollupThree);
    const rollupFourField = document.getElementById(idRollupFour);
    const rollupFiveField = document.getElementById(idRollupFive);

    const latestRollupsResp = await fetch(pathLatestRollups);

    if (latestRollupsResp.ok) {
        const latestRollupsJSON = JSON.parse(await latestRollupsResp.text());
        rollupOneField.innerText = latestRollupsJSON[0];
        rollupTwoField.innerText = latestRollupsJSON[1];
        rollupThreeField.innerText = latestRollupsJSON[2];
        rollupFourField.innerText = latestRollupsJSON[3];
        rollupFiveField.innerText = latestRollupsJSON[4];
    } else {
        const errMsg = "Failed to fetch latest rollups.";
        rollupOneField.innerText = errMsg;
        rollupTwoField.innerText = errMsg;
        rollupThreeField.innerText = errMsg;
        rollupFourField.innerText = errMsg;
        rollupFiveField.innerText = errMsg;
    }
}

// Updates the list of latest transactions.
async function updateLatestTxs() {
    const txOneField = document.getElementById(idTxOne);
    const txTwoField = document.getElementById(idTxTwo);
    const txThreeField = document.getElementById(idTxThree);
    const txFourField = document.getElementById(idTxFour);
    const txFiveField = document.getElementById(idTxFive);

    const latestTxsResp = await fetch(pathLatestTxs);

    if (latestTxsResp.ok) {
        const latestTxsJSON = JSON.parse(await latestTxsResp.text());
        txOneField.innerText = latestTxsJSON[0];
        txTwoField.innerText = latestTxsJSON[1];
        txThreeField.innerText = latestTxsJSON[2];
        txFourField.innerText = latestTxsJSON[3];
        txFiveField.innerText = latestTxsJSON[4];
    } else {
        const errMsg = "Failed to fetch latest transactions.";
        txOneField.innerText = errMsg;
        txTwoField.innerText = errMsg;
        txThreeField.innerText = errMsg;
        txFourField.innerText = errMsg;
        txFiveField.innerText = errMsg;
    }
}

// Displays the rollup based on the rollup number or transaction hash that has been clicked on.
async function displayClickedItem(event) {
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
        body: rollupID,
        method: methodPost
    });

    if (!rollupResp.ok) {
        rollupArea.innerText = "Failed to fetch rollup.";
        blockArea.innerText = "Failed to fetch block.";
        decryptedTxsArea.innerText = "Failed to decrypt transaction blob.";
        resultPane.scrollIntoView();
        return;
    }

    const rollupJSON = JSON.parse(await rollupResp.text());
    rollupArea.innerText = JSON.stringify(rollupJSON, null, "\t");

    const blockResp = await fetch(pathBlock, {
        body: rollupJSON[jsonKeyHeader][jsonKeyL1Proof],
        method: methodPost
    });

    if (blockResp.ok) {
        const blockJSON = JSON.parse(await blockResp.text());
        blockArea.innerText = JSON.stringify(blockJSON, null, "\t");
    } else {
        blockArea.innerText = "Failed to fetch block.";
    }

    const encryptedTxBlob = rollupJSON[jsonKeyEncryptedTxBlob]
    const decryptTxBlobResp = await fetch(pathDecryptTxBlob, {
        body: encryptedTxBlob,
        method: methodPost
    });

    if (decryptTxBlobResp.ok) {
        const txBlobJSON = JSON.parse(await decryptTxBlobResp.text());
        decryptedTxsArea.innerText = JSON.stringify(txBlobJSON, null, "\t");
    } else {
        decryptedTxsArea.innerText = "Failed to decrypt transaction blob.";
    }

    resultPane.scrollIntoView();
}

const initialize = () => {
    // Updates the page's stats and latest rollups and transactions.
    setInterval(async () => {
        await updateStats();
        await updateLatestRollups();
        await updateLatestTxs();
    }, 1000);

    // Handles searches for rollups.
    document.getElementById(idFormGetRollup).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();
        const rollupNumber = document.getElementById(idRollupID).value;
        await displayRollup(rollupNumber);
    });

    // Handles clicks on the latest rollups or transactions.
    document.getElementById(idRollupOne).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupTwo).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupThree).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupFour).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idRollupFive).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxOne).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxTwo).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxThree).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxFour).addEventListener(typeClick, displayClickedItem);
    document.getElementById(idTxFive).addEventListener(typeClick, displayClickedItem);
}

window.addEventListener(eventDomLoaded, initialize);