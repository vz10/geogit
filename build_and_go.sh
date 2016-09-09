#!/usr/bin/env sh

echo "Building geoGIT.go..."
cd /app/
go get github.com/lib/pq
go get github.com/tomnomnom/linkheader
go build geoGIT.go

echo "Running geoGIT..."
./geoGIT
