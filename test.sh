#!/usr/bin/env bash
set -e
# Schema
cd test_data && graphql-codegen --config codegen.yml && cd ..
# Build
go build
# Gen
./spacex-cli generate --queries test_data/queries/ --schema test_data/schema.json --target client --name ApiClient --output test_data/output/spacex.ts
./spacex-cli generate --queries test_data/queries/ --schema test_data/schema.json --target typescript --output test_data/output/spacex.web.ts