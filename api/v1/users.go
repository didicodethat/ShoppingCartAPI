package v1

import (
	"Business/ShoppingCart/api"
	"Business/ShoppingCart/database"
	"errors"
	"net/http"
	"net/mail"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const SessionUserIDKey = "user_id"

type LoginAttempt struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

type CreateUserAttempt struct {
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password"`
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

	err = result.SetPassword(c.Password)
	return result, err
}

func CreateUser(c *gin.Context) {
	var userToCreate CreateUserAttempt
	c.Bind(&userToCreate)
	user, err := userToCreate.toUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	result := database.DB.Create(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusConflict, api.GetApiError(api.ApiDuplicatedUserEmail))
			return
		}
		c.JSON(http.StatusInternalServerError, api.UndefinedApiError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
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

	session.Set(SessionUserIDKey, user.ID)
	session.Save()

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}
