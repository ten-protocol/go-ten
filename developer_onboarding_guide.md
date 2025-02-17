# Developer Onboarding Guide

## Prerequisites
- [go](https://go.dev/doc/install) (version > 1.20.4)
- [gofumpt](https://github.com/mvdan/gofumpt)
- [golangci-lint](https://golangci-lint.run/) (version 1.52.2)

# Update Setup quicksheet
- Install golangci-lint :
  `curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2`
- Install go1.20.4 : 
  [Multiversion install](https://go.dev/doc/manage-install) example:
```
$ go install golang.org/dl/go1.20.4@latest
$ go1.20.4 download
# update your .zshrc script so that go is in the PATH
$ export PATH="/Users/pedro/sdk/go1.20.4/bin/:/Users/pedro/go/bin:$PATH"

```

## Creating New PR and Getting It Ready to Be Merged to the Main Branch

To start with, you need to create a new branch where you'll develop the new features. Generally, we use a branch naming convention that follows the format:  `<your_name>/<feature_description>`.

You can use an IDE of your choice. Most of the team uses Goland for which you can request a license from Gavin.

Once your branch is pushed to GitHub, you should open a Pull Request (or Draft Pull Request if you  need some preliminary opinion about the approach, or you want to run the tests). This way, other team members can review your code and approve it upon completion.

Whenever new commits are pushed to Github, a pipeline is initiated that runs linters, tests, and so forth.

Before pushing your commits to Github, it's recommended to run the following commands. This ensures that your code is properly formatted and any bugs are fixed:

```
gofumpt -l -w .
golangci-lint run --verbose --fast
```

When your branch is ready for review, another developer must approve the changes. Only then can it be merged into the `main` branch. For merging, we normally use `Squash and merge` (it appears in dropdown once PR is approved).
It's recommended to delete your branch from Github after it has been merged.

## Running End-to-End Tests

End-to-end tests are executed after new commits are merged into the `main` branch. If the tests fail, you'll be notified in the `#testnet-continuous-integration` Discord channel.

If you're making substantial changes that could affect E2E tests, you can manually run these tests on your branch.

To do so, follow these steps:

- Go to the [ten-test](https://github.com/ten-protocol/ten-test/actions) repository
- Click on the `Run local tests` [link](https://github.com/ten-protocol/ten-test/actions/workflows/run_local_tests.yml)
- Click `Run workflow`, enter the name of your branch, and then click `Run workflow`

## TEN smart contracts

Documentation for TEN smart contracts is available [here](https://github.com/ten-protocol/go-ten/blob/main/contracts/README.md).