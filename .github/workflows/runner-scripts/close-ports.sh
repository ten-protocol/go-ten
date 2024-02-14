#!/usr/bin/env bash

# Ensure any fail is loud
set -euo pipefail

# Function to check if a command exists
command_exists() {
    type "$1" &> /dev/null
}

echo "start script and set ports vars"

lowest_port=8000
highest_port=58000
additional_ports=(80 81 99)

# Check if lsof command exists
if ! command_exists lsof; then
    echo "Error: lsof command not found. Please install lsof."
    exit 1
fi

echo "fetch pids"
# Use || true to prevent script exit in case of command failure
echo "cmd: lsof -iTCP:$lowest_port-$highest_port -sTCP:LISTEN -t"
pids=$(lsof -iTCP:$lowest_port-$highest_port -sTCP:LISTEN -t)

# Check if the lsof command was successful
if [ -z "$pids" ] && [ $? -ne 0 ]; then
    echo "lsof command failed. Check if you have the necessary permissions."
    exit 1
fi

echo "list pids $pids and process range first"
if [ -z "$pids" ]; then
    echo "No processes are listening on ports from $lowest_port to $highest_port"
else
    for pid in $pids; do
        echo "Process $pid is listening on one of the ports from $lowest_port to $highest_port"
    done
    for pid in $pids; do
        echo "Killing process $pid on one of the ports from $lowest_port to $highest_port"
        kill $pid || echo "Failed to kill process $pid"
    done
fi
echo "range done"
echo "Additional ports: ${additional_ports[@]}"
for port in "${additional_ports[@]}"; do
    echo "cmd: lsof -ti TCP:$port"
    pids=$(lsof -ti TCP:$port)

    # Check if the lsof command was successful
    if [ -z "$pids" ] && [ $? -ne 0 ]; then
        echo "lsof command failed for port $port. Check if you have the necessary permissions."
        continue
    fi

    if [ -z "$pids" ]; then
        echo "No processes are listening on port $port"
    else
        for pid in $pids; do
            echo "Killing process $pid on port $port"
            kill $pid || echo "Failed to kill process $pid"
        done
    fi
done
