Util package for retrying, functions and tools for situations where we are retrying/polling for a change.

Based around a 'RetryStrategy' interface that allows different approaches such as:
- retry X times
- retry for X seconds
- retry for 2 minutes but backing off

