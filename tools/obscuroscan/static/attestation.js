"use strict";

const eventDomLoaded = "DOMContentLoaded";

const idAttestation = "attestation";

const pathAttestation = "/attestation/";

// Updates the displayed stats.
async function displayAttestation() {
    const fieldAttestation = document.getElementById(idAttestation);

    const respAttestation = await fetch(pathAttestation);
    if (respAttestation.ok) {
        const attestationJSON = JSON.parse(await respAttestation.text());
        fieldAttestation.innerText = JSON.stringify(attestationJSON, null, "\t");
    } else {
        fieldAttestation.innerText = "Failed to fetch attestation.";
    }
}

const initialize = async () => {
    await displayAttestation();
}

window.addEventListener(eventDomLoaded, initialize);