"use strict";

const eventDomLoaded = "DOMContentLoaded";
const idBlockHead = "headBlock";
const idHeadRollup = "headRollup";
const pathHeadBlock = "/headblock/";
const pathHeadRollup = "/headrollup/";

const initialize = () => {
    const blockHeadArea = document.getElementById(idBlockHead);
    const headRollupArea = document.getElementById(idHeadRollup);

    setInterval(async () => {
        const headBlockResp = await fetch(pathHeadBlock);

        if (headBlockResp.ok) {
            const json = JSON.parse(await headBlockResp.text())
            blockHeadArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            blockHeadArea.innerText = "Failed to fetch head block."
        }
    }, 1000);

    setInterval(async () => {
        const headRollupResp = await fetch(pathHeadRollup);

        if (headRollupResp.ok) {
            const json = JSON.parse(await headRollupResp.text())
            headRollupArea.innerText = JSON.stringify(json, null, "\t");
        } else {
            headRollupArea.innerText = "Failed to fetch head rollup."
        }
    }, 1000);
}

window.addEventListener(eventDomLoaded, initialize);