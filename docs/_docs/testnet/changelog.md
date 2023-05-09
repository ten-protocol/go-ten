---
---
# Obscuro Testnet Change Log

# May 2023-05-04 (v0.13)
* A list of the PRs merged in this release is as below;
    * `6f8264f7` Fix the we cross-os build (#1244)
    * `325bb0c8` Adding system errors (#1237)
    * `f346e239` Add wait group wait() timeout function to util (#1242)
    * `b759f55a` Remove error from proto message payload (#1241)
    * `fdc9b035` Fix broken link for the pdf version of the whitepaper. (#1238)
    * `61b9e17f` Update some of the @matt todos (#1240)
    * `13c87ac1` Updates after backlog review part three (#1239)
    * `a551408b` Layout update for obscuroscan (#1200)
    * `eeb1623b` Locking the docker image alpine to 16 (#1233)
    * `6758cb32` Update wallet extension docs and include a link to obx faucet (#1235)
    * `a1a30083` Logs are now debuggable from the debug_log visibility endpoint (#1231)
    * `3497be98` Only perform the l2 deployment on a scheduled basis (#1230)
    * `3bda160e` Reintroducing encrypted enclave errors (#1197)
    * `56557a2f` Host to retry forever waiting for the enclave to come up (#1229)

# April 2023-04-13 (v0.12)
* Robustness fixes for event log subscriptions, including ensuring no dropping of event logs, removal of any duplicates, 
  and ensuring correct ordering based on their block insertion. 
* Event log subscriptions will now infer the correct viewing key to use for relevancy checks performed by an Obscuro
  node. This is based on filter criteria in the subscription, e.g. when a particular address field is included 
  in the filter list, if a viewing key for that address is registered, that key will be used. Note that subscriptions to 
  all events, where those events might contain data used in the relevancy checks, will not be able to infer the viewing 
  key to use, and as such will use the first registered key made through that instance of the wallet extension. 
* A list of the PRs merged in this release is as below;
    * `25dc4cd9` Block provider now contributes to health status (#1225)
    * `978e3509` Renaming docker volumes (#1224)
    * `827d67bc` Workaround for xchain message finality (#1220)
    * `56589db7` Persist enclave key (#1223)
    * `845d7d34` Is debug namespace enables flag (#1222)
    * `a3b220d1` Disabled notifications for npm (#1221)
    * `96e4029b` Adding debug_traceTransaction endpoint and functionality (#1214)
    * `37835960` Add different manual tests (#1217)
    * `c1425fbf` Reconnect client on the block provider (#1216)
    * `b0983690` Fix npe (#1215)
    * `e5efecc6` Fix sqlite init logic (#1213)
    * `05a7065f` Normalise data field to fix event dups bug (#1212)
    * `0fc9d164` Fix bug and improve log (#1211)
    * `bb92ae19` Add more logging in the verbose we (#1210)
    * `02d885ea` Review of todos part two (#1209)
    * `c1033393` Setup ego signer for testnets (#1208)
    * `7c2e8985` Pull the updated docker images in testnet upgrade script (#1207)
    * `13d8f102` Fix the block hash of log messages and receipts (#1206)
    * `4766f27b` Pedro/log tx rollup hash (#1204)
    * `6e81b9e2` Remove -a option (#1203)
    * `81aacd56` Eth2network now pushed as part of the deploy flow (#1201)
    * `e8903b38` Add logic to discover the authenticated client to be used for subscriptions (#1196)
    * `cad27131` Fix filtering logic (#1199)
    * `90b04263` Fix local network launcher for enclave persistence (#1198)
    * `48acdc3f` Use docker mount for enclave persistence (#1195)
    * `31562358` Rename offchain to obscall and fix tx execution logs (#1194)
    * `f3f73da6` Changes to the testnet launcher, better logging (#1192)
    * `124eb4cf` Fix names of types aliases (#1193)
    * `c6316cbb` Store events in sql table (#1182)
    * `17a7c44b` Todo review number 1 edits (#1191)
    * `927b970d` L1_host not output properly in upgrade gh action (#1190)
    * `013e8f32` Wipe and re-clone obscuro code on testnet node upgrade (#1186)
    * `f7d4cd35` Revert "refactoring how errors are done. (#1153)" (#1189)
    * `8dd0f1f6` Resolve faucet docs images (#1188)
    * `24659498` Add back in the faucet steps using discord (#1187)
    * `00fedf0f` Run scheduled dev-testnet deployment at 3:05am tues to sat (#1185)
    * `9af0f70c` Add in the release notes for v0.11 (#1177)
    * `45984ad1` Replaced 1.17->1.18 (#1184)
    * `c010d003` Refactoring how errors are done (#1153)
    * `6b96286d` Add docker_buildkit=1 (#1183)
    * `c3b76592` Remove unused host flow controls (#1181)
    * `3afc850e` Fix: node launcher was ignoring the l1 start param (#1180)
    * `18420198` Simplify sqlite creation logic (#1179)
    * `81483e76` Remove db interfaces (#1178)
    * `e7343851` Close database and guard enclave methods (#1175)
    * `9c9337c4` Removing the external api calls section (#1176)
    * `8106a432` Expose underlying sql database (#1171)
    * `235de492` Extract the mgmt contract deployment block for node starts (#1174)

# March 2023-03-22 (v0.11)
* A list of the PRs merged in this release is as below;
    * `1d22729e` Fixed nil ptr case. (#1172)
    * `96565a7d` Bridge and cross chain messaging docs. (#1087)
    * `dfc2b284` Added some logs for event relevancy. (#1169)
    * `915794a4` Process network secrets now iterates only over successful transactions. (#1168)
    * `0d8e947e` Tell prysm node not to look for peers (#1167)
    * `72cf6948` Retry l1 reconnect indefinitely (#1166)
    * `2678359a` Fix dead links. (#1156)e$ 
    * `266fd18c` Fix a dead link (#1157)e$ 
    * `f32bd928` Fixes for 3 dead links. (#1158)
    * `c3147be0` Prysm upgrade + fix compatibility (#1165)
    * `2e38f551` Set sequencer host to use in-mem db as well (#1164)
    * `496e626e` Move network config persistence further down stack (#1161)
    * `290c6c53` Docker testnet launcher to use in-mem db for host (#1162)
    * `6527a019` Update fast_finality.md (#1163)
    * `a483fe3f` Add upgrade option to node script (#1154)
    * `d41ad9ea` Small logging improvements, loosen islatest restriction (#1151)
    * `59762afe` Enable persistence on testnet hosts (#1152)
    * `a6087483` Fix loglvl arg format in testnet script (#1150)
    * `5bd5c21e` Set host and enclave loggers to info by default (#1149)
    * `8b85efe5` Reenable disabled block provider test (#1148)
    * `2e335e8e` Enclave does not process failed rollups (#1138)
    * `fddfeccb` Wallet extension verbose logging (#1147)
    * `824d74f8` Deploy l2 contracts synchronously (#1146)
    * `31b1d017` Add buildkit to docker update step (#1145)
    * `ee335754` Add basic l1 sanity check to networktests tools (#1142)
    * `af7f7259` Updated readme. (#1140)
    * `70197a13` Update the go versions on actions to match the repo version (#1143)
    * `dbda7aec` Fix to the admin address (#1141)
    * `d0a2d02d` Rollup module (#1125)
    * `d4637cc1` Trigger a repository dispatch to the faucet (#1139)
    * `6eb94f9b` Renamed increaseallowance to approve. (#1136)
    * `486892e5` Admin role for bridge. (#1135)
    * `f4e8f2ab` Clean up design documents (#1137)
    * `ec41be07` Draft upgrading design (#1103)
    * `b5a79eb6` Add test that restarts the host as well as the enclave (#1134)
    * `b3b15faa` Fix to the format of the dispatch event (#1133)

# February 2023-02-23 (v0.10)
* A list of the PRs merged in this release is as below;
    * `d81f5f9a` Run a schedule deploy on the l1, and trigger l2 if succesful (#1129)
    * `481dc317` Wrap leveldb so its error types don't leak into our codebase (#1128)
    * `e5d8c398` Resilient to rpc requests while sequencer unknown (#1130)
    * `aa3eaea2` Updated go version and ego version. (#1124)
    * `6de02fb0` Add config to enable host persistence (#1122)
    * `7f7b78fa` Save deployment container logs as build artifacts for fault finding (#1123)
    * `7b640aec` No need to request secret for enclave on host restart (#1117)
    * `71d0ea88` Testnet launcher + node package (#1112)
    * `c33bf2bc` Clarify p2p address syncing (#1116)
    * `1c6f7c7e` Add integration testing and local network tools (#1115)
    * `a4d34075` Implemented scheduled deployment of dev-testnet (#1114)
    * `8adb8970` Rotates the eth2 logs (#1113)
    * `a9c8f9af` Update testnet-local-build_images.sh (#1110)
    * `77cc0dad` Management contract in all vars (#1093)
    * `7bffad4d` Add a single start script to launch testnet with default arguments (#1107)
    * `ff7632fa` Adds a script that waits for the eth2network to be ready (#1108)
    * `76a0abbb` Removes the gethnetwork (#1095)
    * `70767522` Add blockheight and blockhash log keys (#1105)
    * `d4bb4120` Fix comparison for empty hash (#1106)
    * `6e145dce` Update enclave.go (#1104)
    * `55f4fdc9` Fixes for the we tests (#1102)
    * `de1d8bf4` Record latest rollup for every l1 block to avoid scanning back (#1101)
    * `859524b2` Attempt to fix github sleep step bash syntax (#1099)
    * `9196f32b` Aggressive enclave logging to track what what it is processing (#1100)
    * `f56418b9` Fix contract deployer test ports (#1097)
    * `a10e94a7` Add sleep to stagger az cli requests on parallel node creation (#1098)
    * `855c3a36` Fix wallet extension flakyness (#1096)
    * `04af95bd` Fix for contract ci go test integration. reintroduced param names to … (#1094)
    * `2e0dfdc1` Eth2 used in the deployment (#1079)
    * `3ef6a073` Remove speculative execution (#1092)
    * `ba736919` Adding manual tests + minor fixes (#1090)
    * `1d3134b7` Fix sequencer tx production (#1091)
    * `73dc2df8` Rollup lookup should not assume that l1 blocks go back to genesis (#1089)
    * `09932b44` After enclave restart, replay batches to restore missing statedb cache (#1088)
    * `d0bb547f` Provide a clearer way for finding the latest rollup in the chain of a given head block. (#1085)
    * `c610271f` Fixes enclave debug image (#1086)
    * `e2c2babc` Avoid swallowing l1 client error messages (#1081)
    * `54819a09` Ensure l2 head is rewound after l1 forks (#1078)
    * `afafeee2` Call db client batches dbbatches to avoid confusion (#1070)
    * `a8357aa3` Obscuro bridge deployment scripts. (#1046)
    * `c951a9b8` Reduces the time to start a merged eth2 network (#1077)
    * `7a65331c` Make pr template lighter (#1076)
    * `015c6b7e` Juiced the docker. it's now fast. (#1071)
    * `9e7c6a30` Upgrade enclave l2 (#1075)
    * `5e494c36` Fixes the enclave start (#1074)
    * `6f2b8a69` Remove early exit (#1073)
    * `14d25d96` Fix faucet deployment trigger (#1072)
    * `b2e6f8e2` Store the l2 head hash directly in the enclave (#1069)
    * `34e7bbe3` Adds an upgrade testnet workflow (#1049)
    * `f8230934` Adds eth2 network (#1050)
    * `37f52830` Docker image size reduction (#1068)
    * `6f621c0b` No need to manage a mempool on non-sequencer nodes (#1064)
    * `d8f3d64d` Temporarily allow product key sealing edb credentials for testnet upgrades (#1065)
    * `c254ac9c` Removed stuff related to the old bridge (#1057)
    * `c775d8bd` Store batches from rollups. check batches at same height are consistent. (#1063)
    * `69fcdd10` Wallet extension now supports --host cli arg. (#1061)
    * `c5635a75` Fix solidity compilation warnings for unused calldata. warnings are now treated as errors. (#1060)
    * `fde6606c` Adds todo. (#1058)
    * `72c062b8` Requires the full set of batches to be specified when creating a rollup from a set of batches (#1056)
    * `5b7ffbb5` Each rollup references its batches, rather than a set of transactions (#1055)
    * `12229781` Adds comments to clarify some of the log events logic. (#1054)
    * `1a16f9ef` Fix for ci deployment. (#1053)
    * `3c977786` Rollup header points to batch. (#1052)
    * `556a5741` Avoid unnecessary updating of head rollup for block (#1051)

## January 2023-01-17 (v0.9)
* A list of the PRs merged in this release is as below;
    * `35054c64` In mem tests now use containers (#1045)
    * `b736353f` Added layer 2 deployment and updated bits around it (#1044)
    * `d9ef7b8a` Fixes to some hh deployment issues (#1043)
    * `ce18d3fd` Propagates container start error (#1042)
    * `bb1bc76c` Remove gas constants (#1041)
    * `80fcc27d` Hardhat deploy plugin and deployment of the layer 1 bits (#1038)
    * `070c79d7` Added tasks for the wallet extension (#1037)
    * `a0eb02c3` Testnet and dev tesnet now deploy from one workflow (#1035)
    * `10838a4e` Improved validation of incoming rollups (#1040)
    * `47ecb855` Extracted abigen task as a separate one (#1036)
    * `e375f0f5` Obscuroscan now deploys manually (#1033)
    * `d37f325a` Separate batch and rollup in block submission response (#1034)
    * `0d32caf7` Testnet metrics are pushed to datadog (#1025)
    * `35c60d39` Fixed a slight bug in the message bus (#1028)
    * `e30dcb2a` Changed runner from self hosted to ubuntu (#1030)
    * `a3567c0e` Parallel build the images (#941)
  
## January 2023-01-04 (v0.8)
* Predominantly internal changes as part of work on faster finality, persistent L1 and updated bridge. No user 
  visible changes or breaking API changes are made. 
* A list of the PRs merged in this release is as below;
    * `7652b2c5` Fix for issue deploying testnet (#1020)
    * `84bd5d82` Added bridge smart contracts and test. (#993)
    * `8a75bd05` Deployments wait until node 0 is healthy (#1018)
    * `383b9c7b` Updated changelog (#1016)
    * `3622359f` Uses known registry address (#1017)
    * `6f4f876c` Decouples host and rpc server (#1014)
    * `81d92fb3` Locks edb version for the obscuro node (#1012)
    * `1bb248d3` Remove unnecessary start method from enclave interface (#1011)
    * `7e525447` Rpc server is injected into the host (#1010)
    * `7c51dcb6` Docker updates (#1009)
    * `8ac39307` Updates pccs url for edgelessdb + ego (#1008)
    * `74cc6e9b` Change sequencer id var to read from gh secrets (#1006)
    * `d74c2995` Use a curl command to request obx (#1007)
    * `dd9b5e77` Adds metrics + update p2p to use metrics (#1002)
    * `dd6d7ec6` Fix timing issue with sequencer secret (#1004)
    * `80a24e61` Genesis rollup agg field should be the sequencer that produces it (#1003)
    * `d777c4c1` Updates the sgx folder to match the edgelessdb usage (#1000)
    * `521b0156` Add is() support to blockrejecterr for standard errors (#1001)
    * `3ebd06b7` Switching genesis to be a batch, not a rollup (#997)
    * `c3907cb5` Hooks sequencerid flag (#999)
    * `709cd907` Fixes two bugs in batch catchup (#995)
    * `8b1344e1` Fix for oz 4.5.0 dependency (#994)
    * `38ab4925` Return genesis as batch, not rollup (#990)
    * `25e3b73c` Remove ides folders (#992)
    * `e05282e1` Changed contracts to use hardhat for compilation. (#965)
    * `cb561316` Reduce sleep time after unexpected block provider error (#989)
    * `81f52116` Add fast-fail mechanism to break out of retry (#988)
    * `43dcf872` Create two separate header structs for batch and rollup (#985)
    * `b34fde07` Tx injector l1 deposits are estimated (#986)
    * `c40573be` Clean up batch validation on validators; check for missing batches in chain (#984)
    * `4f0cc501` Stop using same struct for rollup and batch headers, so they can diverge (#981)
    * `79bf6b52` Fix main build errors (#982)
    * `12799f2c` Temporarily disable the validation of sequencer signatures (#980)
    * `5e72c139` P2p healthcheck (#962)
    * `c3aaeb24` Check that rollups processed from l1 blocks are sequential (#977)
    * `c2a463c2` Returns clearer error if rollup signature cannot be validated (#978)
    * `37c9246e` Some minimal validation of processed rollups (#976)
    * `5cab9cf6` Send transactions to sequencer, instead of broadcasting them (#975)
    * `e22643db` Add section into pr template (#974)
    * `a6c8415e` Proper handling of genesis batch on the host side (#972)
    * `997737e5` Reintroduce storage of rollups; store genesis rollup and rollups from blocks (#973)
    * `c80b5d56` Check that received batch is produced by sequencer (#971)
    * `1c06e9ec` Removes use of transaction blob crypto in rollup chain. (#970)
    * `2bbde79f` Rename methods on enclave to reflect fact that batches are now source of truth (#969)
    * `061febd8` Fixed broken link. (#967)
    * `1aab5f55` Clean up enclave storage (#968)
    * `21abaa9c` Cleaner production of rollup on sequencer (#964)
    * `7824f769` Produce block submission response outside of rollupchain; some cleanup (#963)
    * `d6199e1f` Clarify `storenewheads` logic (#961)
    * `650f727c` Avoid storing batches twice; small optimisations to batch catch-up (#959)
    * `a0141139` Disabled cross chain messages block hash binding in management contract. (#958)
    * `3ed4b90d` Separate di constructor for container in testing (#957)
    * `7506ebe3` Reenable withdrawals check in sim. (#956)
    * `304d33ab` Remove storage of rollups on the host side (#955)
    * `755a61ca` Create rollups from scratch instead of retrieving them from l1 blocks (#946)
    * `4fd3825a` Hosts now use the addr pk for deployment (#935)
    * `5b6b6380` Gas estimation centralization (#954)
    * `188d1f23` Logging improvements to easily trace what is going on and a band aid on the test to stop showing false positives. (#953)
    * `a3a4f12c` Second attempt at fixing testnet (#952)
    * `b1574942` Fix for the testnet deployment (#951)
    * `838067f8` Standardise container wiring for enclave and host (#949)
    * `5e18b286` Fix security issue. seal secrets with the unique sgx key (#950)
    * `e8e34ac8` Adding stop and status scripts (#948)
    * `13d334da` Obscuro cross chain messaging (#817)
    * `2d327cee` Adds a start node doc (#945)
    * `cfdfd07c` Submit each batch to the enclave if it's successfully stored (#947)
    * `24871c16` Removes issequencerenclave. updates node types. (#944)
    * `9ebcf3f3` Fixes underflow bug. (#942)
    * `8dca53c8` Removes unused flags (#934)
    * `7f854124` Clean up handling of genesis block (#939)
    * `9bb34347` Reverts use of retry.do. (#940)
    * `e27d1242` Pedro/fix get balance (#937)
    * `41580a89` Fixes batch catchup, and fixes bug in sequencer rollup production (#927)
    * `640f065f` Have shorter timeout for awaiting receipt in in-mem sim (#925)
    * `8334dad0` Host id is now generated from pk (#933)
    * `8c06ad68` Quote the owner in the json (#932)
    * `687f910d` Fixes obscuro scan git deploy (#931)
    * `f75f77d4` Remove pkaddress flag (#930)
    * `86c61ec9` Add in user login (#929)
    * `6b302f85` Fix in-mem mock concurrency bug (#928)
    * `53fb50a4` Remove genesis block references outside of in-mem mock code (#922)
    * `2dade86f` Testnet dns now point to node1 (#918)
    * `df5703b1` Avoid duplication in ancestor-checking methods (#920)
    * `84a4a13d` Jamescarlyle webapp obscuroscan (#849)
    * `3d212dad` Allow l1 start block to be configured to avoid all l1 history (#917)
    * `a979dca4` Fix fork block loop and re-merge blockprovider pr (#916)
    * `c4c30343` Removes buggy fetchheadrollup method on storage (#914)
    * `780af91a` Has-get pattern in the db fix (#912)
    * `b695a1a5` Simplifies storage of new heads (#910)
    * `ab28a79b` Remove `headsafterl1block` type; some clean-up (#907)
    * `d3898965` Eth call, estimategas and getbalance calls now respect the block numb… (#902)
    * `483d5ba2` Revert call error not correctly propagated (#905)
    * `437a53df` Reorganises `rollupchain` methods (#904)
    * `f19917f1` Clean up of l1 block submission (#903)
    * `61da560f` Unlink production of rollups from submission of l1 blocks (#899)
    * `da812862` Clarifies fields in `blockstate` (#901)
    * `131b5e4c` Log error if no peers list yet to send batch request to (#900)
    * `9a5f5812` Separate out genesis response from blocksubmissionresponse (#898)
    * `e8e1de09` Clean up `submitl1block` (#897)
    * `1ca96d95` Allow hosts to catch up on missed batches (#887)
    * `ef1ad41d` Revert "use block provider as source of l1 blocks in host (#891)" (#895)
    * `ab23a792` Adds tests to enclave get balance (#894)
    * `2e59cdbd` Use block provider as source of l1 blocks in host (#891)
    * `739b34c1` Reworks confusing error block. (#893)
    * `82528e73` Rollup chain no longer handles encryption and param validation (#892)
    * `75843a00` Fix enclave tests (#890)
    * `abe4ac09` Simplify blockprovider process control to use context (#889)
    * `7ce5047a` Tweaking the faucet to allow for enclave unit tests (#888)
    * `8ad27bd5` Surface enclave errors. (#886)
    * `80b0fef6` Fixes npe when receipt takes longer than expected time (#881)
    * `d23388e4` More enclave db error handling (#883)
    * `41c2e62c` Profiler - address gosec warning on new versions of golangci-lint (#884)
    * `6d116689` Make our errnotfound match ethereum's (#885)
    * `842997c8` Surface errors from accessors_receipts and accessors_metadata (#882)
    * `ddc39783` Surface enclave errors. (#880)
    * `9ccd7cde` Removes dead linteres to remove noise. (#878)
    * `08e16959` Surface more enclave errors (#877)
    * `4641fd93` Rename to be consistent (#876)
    * `8c993fec` More surfacing of enclave db errors (#874)
    * `9e497728` Add file to deploy obscuroscan into dev-testnet (#875)
    * `72fde61d` Update faucet.md (#872)
    * `85ce478b` Extends integration tests of obscuroscan. switch obscuroscan api to return batches, not rollups. (#873)
    * `325d37fb` Tidy up names (#860)
    * `80a77b7e` Surface enclave db errors (#866)
    * `d28d3b27` Revert "testnet dns now point to node1" (#868)
    * `b40a7517` Test obscuroscan's getlatesttxs in integration tests (#871)
    * `84d42fda` Fix dead links in docs (#864)
    * `bea17905` Adds integration test of gettotaltxs. (#870)
    * `e6598799` Return errors from enclave db, instead of ignoring or using a critical log message (#863)
    * `a734947d` Custom error is now a pointer (#865)
    * `6bcadc77` Switch over various obscuroscan api methods to be based on batches, not rollups (#857)

## November 2022-11-22 (v0.7)
  * A variety of stability related issues are fixed within this release. 
  * Inclusion of a health endpoint for system status monitoring. 
  * It is now possible to run an Obscuroscan against a locally deployed testnet. For more information see 
    [building and running a local testnet](https://github.com/obscuronet/go-obscuro/blob/main/README.md#building-and-running-a-local-testnet) 
    in the project readme.
  * Obscuroscan's GitHub Actions [deploy script](https://github.com/obscuronet/go-obscuro/blob/main/.github/workflows/manual-deploy-obscuroscan.yml) has been modified to run the public Testnet Obscuroscan as an Azure web app. This allows access via HTTPS (TLS), which allows app developers to call the Obscuroscan API from other web apps.
  * A list of the relevant PRs addressed in this release is as below;
    * `12a04c40` Checks whether the head rollup is nil (#859)
    * `619d39b4` Clarify that blocks are L1 blocks (#858)
    * `01884de0` Removes endpoint to get L1 height from Obs node (#856)
    * `9b975f3d` eth_getBlockByNumber and eth_getBlockByHash responds based on the batches, not the rollups (#855)
    * `87588e54` Stores batch at the correct point (#854)
    * `f4d37f6e` Remove geth EVM trace logger (#853)
    * `fcc02555` Distribute and store batches (#850)
    * `243f7ef7` Replace panics with logger.Crit in the enclave (#844)
    * `5f97c1a4` Returns errors from DB methods, instead of `found` bools and critical log events (#842)
    * `3dd03cdc` Uses a write lock instead of a read lock (#847)
    * `c039df7d` Locks the subscription list in LogEventManager for threadsafety (#846)
    * `b1bfed47` Gets number of subs in threadsafe way (#845)
    * `76fde61d` Revert "Fix EVM error casting to use pointer variable" (#841)
    * `d225a75c` Adds methods to host DB for batches (#837)
    * `f3d60127` Fix EVM error casting to use pointer variable (#840)
    * `8e21374b` Fix issues with submit block errors (#838)
    * `ddecd719` Fixes concurrency bug in subscription manager (#839)
    * `12f34d46` Create blockprovider and use it for awaitSecret (#813)
    * `9e524e7f` Removes HeaderWithHashes type (#836)
    * `756b7c16` Removes ExtRollup/ExtRollupWithHash split (#835)
    * `17940c7b` Fixing node start out of sync (#832)
    * `81b8d9c8` Testnet DNS now point to node1 (#827)
    * `de9dbc6f` Cleans up the GetLatestTransactions API method (#833)
    * `6932e020` Fixes grabbing a rollup via ObscuroScan (#829)
    * `c9e978f0` Adds booleans to DB methods to indicate whether was found. (#831)
    * `849ea7aa` Fetch latest Rollup Head now returns error (#826)
    * `bc652690` Adding health check endpoint (#825)
    * `8f049ff9` Handle all errors from ethcall and estimate gas (#823)
    * `9107d571` Fixes eth call error propagation (#822)
    * `976d872c` Remove unused test APIs. Rename RPC method constants for clarity (#821)
    * `3a6f197f` Stop in-mem nodes properly. Prune unused in-mem RPC methods (#820)
    * `84e7c615` Provides logger for Obscuroscan (#819)
    * `088d8f50` Dynamic estimate gas (#815)
    * `4478ffbd` Fix the bridge address to pass the checksums (#812)
    * `ef0e04d9` Downgrade the spammy log message (#810)
    * `6cb0d85a` Have sims test the eth_blockNumber endpoint (#809)
    * `d83c201e` Confusing description of `DB.writeRollupNumber`. Minor clean-up (#791)
    * `5faab414` Fix to use the dev build of the contract deployer (#807)

## November 2022-11-08 (v0.6)
  * The Number Guessing Game has been removed from static and auto deployment scripts, and is now hosted 
    [in a sample applications repository](https://github.com/obscuronet/sample-applications). Given the move for 
    Testnet to be long-running (or at least restartable without contract disappearance), the Guessing Game must be 
    persisted across software updates, and redeployed manually if needed in the same way other applications are.
  * The list of sensitive RPC API methods, where the request and response is encrypted in transit, now covers 
    `eth_call`, `eth_estimateGas`, `eth_getBalance`, `eth_getLogs`, `eth_getTransactionByHash`, `eth_getTransactionCount`, 
    `eth_getTransactionReceipt` and `eth_sendRawTransaction`. See the Obscuro
    [documentation](https://docs.obscu.ro/api/sensitive-apis/) for more details. 
  * Calls to wait for a transaction receipt are now blocking, whereas previously they would return an error meaning the
    client side code needed to perform a specific wait and poll loop. The example on how to [programmatically deploy
    a contract](https://docs.obscu.ro/testnet/deploying-a-smart-contract-programmatically/) has been updated accordingly.
  * The ability to start a faucet server against a local testnet deployment is now supported via a docker 
    container. For more information see the Obscuro 
    [readme](https://github.com/obscuronet/go-obscuro#building-and-running-a-local-faucet).
  * Updates to the [Events](https://github.com/obscuronet/go-obscuro/blob/main/design/ux/Events_design.md) design 
    inclusion of the [Fast Finality](https://github.com/obscuronet/go-obscuro/blob/main/design/finality_protocol/fast_finality.md) design.
  * The [Obscuro docs site](https://docs.obscu.ro/) is now searchable. 
  * Testnet is now officially termed `Evan's Cat`.

* ObscuroScan:
  * ObscuroScan supports a single API at [/rollup/](http://testnet.obscuroscan.io/rollup/) which allows web clients to 
    access a JSON representation of rollups and encrypted transactions. Further details 
    [here](https://docs.obscu.ro/testnet/obscuroscan)

## October 2022-10-21 (v0.5)
* Event Subscriptions:
  * Event subscriptions for logs are now supported via the eth_subscribe and eth_getLogs approaches. This has been 
    tested using both the ethers and web3js libraries. Note that eth_newFilter is not currently supported. For more 
    information see [the events design](https://github.com/obscuronet/go-obscuro/blob/main/design/Events_design.md).

## September 2022-09-22 (v0.4)
* Wallet extension:
  * The wallet extension now supports separate ports for HTTP and WebSocket connections. Use the `--port` and `--portWS` 
    command line options respectively for each. For more information see the
    [Wallet extension](https://docs.obscu.ro/wallet-extension/wallet-extension) documentation. 
* Event subscription:
  * An early preview of event subscriptions is available in this release, though note that this is still undergoing 
    testing and feature enhancements and therefore is liable to issues and instability. For more information on the 
    functionality available reach out to the development team on the discord 
    [active testnet developers](https://discord.com/channels/916052669955727371/1004752710077259838) channel. 
* Transaction receipts:
  * Only return receipts for transactions which were included in a canonical rollup.

## September 2022-09-07 (v0.3)
* Tokens / ERC20 contracts
  * The ERC20 'HOC' and 'POC' tokens are now funded with 18 decimal places of precision. Previously funding of 50 
    tokens was erroneously made as 50 10^-18. This means tokens imported into Metamask will display correctly. Note that
    the number guessing game pay to play still costs 1 10^-18 HOC tokens.
* Wallet extension:
  * Viewing keys are now persisted across wallet extension restarts
  * Enhanced logging for registering of viewing keys

## August 2022-08-22 (v0.2)
* Account balances:
  * Added correct calculation of account balances (previously, all accounts were allocated infinite funds).
* Tokens / ERC20 contracts
  * The two pre-installed ERC20 contracts deployed are now named 'HOC' and 'POC', replacing the previous tokens of 'JAM' 
    and 'ETH'. Contract addresses remain the same as before respectively. The tokens have restricted `balanceOf` and 
    `allowance` calls such that only the owner of the account can view details which should be private to them. See 
    `go-obscuro\integration\erc20contract\ObsERC20.sol` for more information. 
  * Testnet now supports a faucet to distribute native OBX on request. Previously pre-funding of accounts meant that 
    no native tokens were required to execute transactions on Obscuro - this is now not the case and native tokens 
    must be requested. Allocation of native OBX, along with HOC and POC tokens is currently not supported automatically 
    and a request to Obscuro Labs should be made on the Faucet Discord channel.
* Gas prices:
  * The node operator can configure the minimum gas price their aggregator will accept on startup.
* Wallet extension 
  * The wallet extension now supports multiple viewing keys through a single running instance. For more information see 
    the [wallet extension design doc](https://github.com/obscuronet/go-obscuro/blob/main/design/ux/wallet_extension.md), 
    specifically `Handling eth_call requests` for how the required viewing key is selected when multiple keys are 
    registered.

## August 2022
* Testnet launch:
  * Testnet preview launched to limited number of application developers.
  * ObscuroScan block explorer for Testnet launched.
  * Number Guessing Game smart contract deployed to Testnet.
* Obscuro Docsite launched.
* Account balances:
  * Added correct calculation of account balances (previously, all accounts were allocated infinite funds).
  * Introduced network faucet account.
  * Obscuro enclaves services can configure the minimum gas price they'll accept
* ``block.difficulty`` will return a true random number generated inside the secure enclave.
