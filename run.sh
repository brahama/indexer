#!/usr/bin/env sh
set -e

echo "Starting indexer in writer mode, connecting to algod node.  Using DB $DATABASE_NAME"

# Start indexer, connecting to node & DB.  Get node API token from shared volume
#./cmd/algorand-indexer/algorand-indexer daemon \
#  -P "host=indexer-db port=5432 user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} sslmode=disable" \
#  --algod-net="http://${ALGORAND_HOST}:8080" --algod-token="$(cat /var/algorand/data/algod.token)"

./cmd/algorand-indexer/algorand-indexer daemon \
  -P "host=${DB_HOST} port=5432 user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" --algod-net="http://${ALGOD_NODE}:4161" --algod-token "${ALGOD_TOKEN}"
