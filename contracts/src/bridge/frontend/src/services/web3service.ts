import { ethers } from 'ethers'
import { useMessageStore } from '@/stores/messageStore'
import ImageGuessGameJson from '@/assets/contract/artifacts/contracts/ImageGuessGame.sol/ImageGuessGame.json'
import ContractAddress from '@/assets/contract/address.json'
import { formatTimeAgo, trackEvent, bigNumberToNumber } from './utils'
import Common from '@/lib/common'
import Web3listener from './web3listener'
import { useGameStore } from '../stores/gameStore'
import { ElNotification } from 'element-plus'

export default class Web3Service {
  constructor(signer) {
    this.contract = new ethers.Contract(ContractAddress.address, ImageGuessGameJson.abi, signer)
    this.signer = signer
  }

  async getGuessHistory() {
    try {
      const challengeId = await this.getChallengeId()
      const historyTx = await this.contract.viewMyGuesses(challengeId)
      // above is the response from the contract, we need to convert it to a more readable format
      // the first array is the x and y coordinates of the guesses
      // the second array is the timestamp of the guesses
      // to get a readable timestamp, we need to fetch the block timestamp from the blockchain using ethers.js and subtract the block timestamp from the timestamp in the response

      const guessCoordinates = historyTx[0]
      const guessTimestamps = historyTx[1]
      const guessHistory = guessCoordinates.map((coordinates, index) => {
        const x = bigNumberToNumber(coordinates[0])
        const y = bigNumberToNumber(coordinates[1])
        const guessTimestamp = bigNumberToNumber(guessTimestamps[index])

        return {
          x,
          y,
          timestamp: formatTimeAgo(bigNumberToNumber(guessTimestamp)),
          win: 'N/A', // this is not available from the contract
          reward: 0 // this is not available from the contract
        }
      })
      return guessHistory
    } catch (error) {
      console.error('Failed to preload guess history - ', error)
    }
  }

  async submitGuess(challengeId, [coordinateX, coordinateY]) {
    // ensuring that coordinatex and coordinatey are not fractional components but whole numbers
    const updatedCoordinates = [Math.round(coordinateX), Math.round(coordinateY)]

    const entryCost = ethers.utils.parseEther(Common.ENTRY_COST)

    // ElNotification({
    //   message: 'Issuing Guess...',
    //   type: 'info'
    // })

    try {
      const submitTx = await this.contract.submitGuess(challengeId, updatedCoordinates, {
        value: entryCost,
        gasLimit: ethers.utils.hexlify(3000000)
      })
      const receipt = await submitTx.wait()
      // ElNotification({
      //   title: 'Issued Guess tx:'
      //   message: receipt.transactionHash,
      //   type: 'success'
      // })

      const web3listener = new Web3listener(this.signer)
      web3listener.startCheckingGuesses(receipt)

      const message = `Your guess has been submitted successfully! The winners will be announced in ${formatTimeAgo(bigNumberToNumber(receipt.events[0].args[4] || 0), false)}...`
      // ElNotification({
      //   message,
      //   type: 'success'
      //   duration: 0,
      // })

      return message
    } catch (e) {
      if (e.reason) {
        // ElNotification({
        //   message: 'Failed to issue Guess - ' + e.reason,
        //   type: 'error'
        // })
        console.error('Failed to issue Guess - ', e.reason)
        return
      }

      // ElNotification({
      //   message:
      //   'Failed to issue Guess - ' + e,
      //   type: 'error'
      // })
      console.error('Failed to issue Guess - ', e)
    }
  }

  async getChallengeId() {
    const messageStore = useMessageStore()
    try {
      const challengeId = await this.contract.currentChallengeIndex()
      const formattedChallengeId = bigNumberToNumber(challengeId)
      return formattedChallengeId
    } catch (error) {
      console.error('Failed to get challenge id - ', error)
    }
  }

  async createChallenge(payload) {
    const messageStore = useMessageStore()

    try {
      // create each challenge with each object in the array
      const createChallengeRes = await Promise.all(
        payload.map(async (challenge) => {
          // Estimate gas for the createChallenge transaction
          const estimatedGas = await this.contract.estimateGas.createChallenge(challenge)

          // Adding 10% buffer Gas
          const gasLimit = estimatedGas.add(estimatedGas.mul(10).div(100))

          const createChallengeTx = await this.contract.createChallenge(challenge, {
            gasLimit: gasLimit.toString()
          })

          const receipt = await createChallengeTx.wait()
          trackEvent('Challenge Created', {
            transactionHash: receipt.transactionHash,
            challengeId: await this.getChallengeId()
          })

          return receipt
        })
      )
      return createChallengeRes
    } catch (error) {
      console.error('Failed to create challenge - ', error)
      messageStore.addMessage('Failed to create challenge - ' + error.reason + ' ...')
    }
  }

  async getChallengePublicInfo() {
    try {
      const challengeId = await this.getChallengeId()
      const challenge = await this.contract.getChallengePublicInfo(challengeId)
      return challenge
    } catch (error) {
      console.error('Failed to get challenge properties - ', error)
    }
  }

  async addAdmin(address) {
    const messageStore = useMessageStore()
    try {
      const addAdminTx = await this.contract.addAdmin(address)
      const receipt = await addAdminTx.wait()
      messageStore.addMessage('Added admin - ' + receipt.transactionHash)
      return receipt
    } catch (error) {
      console.error('Failed to add admin - ', error)
      messageStore.addMessage('Failed to add admin - ' + error.reason + ' ...')
      throw error
    }
  }

  async getPreviousWins() {
    try {
      const gameStore = useGameStore()
      const challengeId = await this.getChallengeId()
      let previousChallenges = []
      for (let id = gameStore.isGameActive ? challengeId - 1 : challengeId; id >= 0; id--) {
        const historyTx = await this.contract.getRevealedChallengeDetails(id)

        previousChallenges.push({
          name: 'Challenge ' + id,
          topGuessesArray: historyTx[0],
          privateImageURL: historyTx[1]
        })
      }

      return previousChallenges
    } catch (error) {
      console.error('Failed to preload guess history - ', error)
    }
  }
}
