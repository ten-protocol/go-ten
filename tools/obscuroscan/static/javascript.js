const eventDomLoaded = "DOMContentLoaded";
const idBlockHead = "blockHead";
const idHeadRollup = "headRollup";
const pathBlockHead = "/blockhead/";
const pathHeadRollup      = "/headrollup/";

const initialize = () => {
    const blockHeadArea = document.getElementById(idBlockHead);
    const headRollupArea = document.getElementById(idHeadRollup);

    setInterval(async () => {
        const headBlockResp = await fetch(pathBlockHead);

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