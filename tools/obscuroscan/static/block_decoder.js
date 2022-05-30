'use strict';

const eventDomLoaded = "DOMContentLoaded";
const formDecryptRollup = "form-decrypt-rollup"
const typeSubmit = "submit"

const initialize = () => {
    document.getElementById(formDecryptRollup).addEventListener(typeSubmit, function(e) {
        e.preventDefault(); // before the code
        console.log('intercepted');
    })
}

window.addEventListener(eventDomLoaded, initialize);