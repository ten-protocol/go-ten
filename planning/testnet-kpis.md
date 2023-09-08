# Measuring Testnet Success
The launch of a testnet in Web3 is a significant milestone for all projects. It demonstrates the capability of the solution, whether the project is likely to meet the promises shared with investors and with the community and it signifies a confident step towards mainnet. The final iterations of a testnet should very closely emulate mainnet and give assurances to users, developers and investors that the final product will be of a high quality with a high chance of success. As a result, the impression left by testnet is crucial to the expected success of a project and how it is perceived. The primary contributors of success for Obscuro testnet are whether it is attractive, it is being used and whether or not users have a positive experience with it. Making the degree of success quantifiable means defining mesaurable success criteria and collecting data to know whether those criteria are being met.

## Testnet Success Criteria
Determining success for testnet will be a data-driven exercise, this being the best way to make measurable and repeatable observations. These observations can subsquently feed into decision-making with outcomes again being measured and compared. Included in the measurements will be criteria which, on the face of it, do not provide value however they have gained traction in the Web3 communicty as success indicators by which projects are judged. We need commentators to be able to compare Obscuro to other projects using like-for-like data points which the Web3 community are comfortable with even if they offer little value, or can even be misleading. For example, total number of transactions in a given period of time is a data point commonly used to compare projects yet it can be easily gamed.  
The testnet success critera have been expressed below in the form of Key Performance Indicators with the rationale for their inclusion, the data source, the actual metric and the target value.

## Testnet KPIs
Key performance indicators (KPIs) will be used to determine the amount of Testnet traction over time.

### Selected KPIs for Users
| KPI NAME | RATIONALE | SOURCE | METRIC | TARGET |
|--|--|--|--|--|
| Testnet uptime | Captures how robust and ready for mainnet Obscuro is | DataDog avg:system.uptime{*} | Average Testnet uptime over the last 4 weeks|99.9%|
| Number of wallets connected to Obscuro Gateway| Good proxy for the number of active users. Straightforward to capture.| Datadog? |Number of daily connections|500|
| Number of transactions| Typical guage of the amount of activity on testnet (even though it can be gamed).| Datadog? |Number of daily transactions|2000|
| Number of RPC requests| Alternative guage of the amount of activity on testnet. Can also show where RPC performance degrades.| Datadog? |Number of daily RPC requests|2000|

### Selected KPIs for Developers
| KPI NAME | RATIONALE | SOURCE | METRIC | TARGET |
|--|--|--|--|--|
| Testnet documentation page views (https://docs.obscu.ro/testnet/) | Gauges interest from developers in building on testnet. A good early indicator. Straightforward to capture | Google Analytics | Number of unique page views in the last 4 weeks |50|
|DApps in development| This is difficult to measure as it’s not something we can look to the network to accurately measure. We can use our business development CRM data as an indicator. | Airtable "Partnerships" tab| Number of partners with status changed to "Soft Commitment" in the last 4 weeks|40|
| New dApps deployed on Testnet| Indicates how much new activity and commitment of effort from dApp builders| ObscuroScan?| Number of new dApp addresses deployed in the last 4 weeks|20|

## Discounted KPIs
These KPIs were considered and discounted for now.

- Wallet extension downloads
    - Good proxy for the number of active users.
    - Straightforward to capture.
    - Discounted because Obscuro Gateway is the recommended wallet connection method.
- Ported Solidity dApps– original vs fork
    - One of Obscuro’s key promises is that existing Solidity dApps will just work on Obscuro. This KPI helps track whether Obscuro is fulfilling this promise. The reason for differentiating between original fork e.g. Uniswap vs PancakeSwap is it further provides insight into whether original teams believe in Obscuro.
    - Discounted because very difficult to determine of the dApps on Testnet which have been ported.
- Number of nodes
    - Nodes secure the network and track general interest in Obscuro. This KPI also provides an input into the staking function and has an impact on tokenomics.
    - Discounted for now because Obscuro Labs will run the nodes. When validator nodes can be run by others, this KPI will be introduced.
- Faucet requests
    - Tracks the number of tokens being requested and whether users are coming back for more. Straightforward to capture.
    - Discounted because it can be gamed too easily leading to misleading results.
- Number of organisations building on Testnet
    - Funded organisations building on Obscuro vs individuals. Good indicator of the long-term health of dApps being built on Obscuro and expands on number of new dApps.
    - Discounted because large effort required to determine who is building on Testnet.

