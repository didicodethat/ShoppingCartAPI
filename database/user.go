package database

import (
	"Business/ShoppingCart/core"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email              string          `gorm:"uniqueIndex" json:"email"`
	Password           string          `json:"-"`
	ShoppingLists      []*ShoppingList `gorm:"many2many:user_shopping_lists;" json:"shoppingLists"`
	OwnedShoppingLists []ShoppingList  `gorm:"foreignKey:UserID" json:"ownedShoppingLists"`
}

func encryptPassword(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword+core.Secrets.PasswordSecret), 14)

	return string(bytes), err
}

// Compares the saved password witht he password sent
func (u *User) ComparePassword(rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(rawPassword+core.Secrets.PasswordSecret))
}

// Use this, and don't set the password directly, because of encryption
func (u *User) SetPassword(rawPassword string) error {
	encrypted, err := encryptPassword(rawPassword)
	if err != nil {
		return err
	}
	u.Password = encrypted
	return err
}

func (u *User) OwnedLists() []ShoppingList {
	var shoppingLists []ShoppingList
	if err := DB.Model(&u).Association("OwnedShoppingLists").Find(&shoppingLists); err != nil {
		panic(err)
	}
	return shoppingLists
}
