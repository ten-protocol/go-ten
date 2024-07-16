#!/bin/bash

log_processes() {
  local process=$1
  echo "Current $process processes:"
  pgrep -a -f $process
}

# Function to forcefully kill processes by name
force_kill_processes() {
  local processes=("$@")
  for process in "${processes[@]}"; do
    echo "Attempting to terminate all $process processes"

    log_processes $process

    pkill -TERM -f $process
    sleep 5

    pkill -KILL -f $process
    sleep 3

    for i in {1..5}; do
      if ! pgrep -f $process > /dev/null; then
        echo "All $process processes terminated successfully"
        break
      else
        echo "Reattempting to kill remaining $process processes"
        pkill -KILL -f $process
        sleep 3
      fi
    done

    if pgrep -f $process > /dev/null; then
      echo "Failed to terminate all $process processes"
      exit 1
    fi

    log_processes $process
  done
}

processes=("geth" "beacon-chain")

force_kill_processes "${processes[@]}"

echo "All specified processes terminated successfully"
exit 0