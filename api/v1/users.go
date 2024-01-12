package v1

import (
	"Business/ShoppingCart/api"
	"Business/ShoppingCart/database"
	"net/http"
	"net/mail"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginAttempt struct {
	Email    string
	Password string
}

type CreateUserAttempt struct {
	Email           string
	Password        string
	ConfirmPassword string
}

func (c *CreateUserAttempt) toUser() (database.User, error) {
	result := database.User{}
	_, err := mail.ParseAddress(c.Email)
	if err != nil {
		return result, err
	}
	result.Email = c.Email
	if c.Password != c.ConfirmPassword {
		return result, api.GetApiError(api.ApiPasswordConfirmationMismatch)
	}

	result.SetPassword(c.Password)
	return result, nil
}

func CreateUser(c *gin.Context) {
	var userToCreate CreateUserAttempt
	c.Bind(&userToCreate)
	user, err := userToCreate.toUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	database.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func LogIn(c *gin.Context) {
	session := sessions.Default(c)
	var loginAttempt LoginAttempt
	var user database.User
	c.Bind(&loginAttempt)
	database.DB.Model(&database.User{}).Where("email = ?", loginAttempt.Email).First(&user)

	if user.ComparePassword(loginAttempt.Password) != nil {
		c.JSON(http.StatusUnauthorized, api.GetApiError(api.ApiWrongCredentials))
		return
	}

	session.Set("user_id", user.ID)

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}
