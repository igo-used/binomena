#!/bin/bash

# Print environment for debugging
echo "Environment variables:"
echo "PORT: $PORT"
echo "NODE_ID: $NODE_ID"
echo "DATA_DIR: $DATA_DIR"

# Create data directory if it doesn't exist
if [ -n "$DATA_DIR" ]; then
  echo "Setting up data directory: $DATA_DIR"
  mkdir -p $DATA_DIR
  
  # Create a data subdirectory for blockchain data
  mkdir -p $DATA_DIR/blockchain
  mkdir -p $DATA_DIR/wallets
  
  echo "Data directory created and ready"
fi

# Start the blockchain node with supported flags
./binomena --api-port $PORT --p2p-port 9000 --id "$NODE_ID"