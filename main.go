package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	tigris "github.com/tigrisdata/tigris-client-go/client"
	"github.com/tigrisdata/tigris-client-go/config"
	"github.com/tigrisdata/tigris-client-go/filter"
	"github.com/tigrisdata/tigris-client-go/update"
)

type User struct {
	Id      int `json:"id" binding:"required" tigris:"primary_key"`
	Name    string
	Balance float64
}

type Order struct {
	Id     int `json:"id" binding:"required" tigris:"primary_key"`
	UserId int

	Products []Product
}

type Product struct {
	Id       int `json:"id" binding:"required" tigris:"primary_key"`
	Name     string
	Quantity int
	Price    float64
}

func setupReadRoute[T interface{}](r *gin.Engine, db *tigris.Database, name string) {
	r.GET(fmt.Sprintf("/%s/read/:id", name), func(c *gin.Context) {
		coll := tigris.GetCollection[T](db)

		id, _ := strconv.Atoi(c.Param("id"))

		u, err := coll.ReadOne(c, filter.Eq("id", id))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, u)
	})
}

func setupCRUDRoutes[T interface{}](r *gin.Engine, db *tigris.Database, name string) {
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

	setupReadRoute[T](r, db, name)

	r.DELETE(fmt.Sprintf("/%s/delete/:id", name), func(c *gin.Context) {
		coll := tigris.GetCollection[T](db)

		id, _ := strconv.Atoi(c.Param("id"))

		if _, err := coll.Delete(c, filter.Eq("id", id)); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Status": "DELETED"})
	})
}

// Create an order transactionally,
// taking into account user balances and product stock.
func setupCreateOrderRoute(r *gin.Engine, db *tigris.Database) {
	r.POST("/orders/create", func(c *gin.Context) {
		var o Order
		if err := c.Bind(&o); err != nil {
			return
		}

		err := db.Tx(c, func(ctx context.Context, tx *tigris.Tx) error {
			users := tigris.GetTxCollection[User](tx)

			u, err := users.ReadOne(ctx, filter.Eq("id", o.UserId))
			if err != nil {
				return err
			}

			products := tigris.GetTxCollection[Product](tx)

			orderTotal := 0.0

			for _, v := range o.Products {
				p, err := products.ReadOne(ctx, filter.Eq("id", v.Id))
				if err != nil {
					return err
				}

				if p.Quantity < v.Quantity {
					return fmt.Errorf("low stock on product %v", p.Name)
				}

				_, err = products.Update(ctx, filter.Eq("id", v.Id),
					update.Set("Quantity", p.Quantity-v.Quantity))
				if err != nil {
					return err
				}

				orderTotal += p.Price * float64(v.Quantity)
			}

			if orderTotal > u.Balance {
				return fmt.Errorf("low balance. order total %v", orderTotal)
			}

			_, err = users.Update(ctx, filter.Eq("id", o.UserId),
				update.Set("Balance", u.Balance-orderTotal))
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Status": "PURCHASED"})
	})
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := tigris.OpenDatabase(ctx, &tigris.DatabaseConfig{Config: config.Config{URL: "localhost:8081"}},
		"shop", &User{}, &Product{}, &Order{})
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