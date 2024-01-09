package v1

import (
	"Business/ShoppingCart/api"
	"Business/ShoppingCart/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Parses the list id and adds a JSON error to the context as a side effect for convinience
func listIdParse(c *gin.Context) (*uint, error) {
	listId, err := strconv.ParseUint(c.Param("listId"), 10, 64)
	if err != nil {
		errorMessage := api.GetApiError(api.ApiErrorWrongParamType)

		c.JSON(http.StatusBadRequest, errorMessage)
		return nil, &errorMessage
	}
	ulistId := uint(listId)
	return &ulistId, nil
}

func ListItemsSUM(c *gin.Context) {
	var result float32
	listId, err := listIdParse(c)
	if err != nil {
		return
	}
	database.DB.Raw(`SELECT SUM(amount * unitary_price) from list_items where shopping_list_id = ?`, listId).Scan(&result)
	c.JSON(200, result)
}

func CreateItem(c *gin.Context) {
	var listItem database.ListItem
	listId, err := listIdParse(c)
	if err != nil {
		return
	}
	c.Bind(&listItem)
	listItem.ShoppingListID = *listId
	database.DB.Create(&listItem)
	c.JSON(http.StatusCreated, listItem)
}

func GetItems(c *gin.Context) {
	type Result struct {
		ID           uint
		Name         string
		Description  string
		Amount       float32
		UnitaryPrice float32
	}
	var listItems []Result
	listId, err := listIdParse(c)
	if err != nil {
		return
	}

	database.DB.Select([]string{"id", "name", "description", "amount", "unitary_price"}).Table("list_items").Where("shopping_list_id = ?", listId).Scan(&listItems)
	c.JSON(http.StatusOK, listItems)
}
