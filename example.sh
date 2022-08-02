#!/bin/bash

set -e

# tigris local up

go build .
./tigris-starter-go >tigris-starter.log 2>&1 &
PID=$!

sleep 0.5

# first parameter is path
# second parameter is document to write
request() {
  curl -X POST "localhost:8080/$1" -H 'Content-Type: application/json' -d "$2"
  echo
}

request users/create '{"Name":"John","Balance":100}' #Id=1
request users/create '{"Name":"Jane","Balance":200}' #Id=2

request products/create '{"Name":"Avocado","Price":10,"Quantity":5}' #Id=1
request products/create '{"Name":"Gold","Price":3000,"Quantity":1}' #Id=2

#low balance
request orders/create '{"UserId":1,"Products":[{"Id":2,"Quantity":1}]}' || true
# low stock
request orders/create '{"UserId":1,"Products":[{"Id":1,"Quantity":10}]}' || true

request orders/create '{"UserId":1,"Products":[{"Id":1,"Quantity":5}]}' #Id=1

curl localhost:8080/users/read/1
echo
curl localhost:8080/products/read/1
echo
curl localhost:8080/orders/read/1
echo

# search
request users/search '{"q":"john"}'
request products/search '{"q":"avocado","searchFields": ["Name"]}'

kill $PID

#tigris local down