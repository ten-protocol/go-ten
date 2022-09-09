import secrets
import requests
import json
import time
import logging
import random
from web3 import Web3
from solcx import compile_source
from eth_account.messages import encode_defunct

WPORT = 3000
WHOST = '127.0.0.1'
LOWER = 0
UPPER = 100
FAUCET_URL = 'http://testnet-faucet.uksouth.azurecontainer.io/fund/obx'

guesser = '''
// SPDX-License-Identifier: MIT
pragma solidity 0.8;

contract Guesser {
    address public owner;
    uint256 public number;

    constructor(uint256 _initialNumber) {
        owner = msg.sender;
        number=_initialNumber;
    }

    function guess(uint256 i) view public returns (int) {
        if (i<number) return 1;
        if (i>number) return -1;
        return 0;
    }

    function destroy() public {
        require(msg.sender == owner, "You are not the owner");
        selfdestruct(payable(address(this)));
    }
}
'''

def guess(contract, max_guesses=100):
    lower = LOWER
    upper = UPPER
    nguess = 0
    while True:
        nguess += 1
        if nguess > max_guesses:
            logging.warn("Exceeded guess count ... exiting")
            return None

        guess = random.randrange(lower, upper)
        ret = contract.functions.guess(guess).call()
        if ret == 1:
            logging.info("Guess is %d, need to go higher" % guess)
            lower = guess + 1
        elif ret == -1:
            logging.info("Guess is %d, need to go lower" % guess)
            upper = guess
        else:
            logging.info("You've guessed the secret %s" % guess)
            return guess


def run():
    # connect to the network, create a local private key and convert into the account
    w3 = Web3(Web3.HTTPProvider('http://%s:%d' % (WHOST, WPORT)))
    private_key = secrets.token_hex(32)
    account = w3.eth.account.privateKeyToAccount(private_key)
    logging.info('Using account with address %s' % account.address)

    # request native OBX from the faucet server
    headers = {'Content-Type': 'application/json'}
    data = {"address": account.address}
    requests.post(FAUCET_URL, data=json.dumps(data), headers=headers)

    # generate a viewing key for this account, sign and post it to the wallet extension
    headers = {'Accept': 'application/json', 'Content-Type': 'application/json'}
    data = {"address": account.address}
    response = requests.post('http://%s:%d/generateviewingkey/' % (WHOST, WPORT), data=json.dumps(data), headers=headers)
    signed_msg = w3.eth.account.sign_message(encode_defunct(text='vk' + response.text), private_key=private_key)

    data = {"signature": signed_msg.signature.hex(), "address": account.address}
    response = requests.post('http://%s:%d/submitviewingkey/' % (WHOST, WPORT), data=json.dumps(data), headers=headers)

    # compile the guessing game and build the deployment transaction
    logging.info('Compiling the guessing game application')
    compiled_sol = compile_source(guesser, output_values=['abi', 'bin'], solc_binary='/opt/homebrew/bin/solc')
    contract_id, contract_interface = compiled_sol.popitem()
    bytecode = contract_interface['bin']
    abi = contract_interface['abi']
    contract = w3.eth.contract(abi=abi, bytecode=bytecode)
    build_tx = contract.constructor(random.randrange(LOWER, UPPER)).buildTransaction(
        {
            'from': account.address,
            'nonce': w3.eth.getTransactionCount(account.address),
            'gasPrice': 1499934385,
            'gas': 720000,
            'chainId': 777
        }
    )

    # Sign the transaction and send to the network
    logging.info('Signing and sending raw transaction')
    signed_tx = account.signTransaction(build_tx)
    tx_hash = None
    try:
        tx_hash = w3.eth.sendRawTransaction(signed_tx.rawTransaction)
    except Exception as e:
        logging.error('Error sending raw transaction %s' % e)
        return

    # wait for the transaction receipt and check the status
    logging.info('Waiting for transaction receipt')
    start = time.time()
    tx_receipt = None
    while True:
        if (time.time() - start) > 60:
            logging.error('Timed out waiting for transaction receipt ... aborting')
            return

        try:
            tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
            if tx_receipt.status == 0:
                logging.error('Transaction receipt has failed status ... aborting')
                return
            else:
                logging.info('Received transaction receipt')
                break
        except Exception as e:
            time.sleep(1)

    # construct the contract using the contract address
    logging.info('Contract address is %s' % tx_receipt.contractAddress)
    contract = w3.eth.contract(address=tx_receipt.contractAddress, abi=abi)

    # guess the number
    logging.info('Starting guessing game')
    guess(contract)


if __name__ == '__main__':
    logging.getLogRecordFactory()
    logging.basicConfig(level=logging.INFO, format='%(asctime)s %(levelname)s %(message)s')
    run()

