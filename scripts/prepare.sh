#!/bin/bash

go mod download

# sudo apt-get update
# sudo apt-get install -y postgresql postgresql-contrib
# sudo systemctl start postgresql || true
# sudo systemctl enable postgresql || true

export PGPASSWORD=val1dat0r
psql -h localhost -p 5432 -U validator -f migrations/1_initilize_schema.sql
psql -h localhost -p 5432 -U validator -f migrations/2_alter_tables.sql
# sudo -U validator psql -h localhost -f migrations/3_grant_access.sql
