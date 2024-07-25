#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

# Wait for 1 minute (60 seconds) before starting to check if PostgreSQL is ready
sleep 30

until pg_isready -h "$host" -p "5432"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd