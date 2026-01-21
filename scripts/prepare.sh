#!/bin/bash

go mod download

export PGPASSWORD=val1dat0r
psql -h localhost -p 5432 -U validator -d project-sem-1 -f migrations/1_initilize_schema.sql
