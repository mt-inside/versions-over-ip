#!/usr/bin/env bash

( go run cmd/server/server.go ) & pid=$!
go run cmd/client/client.go

echo pid is $pid
ps $pid
echo other "server"s
pgrep server
echo "go run forks?"

echo killing $pid
kill $pid
