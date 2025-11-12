#!/bin/bash

# Database configuration

export DB_NAME="wgtrade"
export DB_USER="postgres"
export DB_PASSWORD="s#Dvptj58Uf!e4&1kkaf"
export DB_PORT="5532"
export DB_HOST="easypanel.wgtradepromo.com.br"

export DB_HOST="172.23.48.1"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="wgtrade"

# Build and run the importer
go build -o importador
./importador
