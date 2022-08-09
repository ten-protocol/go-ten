"use strict";

const eventDomLoaded = "DOMContentLoaded";

const idAttestation = "attestation";
const idAttestationReport = "attestationReport";

const pathAttestation = "/api/attestation/";
const pathAttestationReport = "/api/attestationreport/"

// Updates the displayed stats.
async function displayAttestation() {
    const fieldAttestation = document.getElementById(idAttestation);
    const fieldAttestationReport = document.getElementById(idAttestationReport);

    const respAttestation = await fetch(pathAttestation);
    if (respAttestation.ok) {
        const attestationJSON = JSON.parse(await respAttestation.text());
        fieldAttestation.innerText = JSON.stringify(attestationJSON, null, "\t");
    } else {
        fieldAttestation.innerText = "Failed to fetch attestation.";
    }

    const respAttestationReport = await fetch(pathAttestationReport);
    if (respAttestationReport.ok) {
        const attestationReportJSON = JSON.parse(await respAttestationReport.text());
        fieldAttestationReport.innerText = JSON.stringify(attestationReportJSON, null, "\t");
    } else {
        fieldAttestationReport.innerText = "Failed to fetch attestation.";
    }
}

const initialize = async () => {
    await displayAttestation();
}

window.addEventListener(eventDomLoaded, initialize);