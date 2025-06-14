#!/bin/sh

set -e 
echo "run db migration..."
./migrate -path ./migration -database "$DB_SOURCE" -verbose up

echo "start the app..."
exec "$@"