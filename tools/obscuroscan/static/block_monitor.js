"use strict";

const eventDomLoaded = "DOMContentLoaded";
const typeSubmit = "submit";
const idNumRollups = "numRollups";
const idFormGetRollup = "form-get-rollup";
const idRollupNumber = "rollupNumber";
const idBlock = "block";
const idRollup = "rollup";
const pathNumRollups = "/numrollups/";
const pathBlock = "/block/";
const pathRollup = "/rollup/";
const methodPost = "POST";
const jsonKeyHeader = "Header";
const jsonKeyL1Proof = "L1Proof";

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

    document.getElementById(idFormGetRollup).addEventListener(typeSubmit, async (event) => {
        event.preventDefault();

        const rollupNumber = document.getElementById(idRollupNumber).value;
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
    });
}

window.addEventListener(eventDomLoaded, initialize);