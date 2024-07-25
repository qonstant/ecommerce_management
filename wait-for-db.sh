#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

# Timeout in seconds
timeout=30

# Start time
start_time=$(date +%s)

while ! pg_isready -h "$host" -p "5432"; do
  current_time=$(date +%s)
  elapsed_time=$((current_time - start_time))
  
  if [ $elapsed_time -ge $timeout ]; then
    >&2 echo "Postgres is unavailable - timeout after $timeout seconds"
    exit 1
  fi
  
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
