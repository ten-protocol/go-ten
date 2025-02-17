These packages contain code related to the simulation and testing an end-to-end TEN network.

To include the Docker tests when running the tests, build the Docker images using the instructions in the 
`dockerfiles/` folder, then run the tests with the `docker` tag (e.g. `go test -v -tags docker ./...`).
