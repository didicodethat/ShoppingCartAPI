package v1

import (
	"Business/ShoppingCart/api"
	"Business/ShoppingCart/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	ContextShoppingList = "ShopppingList"
	ContextUser         = "User"
)

func userFromContext(c *gin.Context) database.User {
	user, ok := c.Get(ContextUser)

	if !ok {
		panic("Only use userFromContext in authenticaded areas, make sure to set the user parameter on a middleware before this one")
	}

	return user.(database.User)
}

func shoppingListFromContext(c *gin.Context) database.ShoppingList {
	list, ok := c.Get(ContextShoppingList)

	if !ok {
		panic("Only use shoppingListFromContext in authenticaded areas, make sure to set the user parameter on a middleware before this one")
	}

	return list.(database.ShoppingList)
}

// Middleware that loads the current list based on the :listId param of the url
// Only call it after authenticated or if the user variable is set on the Context.
// it will also only load the list if the user has access to the current list
func LoadList() gin.HandlerFunc {
	return func(c *gin.Context) {
		listId, err := strconv.ParseUint(c.Param("listId"), 10, 64)
		user := userFromContext(c)

		if err != nil {
			errorMessage := api.GetApiError(api.ApiErrorWrongParamType)
			c.JSON(http.StatusBadRequest, errorMessage)
			return
		}

		list, err := database.FindShoppingList(uint(listId))

		if err != nil {
			errorMessage := api.GetApiError(api.ApiEntityNotFound)
			c.JSON(http.StatusBadRequest, errorMessage)
			return
		}

		if !list.UserHasAccess(user.ID) {
			errorMessage := api.GetApiError(api.ApiForbidden)
			c.JSON(http.StatusForbidden, errorMessage)
			return
		}

		c.Set(ContextShoppingList, list)

		c.Next()
	}
}

func CreateList(c *gin.Context) {
	var shoppingListRequest database.CreateShoppingListRequest
	c.Bind(&shoppingListRequest)
	shoppingList := database.ShoppingList{
		Name:        shoppingListRequest.Name,
		Description: shoppingListRequest.Description,
	}
	shoppingList.Owner = userFromContext(c)
	if err := database.DB.Create(&shoppingList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, shoppingList)
}

func GetLists(c *gin.Context) {
	user := userFromContext(c)
	c.JSON(http.StatusOK, user.OwnedLists())
}

func CreateItem(c *gin.Context) {
	var listItem database.ListItem
	c.Bind(&listItem)
	list := shoppingListFromContext(c)
	listItem.OwnerList = list

	database.DB.Create(&listItem)
	c.JSON(http.StatusCreated, listItem)
}

func GetList(c *gin.Context) {
	c.JSON(http.StatusOK, shoppingListFromContext(c))
}

// TODO: Create the patch and delete routes.
