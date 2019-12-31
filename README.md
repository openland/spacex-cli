# spacex-cli
Codegen CLI for SpaceX GraphQL

## Usage
spacex-cli generate --queries ./queries/ --schema ./schema.json --target kotlin --package com.spacex.graphql -output ./android/app/src/main/java/com/spacex/graphql
spacex-cli generate --queries ./queries/ --schema ./schema.json --target swift -output ./ios/Operations.swift
spacex-cli generate --queries ./queries/ --schema ./schema.json --target typescript -output ./Operations.ts
