package database

import (
	"Business/ShoppingCart/core"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex"`
	Password      string
	ShoppingLists []*ShoppingList `gorm:"many2many:user_shopping_lists;"`
}

func encryptPassword(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword+core.Secrets.PasswordSecret), 14)

	return string(bytes), err
}

func (u *User) ComparePassword(rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword+core.Secrets.PasswordSecret))
}

func (u *User) SetPassword(rawPassword string) error {
	encrypted, err := encryptPassword(rawPassword)
	if err != nil {
		return err
	}
	u.Password = encrypted
	return err
}
