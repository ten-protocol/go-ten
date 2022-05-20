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
            blockHeadArea.innerText = await headBlockResp.text();
        } else {
            blockHeadArea.innerText = "Failed to retrieve head block."
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