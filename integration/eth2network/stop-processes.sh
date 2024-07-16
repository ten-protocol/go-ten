#!/bin/bash

force_kill_processes() {
  local processes=("$@")
  for process in "${processes[@]}"; do
    echo "Attempting to terminate all $process processes"

    pkill -TERM -f $process
    sleep 2

    pkill -KILL -f $process
    sleep 1

    for i in {1..5}; do
      if ! pgrep -f $process > /dev/null; then
        echo "All $process processes terminated successfully"
        break
      else
        echo "Reattempting to kill remaining $process processes"
        pkill -KILL -f $process
        sleep 1
      fi
    done

    if pgrep -f $process > /dev/null; then
      echo "Failed to terminate $process processes"
      exit 1
    fi
  done
}

processes=("geth" "beacon-chain")

force_kill_processes "${processes[@]}"

echo "All specified processes terminated successfully"
exit 0