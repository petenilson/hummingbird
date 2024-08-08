#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status.

echo "Creating Docker container..."
docker run --name loadtest -p 5432:5432 -e POSTGRES_USER=ledger_user -e POSTGRES_PASSWORD=ledger_password -e POSTGRES_DB=ledger_db -d postgres:16-alpine


echo "Waiting for container..."
sleep 3

echo "Starting webserver..."
./hummingbird &SERVER_PID=$!

echo "Waiting for webserver..."
sleep 2

echo "Creating 1000 test accounts..."
for i in {1..1000}
do
   curl -s -X POST http://localhost:8000/accounts -H "Content-Type: application/json" -d "{\"name\":\"test_account_$i\"}" > /dev/null 2>&1
done
echo "Finished creating test accounts."

echo "Running load test..."
wrk -t12 -c400 -d30s -s wrk_script.lua http://localhost:8000

echo "Stopping server..."
kill $SERVER_PID

echo "Cleaning up Docker container..."
docker stop loadtest
docker rm loadtest

echo "Test completed and environment cleaned up."
