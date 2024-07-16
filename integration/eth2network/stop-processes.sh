#!/bin/bash

# Function to kill processes using specific ports
kill_processes_by_port() {
  local port=$1
  local protocol=$2
  echo "Attempting to kill processes using $protocol port $port"

  if [ "$protocol" = "udp" ]; then
    lsof -i udp:$port | awk 'NR>1 {print $2}' | xargs -r kill -9
  else
    lsof -i tcp:$port | awk 'NR>1 {print $2}' | xargs -r kill -9
  fi

  sleep 3

  if lsof -i $protocol:$port > /dev/null; then
    echo "Failed to kill processes using $protocol port $port"
    exit 1
  else
    echo "Successfully killed processes using $protocol port $port"
  fi
}

ports_to_kill=("12000/udp" "30303/tcp")

# Kill processes by port
for port in "${ports_to_kill[@]}"; do
  IFS=/ read port_num protocol <<< "$port"
  kill_processes_by_port $port_num $protocol
done

echo "All specified processes terminated successfully"
exit 0