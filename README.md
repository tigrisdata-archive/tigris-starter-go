# Tigris Getting Started Golang Application

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

### Insert users

Run following commands to create two users: Jane and John

```shell
curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"Name":"John","Balance":100}'
curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"id":2,"Name":"Jane","Balance":200}'
```

### Insert products

Run the following commands to insert two products: Avocado and Gold

```shell
curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"Name":"Avocado","Price":10,"Quantity":5}'
curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"id":2,"Name":"Gold","Price":3000,"Quantity":1}'
```

### Place some orders

#### Low balance

The next order will fail because John is trying to purchase 1 unit of Gold which costs 3000,
while John's balance is 100.

```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"UserId":1, "Products" : [{"id":2,"Quantity":1}]}'
```

#### Low stock

The next order will fail because Jane is trying to purchase 30 Avocados which costs 300, while
Jane's balance is 200.

```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	-d '{"id":1,"UserId":2, "Products" : [{"id":1,"Quantity":30}]}'
```

#### Successful purchase

The next order succeeds because John is purchasing 5 Avocados, which costs 50 and
John's balance is 100, which is enough for the purchase.

```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"UserId":1, "Products" : [{"id":1,"Quantity":5}]}'
```

### Check the balances and stock

Now check that John's balance and Avocado stock is changed accordingly.

```shell
curl localhost:8080/users/read/1
curl localhost:8080/products/read/1
```

# License

This software is licensed under the [Apache 2.0](LICENSE).
