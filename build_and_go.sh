#!/usr/bin/env sh

echo "Building geoGIT.go..."
cd /app/
go get github.com/lib/pq
go build geoGIT.go

echo "Running geoGIT..."
./geoGIT
