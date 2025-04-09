#!/bin/bash

# Create data directory if it doesn't exist
mkdir -p ./data

# Start the blockchain node
./binomena --api-port $PORT --p2p-port 9000 --id "$NODE_ID" --data-dir ./data
