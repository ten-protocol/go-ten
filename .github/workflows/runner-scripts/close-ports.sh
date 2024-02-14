#!/usr/bin/env bash

# Ensure unset variables cause an exit and that a pipeline fails if any command fails
set -uo pipefail

echo "start script and set ports vars"

lowest_port=8000
highest_port=58000
additional_ports=(80 81 99)

echo "fetch pids"
pids=$(lsof -iTCP:$lowest_port-$highest_port -sTCP:LISTEN -t)
lsof_exit_status=$?

echo "list pids $pids and process range first"
if [ $lsof_exit_status -eq 0 ]; then
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
elif [ $lsof_exit_status -eq 1 ]; then # lsof returns 1 if no processes are found
    echo "No processes found on ports from $lowest_port to $highest_port, which is expected."
else
    echo "lsof command failed with exit status $lsof_exit_status"
fi
echo "range done"

echo "Additional ports: ${additional_ports[@]}"
for port in "${additional_ports[@]}"; do
    pids=$(lsof -ti TCP:$port)
    lsof_exit_status=$?

    if [ $lsof_exit_status -eq 0 ]; then
        if [ -z "$pids" ]; then
            echo "No processes are listening on port $port"
        else
            for pid in $pids; do
                echo "Killing process $pid on port $port"
                kill $pid || echo "Failed to kill process $pid"
            done
        fi
    elif [ $lsof_exit_status -eq 1 ]; then # lsof returns 1 if no processes are found
        echo "No processes found on port $port, which is expected."
    else
        echo "lsof command for port $port failed with exit status $lsof_exit_status"
    fi
done
