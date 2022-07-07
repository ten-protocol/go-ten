#!/usr/bin/env bash

#
# This script cleans up the environment BEFORE any job run
#
# It's expected to exist in /home/obscuro/scripts/pre-run-cleanup.sh
# Has to be manually deployed, updating the CI will not update the script.
# More details : https://docs.github.com/en/actions/hosting-your-own-runners/running-scripts-before-or-after-a-job#triggering-the-scripts
#


# Ensure any fail is loud and explicit
set -euo pipefail

still_running_containers=$(docker ps -q)

if [[ ${still_running_containers} ]]; then
    # Stop any running container
    echo "Killing containers ${still_running_containers}"
    docker kill ${still_running_containers}
fi

# Clean docker
docker system prune -a -f --volumes

# kill any running application
sudo pkill -f geth
sudo pkill -f enclave
sudo pkill -f host

# kill any running application based on the port
echo 'Ports open before port kill'
sudo ss -pntl

# Disabling the port kill for now
#low=30000;
#high=39999;
#for i in $(seq $low $high); do
#  sudo lsof -i :$i | tail -n +2 | awk '{system("kill -s 9 " $2)}';
#done
#
#echo 'After port kill'
#sudo ss -pntl

echo 'Ready to run job'
