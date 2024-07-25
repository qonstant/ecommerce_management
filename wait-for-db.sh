#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

# Timeout in seconds
timeout=30

# Start time
start_time=$(date +%s)

# Run the command in the background and get its PID
(
  while ! pg_isready -h "$host" -p "5432"; do
    >&2 echo "Postgres is unavailable - sleeping"
    sleep 1
  done
  >&2 echo "Postgres is up - executing command"
) &

# PID of the background process
bg_pid=$!

# Wait for either the timeout or the background process to complete
while kill -0 "$bg_pid" 2> /dev/null; do
  current_time=$(date +%s)
  elapsed_time=$((current_time - start_time))
  
  if [ $elapsed_time -ge $timeout ]; then
    >&2 echo "Timeout after $timeout seconds - exiting"
    kill -9 "$bg_pid" 2> /dev/null
    exit 1
  fi
  
  sleep 1
done

# Execute the command if the background process completed successfully
exec $cmd
