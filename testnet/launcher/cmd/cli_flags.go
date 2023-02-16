package main

// Flag names.
const (
	numberNodesFlag = "number_nodes"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		numberNodesFlag: "Number of obscuro nodes to deploy",
	}
}
