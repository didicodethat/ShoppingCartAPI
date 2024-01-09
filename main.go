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
		v1Routes.GET("/lists", v1.GetLists)
		v1Routes.PUT("/lists", v1.CreateList)

		listRoutes := v1Routes.Group("/lists/:listId")
		{
			listRoutes.GET("/sum", v1.ListItemsSUM)
			listRoutes.GET("/items", v1.GetItems)
			listRoutes.PUT("/items", v1.CreateItem)
		}
	}

	r.Run()
}
