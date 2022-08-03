# Obscuroscan

See the documentation [here](https://docs.obscu.ro/wallet-extension/wallet-extension.html).

## Developer notes

### Updating the guessing game

* Clone the source: `https://github.com/obscuronet/number-guessing-game`
* Compile the source: `npm install && npx hardhat compile && npm run build`
* Rename `public/index.html` to `public/game.html`
* Copy the contents of the `public` folder into `tools/obscuroscan/static`
