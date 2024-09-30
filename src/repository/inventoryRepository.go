package repository

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetAllInventoryItems(userId uuid.UUID) ([]dao.InventoryItem, error)
	IncreaseItemQuantity(userId uuid.UUID, plant constant.Plant, amount int) error
	GetItemQuantity(userId uuid.UUID, plant constant.Plant) (int, error)
}

type InventoryRepositoryImpl struct {
	db *gorm.DB
}

func (u InventoryRepositoryImpl) GetAllInventoryItems(userId uuid.UUID) ([]dao.InventoryItem, error) {
	var inventoryItems []dao.InventoryItem
	err := u.db.Where("user_id = ? AND plant IN ?", userId, constant.Plants).Find(&inventoryItems).Error
	if err != nil {
		return nil, err
	}
	existingPlants := make(map[constant.Plant]bool)
	for _, item := range inventoryItems {
		existingPlants[item.Plant] = true
	}
	var newItems []dao.InventoryItem
	for _, plant := range constant.Plants {
		if !existingPlants[plant] {
			newItems = append(newItems, dao.InventoryItem{
				UserID: userId,
				Plant:  plant,
			})
		}
	}
	if len(newItems) > 0 {
		if err := u.db.Create(&newItems).Error; err != nil {
			return nil, err
		}
		inventoryItems = append(inventoryItems, newItems...)
	}

	orderedInventoryItems := make([]dao.InventoryItem, 0, len(constant.Plants))
	for _, plant := range constant.Plants {
		for _, item := range inventoryItems {
			if item.Plant == plant {
				orderedInventoryItems = append(orderedInventoryItems, item)
				break // Move to the next plant once found
			}
		}
	}

	return orderedInventoryItems, nil
}

func (u InventoryRepositoryImpl) IncreaseItemQuantity(userId uuid.UUID, plant constant.Plant, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	var item dao.InventoryItem
	err := u.db.Where("user_id = ? AND plant = ?", userId, plant).First(&item).Error
	if err != nil {
		return err // Item not found or other error
	}
	item.Quantity += amount // Assuming you have a Quantity field
	if err := u.db.Save(&item).Error; err != nil {
		return err
	}

	return nil
}
func (u InventoryRepositoryImpl) GetItemQuantity(userId uuid.UUID, plant constant.Plant) (int, error) {
	var item dao.InventoryItem
	err := u.db.Where("user_id = ? AND plant = ?", userId, plant).First(&item).Error
	if err != nil {
		return 0, err // Item not found or other error
	}
	return item.Quantity, nil
}

func InventoryRepositoryInit(db *gorm.DB) *InventoryRepositoryImpl {
	_ = db.AutoMigrate(&dao.InventoryItem{})
	return &InventoryRepositoryImpl{
		db: db,
	}
}
