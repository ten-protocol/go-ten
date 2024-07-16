#!/bin/bash

# Function to forcefully kill processes by name
force_kill_processes() {
  local processes=("$@")
  for process in "${processes[@]}"; do
    pkill -9 $process
    sleep 1 # Give the system some time to terminate the processes

    # Verify processes are terminated
    for i in {1..5}; do
      if ! pgrep -f $process > /dev/null; then
        echo "All $process processes terminated successfully"
        break
      else
        echo "Reattempting to kill remaining $process processes"
        pkill -9 $process
        sleep 1
      fi
    done

    if pgrep -f $process > /dev/null; then
      echo "Failed to terminate all $process processes"
      exit 1
    fi
  done
}

# Processes to kill
processes=("geth" "beacon-chain")

# Forcefully kill specified processes
force_kill_processes "${processes[@]}"

echo "All specified processes terminated successfully"
exit 0