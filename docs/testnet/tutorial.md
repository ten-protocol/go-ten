# Getting Started
In this tutorial, you build your own Obscuro dApp from the start. This is a good way to experience a typical development process as you learn how Obscuro dApps are built, design concepts, tools and terminology.

# Lets go
This 'Guessing Game' tutorial provides an introduction to the fundamentals of Obscuro and shows you how to:
- Set up your local Ethereum development environment for Obscuro.
- Use Obscuro concepts to develop a dApp

What is the *Guessing Game*?

The Guessing Game is a dApp that you will build that demonstrates a basic Obscuro use case, which is a simple number guessing game. The contract generates a random secret number when it's deployed, which is never revealed to an operator or end-user because of the privacy benefits of Obscuro. The goal of the game is to guess this number, and each time an attempt is made, an entrance fee of 1 token is paid. If a user correctly guesses the number, the contract will pay out all of the accumulated entrance fees to them, and reset itself with a new random number.

Without Obscuro, it would be possible to look up the internal state of the contract and cheat, and the game wouldn't work.

The dApp has many of the features that you'd expect to find, including:
- Generating a random number known only to the dApp
- Allow users to guess the numbers by depositing a token for each play
- Colder/Warmer functionality to help guide users to the correct number
- Events to alert users whether they have won or not
- Track all users plays
- Contract owner privileges to start and stop the game

