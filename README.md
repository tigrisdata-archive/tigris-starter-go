# Tigris Getting Started Golang Application

The code in this repo shows how to integrate Tigris with the backend of
microservice architecture application, using simplified eCommerce service implementation.

For more information please refer to: [Tigris documentation](https://docs.tigrisdata.com)

## Clone the repo

```shell
git clone https://github.com/tigrisdata/tigris-starter-go.git
cd tigris-starter-go
```

## Install Tigris CLI

### macOS
```shell
curl -sSL https://tigris.dev/cli-macos | sudo tar -xz -C /usr/local/bin
```

### Linux
```shell
curl -sSL https://tigris.dev/cli-linux | sudo tar -xz -C /usr/local/bin
```

## Start local Tigris instance
```shell
tigris local up
```

## Compile and start the application
```shell
go build .
./tigris-starter-go
```

## Test the application in new terminal window

Note: For repeatability of the following example, document IDs are hardcoded.
To generate random ID, just remove "_id" field from the request.

### Insert users

Run following commands to create two users: Jane and John

```shell
curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"Name":"John","Balance":100,"_id":"11111111-1111-1111-1111-111111111111"}'
curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"Name":"Jane","Balance":200,"_id":"22222222-2222-2222-2222-222222222222"}'
```

### Insert products

Run the following commands to insert two products: Avocado and Gold

```shell
curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"Name":"Avocado","Price":10,"Quantity":5,"_id":"11111111-1111-1111-1111-111111111111"}'
curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"Name":"Gold","Price":3000,"Quantity":1,"_id":"22222222-2222-2222-2222-222222222222"}'
```

### Place some orders

#### Low balance

The next order will fail because John is trying to purchase 1 unit of Gold which costs 3000,
while John's balance is 100.

```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"UserId":"11111111-1111-1111-1111-111111111111", "Products" : [{"_id":"22222222-2222-2222-2222-222222222222","Quantity":1}],"_id":"11111111-1111-1111-1111-111111111111"}'
```

#### Low stock

The next order will fail because Jane is trying to purchase 30 Avocados which costs 300, while
Jane's balance is 200.

```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	-d '{"UserId":"11111111-1111-1111-1111-111111111111", "Products" : [{"_id":"11111111-1111-1111-1111-111111111111","Quantity":30}],"_id":"11111111-1111-1111-1111-111111111111"}'
```

#### Successful purchase

The next order succeeds because John is purchasing 5 Avocados, which costs 50 and
John's balance is 100, which is enough for the purchase.

```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"UserId":"11111111-1111-1111-1111-111111111111", "Products" : [{"_id":"11111111-1111-1111-1111-111111111111","Quantity":5}],"_id":"11111111-1111-1111-1111-111111111111"}'
```

### Check the balances and stock

Now check that John's balance and Avocado stock is changed accordingly.

```shell
curl localhost:8080/users/read/11111111-1111-1111-1111-111111111111
curl localhost:8080/products/read/11111111-1111-1111-1111-111111111111
curl localhost:8080/orders/read/11111111-1111-1111-1111-111111111111
```

# License

This software is licensed under the [Apache 2.0](LICENSE).
