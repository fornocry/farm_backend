package service

import (
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
	"crazyfarmbackend/src/domain/dtob"
	"crazyfarmbackend/src/repository"
	"fmt"
	"github.com/gin-gonic/gin"
)

type InventoryService interface {
	GetAllItems(c *gin.Context) (dto.GetAllItemsResponse, error)
}

type InventoryServiceImpl struct {
	inventoryRepository repository.InventoryRepository
}

func (u InventoryServiceImpl) GetAllItems(c *gin.Context) (dto.GetAllItemsResponse, error) {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		return dto.GetAllItemsResponse{}, fmt.Errorf("failed to get user from context")
	}
	items, err := u.inventoryRepository.GetAllInventoryItems(user.ID)
	if err != nil {
		return dto.GetAllItemsResponse{}, err
	}
	var dtoItems []dto.InventoryItem
	for _, item := range items {
		dtoItems = append(dtoItems, dtob.ConstructInventoryItemByModel(item))
	}
	response := dto.GetAllItemsResponse{
		Items: dtoItems,
	}

	return response, nil
}

func InventoryServiceInit(inventoryRepository repository.InventoryRepository) *InventoryServiceImpl {
	return &InventoryServiceImpl{
		inventoryRepository: inventoryRepository,
	}
}
