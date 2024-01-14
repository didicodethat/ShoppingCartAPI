package main

import (
	"Business/ShoppingCart/api"
	v1 "Business/ShoppingCart/api/v1"
	"Business/ShoppingCart/core"
	"Business/ShoppingCart/database"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
)

//go:embed settings/secrets.toml
var SecretsToml []byte

// user must be authenticated to access this route
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(v1.SessionUserIDKey)
		if userID == nil {
			c.JSON(http.StatusUnauthorized, api.GetApiError(api.ApiRestrictedArea))
			return
		}
		var user database.User
		database.DB.Table("users").Where("id = ?", userID).First(&user)
		c.Set(v1.ContextUser, user)
		c.Next()
	}
}

func main() {
	core.LoadSecrets(SecretsToml)
	db, err := database.SetupDB()
	if err != nil {
		panic(fmt.Sprintf("Couln't Setup database, got error: %s", err))
	}

	store := gormsessions.NewStore(db, true, []byte("secret"))
	r := gin.Default()
	r.Use(sessions.Sessions("main_session", store))
	v1Routes := r.Group("/v1")
	{
		authRoutes := v1Routes.Group("/auth")
		{
			authRoutes.PUT("/register", v1.CreateUser)
			authRoutes.POST("/login", v1.LogIn)
		}

		listRoutes := v1Routes.Group("/lists")
		listRoutes.Use(Authenticated())
		{
			listRoutes.GET("/", v1.GetLists)
			listRoutes.PUT("/", v1.CreateList)
			listItemRoutes := listRoutes.Group("/:listId")
			listItemRoutes.Use(v1.LoadList())
			{
				listItemRoutes.GET("/", v1.GetList)
				listItemRoutes.PUT("/items", v1.CreateItem)
			}
		}
	}

	r.Run()
}
