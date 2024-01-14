package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShoppingList struct {
	gorm.Model
	Name        string     `form:"name" json:"name"`
	Description *string    `form:"description" json:"description"`
	UserID      uint       `form:"userId" json:"userId"`
	Owner       User       `json:"owner" gorm:"foreignKey:UserID"`
	Members     []*User    `gorm:"many2many:user_shopping_lists;" json:"members"`
	ListItems   []ListItem `json:"listItems" gorm:"foreignKey:OwnerListID"`
}

type CreateShoppingListRequest struct {
	Name        string  `form:"name" json:"name"`
	Description *string `form:"description" json:"description"`
}

func FindShoppingList(id uint) (ShoppingList, error) {
	var list ShoppingList
	list.ID = id
	result := DB.Preload(clause.Associations).Find(&list)
	return list, result.Error
}

// checks if the user has access to the list
func (s *ShoppingList) UserHasAccess(userID uint) bool {
	if s.Owner.ID == userID {
		return true
	}

	for _, member := range s.Members {
		if member.ID == userID {
			return true
		}
	}

	return false
}
