#!/bin/bash

# Print environment for debugging
echo "Environment variables:"
echo "PORT: $PORT"
echo "NODE_ID: $NODE_ID"

# Start the blockchain node with supported flags
./binomena --api-port $PORT --p2p-port 9000 --id "genesis-node"