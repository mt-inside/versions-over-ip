#!/usr/bin/env bash

( go run cmd/server/server.go ) & pid=$!

go run cmd/client/client.go

pgid=$(cat /proc/$pid/stat | awk '{ print $5 }')  # relies on process name not having spaces
echo "killing group $pgid (under parent $pid)"
kill -- -$pgid # `go run` forks, so this kills the whole process group
