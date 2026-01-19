#!/bin/bash

go mod download

sudo apt-get update
sudo apt-get install -y postgresql postgresql-contrib

sudo systemctl start postgresql || true
sudo systemctl enable postgresql || true

# Ожидание запуска PostgreSQL
sleep 3

sudo -u postgres psql -f migrations/1_initilize_schema.sql
sudo -u postgres psql -f migrations/2_alter_tables.sql
sudo -u postgres psql -f migrations/3_grant_access.sql
