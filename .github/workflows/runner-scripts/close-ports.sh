#!/usr/bin/env bash

# Ensure any fail is loud
set -euo pipefail

lowest_port=8000
highest_port=58000
additional_ports=(80 81 99)

pids=$(lsof -iTCP:$lowest_port-$highest_port -sTCP:LISTEN -t)
if [ -z "$pids" ]; then
  echo "No processes are listening on ports from $lowest_port to $highest_port"
else
  for pid in $pids; do
    echo "Process $pid is listening on one of the ports from $lowest_port to $highest_port"
  done
  for pid in $pids; do
    echo "Killing process $pid on one of the ports from $lowest_port to $highest_port"
    kill -9 $pid || true
  done
fi

echo "Additional ports: ${additional_ports[@]}"
for port in "${additional_ports[@]}"; do
  pids=$(lsof -ti TCP:$port)
  if [ -z "$pids" ]; then
    echo "No processes are listening on port $port"
  else
    for pid in $pids; do
      echo "Killing process $pid on port $port"
      kill $pid || true
    done
  fi
done
