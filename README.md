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

### Insert users

Run following commands to create two users: Jane and John

```shell
curl localhost:8080/users/create \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"Name":"John","Balance":100}'
    
curl localhost:8080/users/create \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"Name":"Jane","Balance":200}'
```

### Insert products

Run the following commands to insert two products: Avocado and Gold

```shell
curl localhost:8080/products/create \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"Name":"Avocado","Price":10,"Quantity":5}'
    
curl localhost:8080/products/create \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"Name":"Gold","Price":3000,"Quantity":1}'
```

### Place some orders

#### Low balance

Let's start off with an order that fails because John is trying to purchase 1
unit of Gold that costs $3000.00, while John's balance is $100.00.

```shell
curl http://localhost:8080/orders/create \
      -X POST \
      -H 'Content-Type: application/json' \
      -d '{"UserId":1,"Products":[{"Id":2,"Quantity":1}]}'
```

#### Low stock

The next order fails as well because Jane is trying to purchase 10 Avocados,
but there is only 5 in the stock.

```shell
curl http://localhost:8080/orders/create \
      -X POST \
      -H 'Content-Type: application/json' \
      -d '{"UserId":2,"Products":[{"Id":1,"Quantity":10}]}'
```

#### Successful purchase

Now an order that succeeds as John purchases 5 Avocados that cost
$50.00 and John's balance is $100.00, which is enough for the purchase.

```shell
curl http://localhost:8080/orders/create \
      -X POST \
      -H 'Content-Type: application/json' \
      -d '{"UserId":1,"Products":[{"Id":1,"Quantity":5}]}'
```

### Check the balances and stock

Now go ahead and confirm that both John's balance and the Avocado stock is
up-to-date.

```shell
curl http://localhost:8080/users/read/1
curl http://localhost:8080/products/read/1
curl http://localhost:8080/orders/read/1
```

### Search
Now, search for users

```shell
curl http://localhost:8080/users/search \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"q":"john"}'
```

Or search for products named "avocado"

```shell
curl localhost:8080/products/search \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{
        "q": "avocado",
        "searchFields": ["Name"]
      }'
```

# License

This software is licensed under the [Apache 2.0](LICENSE).
