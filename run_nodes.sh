#!/bin/bash

# Build the blockchain
go build -o binomena

# Start the first node (bootstrap node)
./binomena --api-port 8080 --p2p-port 9000 --id "node-1" &
BOOTSTRAP_PID=$!

# Wait for the bootstrap node to start
sleep 2

# Get the bootstrap node's address
BOOTSTRAP_ADDRESS=$(curl -s http://localhost:8080/status | grep -o '"peers": \[[^]]*\]' | grep -o '/ip4/[^"]*')

# Start additional nodes
./binomena --api-port 8081 --p2p-port 9001 --bootstrap "$BOOTSTRAP_ADDRESS" --id "node-2" &
NODE2_PID=$!

./binomena --api-port 8082 --p2p-port 9002 --bootstrap "$BOOTSTRAP_ADDRESS" --id "node-3" &
NODE3_PID=$!

# Wait for user to press Ctrl+C
echo "All nodes started. Press Ctrl+C to stop all nodes."
wait $BOOTSTRAP_PID

# Kill all nodes when the bootstrap node is stopped
kill $NODE2_PID $NODE3_PID
