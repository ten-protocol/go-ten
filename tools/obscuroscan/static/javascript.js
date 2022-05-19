const eventDomLoaded = "DOMContentLoaded";
const idBlockHeadHeight = "blockHeadHeight";
const idHeadRollup = "headRollup";
const pathBlockHeadHeight = "/blockheadheight/";
const pathHeadRollup      = "/headrollup/";

const initialize = () => {
    const blockHeadHeightArea = document.getElementById(idBlockHeadHeight);
    const headRollupArea = document.getElementById(idHeadRollup);

    setInterval(async () => {
        const blockHeadHeightResp = await fetch(pathBlockHeadHeight);

        if (blockHeadHeightResp.ok) {
            blockHeadHeightArea.innerText = await blockHeadHeightResp.text();
        } else {
            blockHeadHeightArea.innerText = "Failed to retrieve block head height."
        }
    }, 1000);

    setInterval(async () => {
        const headRollupResp = await fetch(pathHeadRollup);

        if (headRollupResp.ok) {
            headRollupArea.innerText = await headRollupResp.text();
        } else {
            headRollupArea.innerText = "Failed to retrieve head rollup."
        }
    }, 1000);
}

window.addEventListener(eventDomLoaded, initialize);