This package provides a standard way to simulate user activity against a test network.

## Runner
The Runner interface is the core of it, with a Start(), Stop(), Verify() pattern it tries to be flexible for different simulation scenarios.

This allows you to simulate users sending transactions of different types, spamming out viewing key registrations and event 
subscriptions or simulating activity in an on-chain game.

Some tests might even run a combination of Runners to simulate a more diverse and realistic range of network activity.

## Test
The DurationTest pairs a Runner with a Duration and verifier(s) to make it an executable test. This separation is useful
because sometimes we want to use a Runner in complex scenarios or just running indefinitely and don't always want to 
provide a fixed duration.