# Set up your environment
- You'll need to install MetaMask following the instructions [here](https://metamask.io/)
- Then download and run the Obscuro Wallet extension following the instructions [here](https://docs.obscu.ro/wallet-extension/wallet-extension.html)
- Now configure MetaMask following the instructions [here](https://docs.obscu.ro/wallet-extension/configure-metamask.html)
- Finally, check you can open the Remix IDE by visiting the following URL in your browser https://remix.ethereum.org

That's it, you're all set to start building your first dApp on Obscuro.

## 1. To begin, we'll clone a a basic repo from Github through Remix by following these instructions

### a. Click on GitHub from the main screen in Remix
![Step 1 screenshot](https://images.tango.us/workflows/454b9fc4-c6f9-43d1-beef-9cbbff7f0b0b/steps/7cb0a9d9-2e3d-4c33-b10c-4f89dea34206/a3c964d3-6d40-4931-9ddc-8390f5744452.png?crop=focalpoint&fit=crop&fp-x=0.2873&fp-y=0.4869&fp-z=2.9198&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=1812%3A1226)


### b. And enter "https://github.com/obscuronet/tutorial/blob/main/number-guessing-game/contracts/Guess.sol" into form and hit import
![Step 2 screenshot](https://images.tango.us/workflows/454b9fc4-c6f9-43d1-beef-9cbbff7f0b0b/steps/48d5b32f-eda8-4e7d-aefc-889d64147192/70a4e414-27b0-4238-9529-f1d64c7fbffd.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.2292&fp-z=1.8200&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=1812%3A1226)


### c. Select Guess.sol and we're ready
![Step 3 screenshot](https://images.tango.us/workflows/454b9fc4-c6f9-43d1-beef-9cbbff7f0b0b/steps/e35c3f04-92f2-4edf-84ce-ea467feb8797/d401a509-0ab6-4ce8-892d-fbce81c0aeff.png?crop=focalpoint&fit=crop&fp-x=0.1076&fp-y=0.2406&fp-z=2.6507&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=1812%3A1226)


# [Workflow with Ethereum and Github](https://app.tango.us/app/workflow/41b454a9-b560-4d60-8775-3e40c5bc1466?utm_source=markdown&utm_medium=markdown&utm_campaign=workflow%20export%20links)


![Step 1 screenshot](https://images.tango.us/workflows/41b454a9-b560-4d60-8775-3e40c5bc1466/steps/74347d9b-e3fe-48aa-8e74-fecebd22fdfa/80231e65-e0be-4bf2-8150-521721745b2c.png?crop=focalpoint&fit=crop&fp-x=0.3145&fp-y=0.5909&fp-z=3.3270&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=2307%3A1244)



## 2. Exploring Guess.sol
Guess.sol is a Solidity file, and once we're done, it will contain everything we need for the 'Guessing Game'. Inside the file you'll find the classic ERC20 interface and a simple implementation named ERC20Basic representing OGG (Obscuro Guessing Game token) which we'll use as means of players entering to play and the prize pool.

There is also an empty contract called Guess inside Guess.sol, and this is what we'll be extending. The contract Guess will be written the same as any other Solidity contract in Ethereum and it'll run on the EVM so will exhibit the same behhaviours.

Let's start by by defining the key variables we'll need for the game. One to store the secret number that only the contract knows about and another containing the ERC-20 contract.

```
    uint8 private _target
    IERC20 public erc20
```
Upon seeing this, you might be thinking, hang on, I can _target is set to private, but so what? Private variables in Ethereum don't really mean they're private, just that they can't be accessed ot modified by other smart contracts. And you're right, that's the case in Ethereum. But in Obscuro, private variables really are private (unless you explicitly decide to reveal them!). This is because all contracts and contract data is encrypted and only accessible within enclaves.

Our variables should now look like this
![Step 1 screenshot](https://images.tango.us/workflows/dc61d575-7eea-457c-bdfe-8edccf79b366/steps/3d9b387b-a492-4096-aa8a-1d50084be0d5/6d08adce-70d2-4203-8d18-c09999536553.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.5000&fp-z=1.0000&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=413%3A131)


Next let's think about how we'll instantitate this contract. We'll need to:
- create an instance of the ERC-20 contract to manage and hold OGG tokens
- generate a secret random number


Let's extend to the constructor include these:

![Step 2 screenshot](https://images.tango.us/workflows/a8278591-dc09-4e26-97a2-b30776f86179/steps/58146527-d6d4-47c4-a796-cc0b3b4909db/56afa096-9f71-49a0-a85d-94ba3ecc56a4.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.5000&fp-z=1.0000&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=410%3A117)


Now let's write our function to generate the random number and secretly store it within the contract

![Step 1 screenshot](https://images.tango.us/workflows/a8278591-dc09-4e26-97a2-b30776f86179/steps/8b18e51b-b4da-43e6-a752-eea0cc550791/009e7455-a298-4b66-808b-39080f30af24.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.5000&fp-z=1.0000&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=864%3A211)


This isn't great as although the number is arguably random, it's not random to anyone that both knows how it's generated and also knows how to measure the inputs at the time the function is called. We'll talk at the end how we can do better!

Okay, so right now we should have:
- a constructor that initializes the contract
- a way to calculate the random number to be guessed

We now need to add the ability for a player to make an actual guess. But before we do that, we have not thought about the range we might want numbers to appear in. This would allow us to set the difficulty for any particular game. Let's do that now.

Extend the constructor accept a range int and create a variable for it so we have something that looks like this

![Step 1 screenshot](https://images.tango.us/workflows/8f7ef3ef-091f-4b5d-9017-d041449139ae/steps/43aff627-6025-42c1-8e1a-654761d01b90/c3fc2de4-e9fa-4d31-b1a8-abe5331e0ec9.png?crop=focalpoint&fit=crop&fp-x=0.2957&fp-y=0.3023&fp-z=4.0000&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=2307%3A1244)

And we can change our random number function to adhere to the range by including the modulo operator %, like this:

![Step 1 screenshot](https://images.tango.us/workflows/3365689e-af3c-4c4c-9a8a-890f9e1fa4c5/steps/34d0f0b0-799f-41cb-b9ff-ff475af306cd/0ed09dad-0b61-4fd3-9bf8-626a897ac050.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.5000&fp-z=1.0000&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=850%3A198)


Onto the actual guess! We'll need the function to capture several things:
- Check the player has enough tokens to play
- Allow the player to transfer a token to the prize pool to pay for a guess
- Handle the player winning
- Handle the player not winning

Let's go ahead and create a function called *attempt* that takes as input a guess

![Step 1 screenshot](https://images.tango.us/workflows/560f6b90-fcc4-4c59-8a27-bf8a5fc08370/steps/bf2c1e4e-8323-4b89-85fd-e35805aff122/e6271350-4eb1-41d5-9b3d-afaf3058040a.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.5000&fp-z=1.0000&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=473%3A100)

We already have the OGG tokens defined so we can use this to check funds with 

```
    require(erc20.allowance(msg.sender, address(this)) >= 1 ether, "Check the token allowance.");
```
This assumes it costs 1 token per play.

Next we transfer the funds for each play with something like

```
erc20.transferFrom(msg.sender, address(this), 1 ether);
```

Next add code that checks if the guess was correct or not. If it's correct, call the ERC-20 function to transfer all the funds in the contract and reset the game by calling _setNewTarget().

If it's not correct, we do nothing. Your code might look something like this:

![Step 1 screenshot](https://images.tango.us/workflows/7073e589-7d79-491f-b4f1-27db2eb144c1/steps/ad4b9d9a-5793-420f-b21c-3a22ae8e83b5/b016b98e-203a-4362-9c14-6ff2928b9ba8.png?crop=focalpoint&fit=crop&fp-x=0.5000&fp-y=0.5000&fp-z=1.0000&w=1200&mark-w=0.2&mark-pad=0&mark64=aHR0cHM6Ly9pbWFnZXMudGFuZ28udXMvc3RhdGljL21hZGUtd2l0aC10YW5nby13YXRlcm1hcmsucG5n&ar=879%3A246)


This works. You now have an Obscuro contract that will let players play the guessing game. But, it's not very good, right? How do players know if they've won or not? We need events!

In Ethereum events are typically written to the event log that is viewable by everyone. In Obscuro, events are viewable only by who they're intended for. This is done by encrypting the events a key that only the intended recipient knows and is handled for the recipient through the Wallet Extension.





Reset previous tries

Add a way to show how developers can make things public

How might we improve this so that the number is random to even the developers


