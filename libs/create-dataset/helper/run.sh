#!/usr/bin/env bash

fish -c 'psql -c "DROP DATABASE IF EXISTS postdb;"'
fish -c 'psql -c "create database postdb;"'
rm -rf ./media/ ./err.txt ./img.txt ./img-resize.txt
rm -rf ../../db/prisma/migrations
getPWD=$(echo "$(pwd)")
cd ../../db/
bunx prisma generate
npx prisma migrate dev --name "Initial Migration"
cd "$getPWD"
bun create_dataset.ts
