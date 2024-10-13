package constructor

import (
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
)

func ConstructInventoryItemByModel(item dao.InventoryItem) dto.InventoryItem {
	return dto.InventoryItem{
		Plant:    item.Plant,
		Quantity: item.Quantity}
}
