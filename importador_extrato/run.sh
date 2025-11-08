#!/bin/bash

# Database configuration
export DB_HOST="windows-host"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="postgres"

# Build and run the importer
go build -o importador
./importador
