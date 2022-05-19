const eventDomLoaded = "DOMContentLoaded";
const idBlockHeadHeight = "blockHeadHeight";
const pathBlockHeadHeight = "/blockheadheight/";

const initialize = () => {
    const statusArea = document.getElementById(idBlockHeadHeight);

    setInterval(async () => {
        const blockHeadHeightResp = await fetch(pathBlockHeadHeight);

        if (blockHeadHeightResp.ok) {
            statusArea.innerText = await blockHeadHeightResp.text();
        } else {
            statusArea.innerText = "Failed to retrieve block head height."
        }
    }, 2000);
}

window.addEventListener(eventDomLoaded, initialize);