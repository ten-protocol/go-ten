const eventClick = "click";
const eventDomLoaded = "DOMContentLoaded";
const idGetBlockHeadHeight = "getBlockHeadHeight";
const idBlockHeadHeight = "blockHeadHeight";
const pathBlockHeadHeight = "/blockheadheight/";

const initialize = () => {
    const getBlockHeadHeightButton = document.getElementById(idGetBlockHeadHeight);
    const statusArea = document.getElementById(idBlockHeadHeight);

    getBlockHeadHeightButton.addEventListener(eventClick, async () => {
        const blockHeadHeightResp = await fetch(pathBlockHeadHeight);

        if (isOk(blockHeadHeightResp)) {
            statusArea.innerText = `Current block head height: ${await blockHeadHeightResp.text()}`;
        } else {
            statusArea.innerText = "Failed to retrieve block head height."
        }
    })
}

window.addEventListener(eventDomLoaded, initialize);

function isOk(response) {
    return response.status >= 200 && response.status < 300
}