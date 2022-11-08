# Testnet KPIs

Key performance indicators (KPIs) will be used to determine the amount of Testnet traction over time.

## Selected KPIs
| KPI NAME | RATIONALE | SOURCE | METRIC | TARGET |
|--|--|--|--|--|
| Testnet documentation page views (https://docs.obscu.ro/testnet/) | Gauges interest from developers in building on testnet. A good early indicator. Straightforward to capture | Google Analytics | Number of unique page views in the last 4 weeks | Targets will be determined once initial baseline data is captured |
| Testnet uptime | Captures how robust and ready for mainnet Obscuro is | DataDog avg:system.uptime{*} | Average Testnet uptime over the last 4 weeks||
| Number of end user accounts| Captures users using testnet in general and testing new dApps and shows the pipeline from wallet extension download to actively using testnet | TBD | Number of new end user account addresses in the last 4 weeks ||
|DApps in development| This is difficult to measure as it’s not something we can look to the network to accurately measure. We can use our business development CRM data as an indicator. | Airtable "Partnerships" tab| Number of partners with status changed to "Soft Commitment" in the last 4 weeks||
| New dApps deployed on Testnet| Indicates how much new activity and commitment of effort from dApp builders| ObscuroScan?| Number of new dApp addresses deployed in the last 4 weeks||
| Wallet extension downloads| Good proxy for the number of active users. Straightforward to capture.| GitHub|Number of wallet extension downloads in the last 4 weeks||

## Discounted KPIs
These KPIs were considered and discounted.

- Ported Solidity dApps– original vs fork
    - One of Obscuro’s key promises is that existing Solidity dApps will just work on Obscuro. This KPI helps track whether Obscuro is fulfilling this promise. The reason for differentiating between original fork e.g. Uniswap vs PancakeSwap is it further provides insight into whether original teams believe in Obscuro.
    - Discounted because very difficult to determine of the dApps on Testnet which have been ported.
- Number of transactions
    - Gauge the amount of activity on Testnet.
    - Discounted because it can be gamed too easily leading to misleading results.
- Number of nodes
    - Nodes secure the network and track general interest in Obscuro. This KPI also provides an input into the staking function and has an impact on tokenomics.
    - Discounted for now because too early in the lifecycle of Testnet to bring nodes online.
- Faucet requests
    - Tracks the number of tokens being requested and whether users are coming back for more. Straightforward to capture.
    - Discounted because it can be gamed too easily leading to misleading results.
- Number of organisations building on Testnet
    - Funded organisations building on Obscuro vs individuals. Good indicator of the long-term health of dApps being built on Obscuro and expands on number of new dApps.
    - Discounted because large effort required to determine who is building on Testnet.

