# Tech decisions

This document lists the frameworks and tooling used by Ten.

The purpose of this list is to avoid the proliferation of different frameworks and tooling with the same general 
purpose, as this increases the developer overhead when extending Ten.

The introduction of a new framework or tool should be discussed with the team, even if the usage sits outside the 
`go-ten` repo.

## Deployment

### Containerisation

Containers should target [Docker](https://www.docker.com/).

### Cloud

Cloud deployments should target [Azure](https://azure.microsoft.com/en-gb/).

## Testing

### Mocking

Mocking should be performed using [testify](https://github.com/stretchr/testify).

### End-to-end tests

End-to-end tests should be written using [PySys](https://pypi.org/project/PySys/).

## Communications

### RPC

RPC communication should occur over [protocol buffers](https://pkg.go.dev/google.golang.org/protobuf).

## Web

### Servers

Webservers should be written in [Gin](https://github.com/gin-gonic/gin).
