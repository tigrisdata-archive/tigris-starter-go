#!/bin/bash

set -e

go build .
./tigris-starter-go >tigris-starter.log 2>&1 &
PID=$!

sleep 0.5

curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"Name":"John","Balance":100,"_id":"11111111-1111-1111-1111-111111111111"}'
echo
curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"Name":"Jane","Balance":200,"_id":"22222222-2222-2222-2222-222222222222"}'
echo

curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"Name":"Avocado","Price":10,"Quantity":5,"_id":"11111111-1111-1111-1111-111111111111"}'
echo
curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"Name":"Gold","Price":3000,"Quantity":1,"_id":"22222222-2222-2222-2222-222222222222"}'
echo

curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"UserId": "11111111-1111-1111-1111-111111111111", "Products" : [{"_id":"22222222-2222-2222-2222-222222222222","Quantity":1}]}' || true
echo

curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	-d '{"UserId":"11111111-1111-1111-1111-111111111111", "Products" : [{"_id":"11111111-1111-1111-1111-111111111111","Quantity":30}]}' || true
echo

curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"UserId":"11111111-1111-1111-1111-111111111111", "Products" : [{"_id":"11111111-1111-1111-1111-111111111111","Quantity":5}],"_id":"11111111-1111-1111-1111-111111111111"}'
echo

curl localhost:8080/users/read/11111111-1111-1111-1111-111111111111
echo
curl localhost:8080/products/read/11111111-1111-1111-1111-111111111111
echo
curl localhost:8080/orders/read/11111111-1111-1111-1111-111111111111
echo

kill $PID
