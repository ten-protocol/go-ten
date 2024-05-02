/** @type import('hardhat/config').HardhatUserConfig */

require('dotenv').config()
require('@nomicfoundation/hardhat-toolbox')
require('@nomiclabs/hardhat-ethers')

const { PRIVATE_KEY, USER_KEY } = process.env

task('deploy', 'Deploys the ImageGuessGame contract').setAction(async (taskArgs, hre) => {
  const ImageGuessGame = await hre.ethers.getContractFactory('ImageGuessGame')
  const imageGuessGame = await ImageGuessGame.deploy()
  await ImageGuessGame.deployed()

  console.log('ImageGuessGame deployed to:', imageGuessGame.address)
})

module.exports = {
  solidity: '0.8.19',
  paths: {
    sources: './contracts',
    artifacts: './src/assets/contract/artifacts'
  },
  networks: {
    ten: {
      chainId: 443,
      url: `https://testnet.ten.xyz/v1/${USER_KEY}`,
      gasPrice: 2000000000,
      accounts: [`0x${PRIVATE_KEY}`]
    }
  }
}
