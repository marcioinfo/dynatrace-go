#!/bin/bash

cd /app 

swag init -g ./adapters/http/server.go --parseDependency --parseInternal

exec air