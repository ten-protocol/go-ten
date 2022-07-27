"use strict";

const eventDomLoaded = "DOMContentLoaded";
const idNumBlocks = "numBlocks";
const idBlockHead = "headBlock";
const idHeadRollup = "headRollup";
const pathNumBlocks = "/numblocks/";
const pathHeadBlock = "/headblock/";
const pathHeadRollup = "/headrollup/";

const initialize = () => {
    const numBlocksField = document.getElementById(idNumBlocks);
    const blockHeadArea = document.getElementById(idBlockHead);
    const headRollupArea = document.getElementById(idHeadRollup);

    setInterval(async () => {
        const numBlocksResp = await fetch(pathNumBlocks);

        if (numBlocksResp.ok) {
            numBlocksField.innerText = await numBlocksResp.text();
        } else {
            numBlocksField.innerText = "Failed to fetch number of blocks. Cause: " + await numBlocksResp.text()
        }
    }, 1000);

    setInterval(async () => {
        const headBlockResp = await fetch(pathHeadBlock);

        if (headBlockResp.ok) {
            const json = JSON.parse(await headBlockResp.text())
            blockHeadArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            blockHeadArea.innerText = "Failed to fetch head block. Cause: " + await headBlockResp.text()
        }
    }, 1000);

    setInterval(async () => {
        const headRollupResp = await fetch(pathHeadRollup);

        if (headRollupResp.ok) {
            const json = JSON.parse(await headRollupResp.text())
            headRollupArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            headRollupArea.innerText = "Failed to fetch head rollup. Cause: " + await headRollupResp.text()
        }
    }, 1000);
}

window.addEventListener(eventDomLoaded, initialize);