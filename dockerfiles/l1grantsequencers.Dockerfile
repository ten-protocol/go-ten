# --- STAGE 1: Build the binary ---
FROM golang:1.25-alpine AS builder

# Install necessary tools (for CGO or general build needs)
RUN apk add --no-cache build-base

# Enable CGO (can be disabled if not using C deps)
ENV CGO_ENABLED=1

# Set working directory to match your repo root
WORKDIR /home/obscuro/go-obscuro

# Copy mod files early for caching
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy full codebase
COPY . .

# Set working dir to the script's package
WORKDIR /home/obscuro/go-obscuro/testnet/launcher/l1grantsequencers/cmd

# Build the binary (output name is optional but useful)
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o l1grantsequencers

# --- STAGE 2: Slim runtime image ---
FROM alpine:3.18

# Copy built binary into the runtime image
COPY --from=builder /home/obscuro/go-obscuro/testnet/launcher/l1grantsequencers/cmd/l1grantsequencers /usr/local/bin/l1grantsequencers

# Set working directory (optional)
WORKDIR /usr/local/bin

# Default command to run the script
ENTRYPOINT ["/usr/local/bin/l1grantsequencers"]