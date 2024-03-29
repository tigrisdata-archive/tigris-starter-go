package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tigrisdata/tigris-client-go/fields"
	"github.com/tigrisdata/tigris-client-go/filter"
	"github.com/tigrisdata/tigris-client-go/search"
	"github.com/tigrisdata/tigris-client-go/tigris"
)

type User struct {
	Id int32 `tigris:"primaryKey,autoGenerate"`

	Name    string
	Balance float64
}

type Order struct {
	Id int32 `tigris:"primaryKey,autoGenerate"`

	UserId int32

	Products []Product
}

type Product struct {
	Id int32 `tigris:"primaryKey,autoGenerate"`

	Name     string
	Quantity int
	Price    float64
}

func setupReadRoute[T interface{}](r *gin.Engine, db *tigris.Database, name string) {
	r.GET(fmt.Sprintf("/%s/read/:id", name), func(c *gin.Context) {
		coll := tigris.GetCollection[T](db)

		u, err := coll.ReadOne(c, filter.Eq("Id", c.Param("id")))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, u)
	})
}

func setupCRUDRoutes[T interface{}](r *gin.Engine, db *tigris.Database, name string) {
	setupReadRoute[T](r, db, name)
	setupSearchRoute[T](r, db, name)

	r.POST(fmt.Sprintf("/%s/create", name), func(c *gin.Context) {
		coll := tigris.GetCollection[T](db)

		var u T
		if err := c.Bind(&u); err != nil {
			return
		}

		if _, err := coll.Insert(c, &u); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, u)
	})

	r.DELETE(fmt.Sprintf("/%s/delete/:id", name), func(c *gin.Context) {
		coll := tigris.GetCollection[T](db)

		if _, err := coll.Delete(c, filter.Eq("Id", c.Param("id"))); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Status": "DELETED"})
	})
}

// Create an order in a transaction,
// taking into account user balances and product stock.
func setupCreateOrderRoute(r *gin.Engine, db *tigris.Database) {
	r.POST("/orders/create", func(c *gin.Context) {
		var o Order
		// Read the request body into o
		if err := c.Bind(&o); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Perform the read, update and insert on the users, products and orders
		// collections in a transaction. If the function passed to db.Tx
		// returns an error then the transaction will be automatically rolled back.
		// If no error is returned, the transaction will be automatically committed.
		err := db.Tx(c, func(txCtx context.Context) error {
			// Fetch an object of the users collection
			users := tigris.GetCollection[User](db)

			// Read the user with order's UserId
			u, err := users.ReadOne(txCtx, filter.Eq("Id", o.UserId))
			if err != nil {
				return err
			}

			// Fetch an object of the products collection
			products := tigris.GetCollection[Product](db)

			orderTotal := 0.0

			// For every product in the order
			for i := 0; i < len(o.Products); i++ {
				v := &o.Products[i]

				// Fetch the product in the order from the collection
				p, err := products.ReadOne(txCtx, filter.Eq("Id", v.Id))
				if err != nil {
					return err
				}

				// Verify that product quantity in the inventory is more than the
				// product quantity in the order
				if p.Quantity < v.Quantity {
					return fmt.Errorf("low stock on product %v", p.Name)
				}

				// Update the inventory to reduce the product quantity based on the
				// quantity in the order
				if _, err = products.Update(txCtx,
					filter.Eq("Id", v.Id),
					fields.Set("Quantity", p.Quantity-v.Quantity)); err != nil {
					return err
				}

				orderTotal += p.Price * float64(v.Quantity)

				// Remember purchase price in the being created order
				v.Price = p.Price
			}

			// Verify that the user has enough balance to be able to support the
			// order purchase
			if orderTotal > u.Balance {
				return fmt.Errorf("low balance. order total %v", orderTotal)
			}

			// Update the user's balance to account for the order purchase
			if _, err = users.Update(txCtx,
				filter.Eq("Id", o.UserId),
				fields.Set("Balance", u.Balance-orderTotal)); err != nil {
				return err
			}

			orders := tigris.GetCollection[Order](db)

			// Create the order
			_, err = orders.Insert(txCtx, &o)

			// Returning no error means that the transaction will be committed.
			return err
		})

		// If there is no error returned then the transaction was successfully
		// committed and the data has been consistently updated in Tigris.
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, o)
	})
}

// Create routes for searching data in a collection
func setupSearchRoute[T interface{}](r *gin.Engine, db *tigris.Database, name string) {
	r.POST(fmt.Sprintf("/%s/search", name), func(c *gin.Context) {
		coll := tigris.GetCollection[T](db)

		var u search.Request
		if err := c.Bind(&u); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		it, err := coll.Search(c, &u)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		r := &search.Result[T]{}
		for it.Next(r) {
			c.JSON(http.StatusOK, r)
		}
		if err := it.Err(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cloud config
	// tigrisCfg := tigris.Config{URL: "api.preview.tigrisdata.cloud:443", ClientID: "your-tigris-app-id", ClientSecret: "your-tigris-app-secret", Project: "shop"}

	// Local config
	tigrisCfg := tigris.Config{URL: "localhost:8081", Project: "shop"}
	db, err := tigris.OpenDatabase(ctx, &tigrisCfg, &User{}, &Product{}, &Order{})

	if err != nil {
		panic(err)
	}

	r := gin.Default()

	setupCRUDRoutes[User](r, db, "users")
	setupCRUDRoutes[Product](r, db, "products")
	setupReadRoute[Order](r, db, "orders")

	setupCreateOrderRoute(r, db)

	if err := r.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
