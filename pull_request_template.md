### Why is this change needed?

- Provide a description and a link to the underlying ticket

### What changes were made as part of this PR:

- Provide a high level list of the changes made

### Definition of done

- [ ] Unit tests added to cover new or changed functionality 
- [ ] Docs pages updated to cover new or changed functionality
- [ ] [Changelog.md](https://github.com/obscuronet/go-obscuro/blob/main/docs/testnet/changelog.md) updated 
- [ ] End-to-end tests run if required (see below)

### End-to-end tests

Running the end-to-end tests in the [obscuro-test repo](https://github.com/obscuronet/obscuro-test) is currently at the 
discretion of the developer prior to merging a PR. The intention is not to slow down development by having the longer 
running end-to-end tests as a PR gate run on each branch commit etc. To manually trigger a run pre-merge;

- Go to [run_local_tests](https://github.com/obscuronet/obscuro-test/actions/workflows/run_local_tests.yml)
- Click the "Run workflow" button
- Use the main branch of obscuro-test 
- Enter the name of the PR branch for go-obscuro
- Run the workflow

The run takes approximately 20 minutes and any failure will be reported in the workflow output. Should the tests fail 
the docker container output for the local testnet and all test artifacts will be added to the run and can be downloaded
for debugging and analysis, with a retention period of 2 days. 

Note that every PR merge to main will also trigger a run of post-merge tests, so even if you do not manually trigger 
pre-merge, tests will be run post-merge. The intention being to catch any issues fast on main to allow for quick 
resolution. The post-merge test output can be seen at 
[run_merge_tests](https://github.com/obscuronet/obscuro-test/actions/workflows/run_merge_tests.yml) and will show both 
the PR number and author in the list of runs. 


