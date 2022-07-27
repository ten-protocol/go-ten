"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit";
const idNumRollups = "numRollups";
const idFormGetBlock = "form-get-block";
const idBlockNumber = "blockNumber";
const idBlock = "block";
const idRollup = "rollup";
const pathNumRollups = "/numrollups/";
const pathBlock = "/block/";
const pathRollup = "/rollup/";
const methodPost = "POST";

const initialize = () => {
    const numRollupsField = document.getElementById(idNumRollups);
    const blockArea = document.getElementById(idBlock);
    const rollupArea = document.getElementById(idRollup);

    setInterval(async () => {
        const numRollupsResp = await fetch(pathNumRollups);

        if (numRollupsResp.ok) {
            numRollupsField.innerText = await numRollupsResp.text();
        } else {
            numRollupsField.innerText = "Failed to fetch number of rollups. Cause: " + await numRollupsResp.text()
        }
    }, 1000);

    document.getElementById(idFormGetBlock).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();

        const number = document.getElementById(idBlockNumber).value;

        const blockResp = await fetch(pathBlock, {
            method: methodPost,
            body: number
        });

        if (blockResp.ok) {
            const json = JSON.parse(await blockResp.text())
            blockArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            blockArea.innerText = "Failed to fetch block. Cause: " + await blockResp.text()
        }

        const rollupResp = await fetch(pathRollup, {
            method: methodPost,
            body: number
        });

        if (rollupResp.ok) {
            const json = JSON.parse(await rollupResp.text())
            rollupArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            rollupArea.innerText = "Failed to fetch rollup. Cause: " + await rollupResp.text()
        }
    });
}

window.addEventListener(eventDomLoaded, initialize);