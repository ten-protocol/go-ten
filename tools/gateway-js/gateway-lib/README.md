# TEN Gateway JS

A JavaScript library for the gateway, providing streamlined access and functionalities for interacting with the TEN network.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Build](#build)
- [Contribute](#contribute)
- [License](#license)

## Features

- Seamless connection to TEN Network.
- Easy-to-use methods for joining and authenticating.
- External consumption support through CDN or NPM.

## Installation

To install `ten-gateway-js`, use npm:

\`\`\`bash
npm install ten-gateway-js
\`\`\`

## Usage

\`\`\`javascript

const Gateway = require('ten-gateway-js');

const gateway = new Gateway(httpURL, wsURL, provider);
await gateway.join();
await gateway.registerAccount(account);

\`\`\`

## Build

To build for development:

\`\`\`bash
npm run dev
\`\`\`

For production:

\`\`\`bash
npm run build
\`\`\`

The production build will be available for external consumption on GitHub Pages at `https://go-ten.github.io/tennet/gateway.bundle.js`.

## Contribute

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
