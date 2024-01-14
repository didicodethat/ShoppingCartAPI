package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ListItem struct {
	gorm.Model
	Name         string       `form:"name" json:"name"`
	Description  string       `form:"description" json:"description"`
	Amount       float32      `form:"amount" json:"amount"`
	UnitaryPrice float32      `form:"unitaryPrice" json:"unitaryPrice"`
	Order        int32        `form:"order" json:"order"`
	OwnerListID  uint         `json:"ownerListId"`
	OwnerList    ShoppingList `json:"shoppingList" gorm:"foreignKey:OwnerListID"`
}

func (item *ListItem) UserHasAccess(userID uint) bool {
	return item.OwnerList.UserHasAccess(userID)
}

func GetListItem(itemID uint) ListItem {
	var item ListItem
	item.ID = itemID
	DB.Preload("OwnerList.Members").Preload(clause.Associations).Find(&item)
	return item
}
