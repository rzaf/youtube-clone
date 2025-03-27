#!/bin/sh

cd /app

BINARY_PATH="$1"
SERVICE=$(basename "$BINARY_PATH")

if [ -z "$BINARY_PATH" ]; then
  echo "Error: No binary path provided."
  exit 1
fi

# Function to get the checksum of the binary
get_checksum() {
  if [ -f "$BINARY_PATH" ]; then
    sha256sum "$BINARY_PATH" | awk '{print $1}'
  else
    echo ""
  fi
}

# Initial checksum
CHECKSUM=$(get_checksum)
PID=""

# Start the binary function
start_binary() {
  echo "Starting the $SERVICE..."
  "$BINARY_PATH" &
  PID=$!
  echo "Service started with PID $PID"
}

# Stop the binary function
stop_binary() {
  if [ ! -z "$PID" ]; then
    echo "Stopping the $SERVICE..."
    kill "$PID"
    wait "$PID"
    PID=""
  fi
}

# Start the binary initially
start_binary

# Infinite loop to monitor the binary for changes
while true; do
  NEW_CHECKSUM=$(get_checksum)

  if [ "$CHECKSUM" != "$NEW_CHECKSUM" ]; then
    echo "$SERVICE has changed, restarting the service..."
    CHECKSUM="$NEW_CHECKSUM"
    
    # Stop the running binary and restart it
    stop_binary
    start_binary
  fi

  # Sleep for a few seconds before checking again
  sleep 5
done
