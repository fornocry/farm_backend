package repository

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetAllInventoryItems(userId uuid.UUID) ([]dao.InventoryItem, error)
	AdjustItemQuantity(userId uuid.UUID, plant constant.Plant, amount int) error
	GetItemQuantity(userId uuid.UUID, plant constant.Plant) (int, error)
	GetMyFields(userId uuid.UUID) ([]dao.UserField, error)
	GetMyField(userId uuid.UUID, fieldID int) (*dao.UserField, error)
	PlantField(userId uuid.UUID, fieldID int, plant constant.Plant) (dao.UserField, error)
}

type InventoryRepositoryImpl struct {
	db *gorm.DB
}

func (u *InventoryRepositoryImpl) GetAllInventoryItems(userId uuid.UUID) ([]dao.InventoryItem, error) {
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

func (u *InventoryRepositoryImpl) AdjustItemQuantity(userId uuid.UUID, plant constant.Plant, amount int) error {
	if amount == 0 {
		return fmt.Errorf("amount must not be zero")
	}

	var item dao.InventoryItem
	err := u.db.Where("user_id = ? AND plant = ?", userId, plant).First(&item).Error
	if err != nil {
		return err // Item not found or other error
	}

	item.Quantity += amount // This will handle both increase and decrease

	if item.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	if err := u.db.Save(&item).Error; err != nil {
		return err
	}

	return nil
}

func (u *InventoryRepositoryImpl) GetItemQuantity(userId uuid.UUID, plant constant.Plant) (int, error) {
	var item dao.InventoryItem
	err := u.db.Where("user_id = ? AND plant = ?", userId, plant).First(&item).Error
	if err != nil {
		return 0, err // Item not found or other error
	}
	return item.Quantity, nil
}

func (u *InventoryRepositoryImpl) GetMyFields(userId uuid.UUID) ([]dao.UserField, error) {
	var userFields []dao.UserField
	if err := u.db.Where("user_id = ?", userId).Find(&userFields).Error; err != nil {
		return nil, err
	}
	return userFields, nil
}
func (u *InventoryRepositoryImpl) GetMyField(userId uuid.UUID, fieldID int) (*dao.UserField, error) {
	var userFields dao.UserField
	if err := u.db.Where("user_id = ? AND field_id = ?", userId, fieldID).First(&userFields).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if no record found
		}
		return nil, err // Return error for other cases
	}
	return &userFields, nil
}

func (u *InventoryRepositoryImpl) PlantField(userId uuid.UUID, fieldID int, plant constant.Plant) (dao.UserField, error) {
	userField := dao.UserField{
		UserID:  userId,
		FieldID: fieldID,
		Plant:   plant,
	}
	if err := u.db.Save(&userField).Error; err != nil {
		return dao.UserField{}, err
	}

	return userField, nil
}

func InventoryRepositoryInit(db *gorm.DB) *InventoryRepositoryImpl {
	_ = db.AutoMigrate(&dao.InventoryItem{})
	return &InventoryRepositoryImpl{
		db: db,
	}
}
