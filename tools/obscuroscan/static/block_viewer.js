"use strict";

const eventDomLoaded = "DOMContentLoaded";
const idNumBlocks = "numBlocks";
const idBlock = "block";
const idRollup = "rollup";
const pathNumBlocks = "/numblocks/";
const pathBlock = "/block/";
const pathRollup = "/rollup/";

const initialize = () => {
    const numBlocksField = document.getElementById(idNumBlocks);
    const blockArea = document.getElementById(idBlock);
    const rollupArea = document.getElementById(idRollup);

    setInterval(async () => {
        const numBlocksResp = await fetch(pathNumBlocks);

        if (numBlocksResp.ok) {
            numBlocksField.innerText = await numBlocksResp.text();
        } else {
            numBlocksField.innerText = "Failed to fetch number of blocks. Cause: " + await numBlocksResp.text()
        }
    }, 1000);

    setInterval(async () => {
        const blockResp = await fetch(pathBlock);

        if (blockResp.ok) {
            const json = JSON.parse(await blockResp.text())
            blockArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            blockArea.innerText = "Failed to fetch block. Cause: " + await blockResp.text()
        }
    }, 1000);

    setInterval(async () => {
        const rollupResp = await fetch(pathRollup);

        if (rollupResp.ok) {
            const json = JSON.parse(await rollupResp.text())
            rollupArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            rollupArea.innerText = "Failed to fetch rollup. Cause: " + await rollupResp.text()
        }
    }, 1000);
}

window.addEventListener(eventDomLoaded, initialize);