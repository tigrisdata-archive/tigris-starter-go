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

```shell
curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"Name":"John","Balance":100}'
curl -X POST localhost:8080/users/create -H 'Content-Type: application/json' \
	 -d '{"id":2,"Name":"Jane","Balance":200}'
```

### Insert products

```shell
curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"Name":"Avocado","Price":10,"Quantity":5}'
curl -X POST localhost:8080/products/create -H 'Content-Type: application/json' \
	 -d '{"id":2,"Name":"Gold","Price":3000,"Quantity":1}'
```

### Place some orders

#### Low balance
```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"UserId":1, "Products" : [{"id":2,"Quantity":1}]}'
```

#### Low stock
```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	-d '{"id":1,"UserId":1, "Products" : [{"id":1,"Quantity":1000}]}'
```

#### Successful purchase
```shell
curl -X POST localhost:8080/orders/create -H 'Content-Type: application/json' \
	 -d '{"id":1,"UserId":1, "Products" : [{"id":1,"Quantity":5}]}'
```

### Check the balances and stock

```shell
curl localhost:8080/users/read/1
curl localhost:8080/products/read/1
```

# License

This software is licensed under the [Apache 2.0](LICENSE).
