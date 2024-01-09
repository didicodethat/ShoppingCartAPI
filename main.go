package main

import (
	v1 "Business/ShoppingCart/api/v1"
	"Business/ShoppingCart/database"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := database.SetupDB()
	if err != nil {
		panic(fmt.Sprintf("Couln't Setup database, got error: %s", err))
	}

	r := gin.Default()
	v1Routes := r.Group("/v1")
	{
		v1Routes.GET("/lists/:listId/sum", v1.ListItemsSUM)

		v1Routes.GET("/lists", func(c *gin.Context) {

		})

		v1Routes.GET("/lists/:listId/items", v1.GetItems)

		v1Routes.PUT("/lists/:listId/items", v1.CreateItem)
	}

	r.Run()
}
