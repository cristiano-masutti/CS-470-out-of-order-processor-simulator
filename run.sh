#!/bin/bash
set -e

# Usage check
if [ "$#" -ne 2 ]; then
  echo "Usage: ./run.sh </path/to/input.json> </path/to/output.json>"
  exit 1
fi

INPUT=$1
OUTPUT=$2

# Run the simulator from the root directory
./simulator "$INPUT" "$OUTPUT"
