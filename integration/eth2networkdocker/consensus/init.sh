#!/bin/sh

echo "Moving into execution directory"
cd /consensus || exit

mkdir -p beacondata

cp network-keys beacondata/network-keys