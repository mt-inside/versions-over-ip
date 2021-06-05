#!/usr/bin/env bash

( go run cmd/server/server.go ) & pid=$!

go run cmd/client/client.go

echo killing $pid
kill -- -$pid # `go run` forks, so this kills the whole process group (the PID of the leader is the PGID of the group)
