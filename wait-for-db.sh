#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

# Wait for 20 seconds
sleep 20

# Execute the command
exec $cmd
