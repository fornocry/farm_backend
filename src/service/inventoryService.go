package service

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/constructor"
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type InventoryService interface {
	GetAllItems(c *gin.Context) (dto.GetAllItemsResponse, error)
	GetMyFields(c *gin.Context) []dto.UserField
	PlantField(c *gin.Context) dto.UserField
}

type InventoryServiceImpl struct {
	inventoryRepository repository.InventoryRepository
}

func (u *InventoryServiceImpl) GetAllItems(c *gin.Context) (dto.GetAllItemsResponse, error) {
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
		dtoItems = append(dtoItems, constructor.ConstructInventoryItemByModel(item))
	}
	response := dto.GetAllItemsResponse{
		Items: dtoItems,
	}

	return response, nil
}

func (u *InventoryServiceImpl) GetMyFields(c *gin.Context) []dto.UserField {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "")
	}

	userFields, err := u.inventoryRepository.GetMyFields(user.ID)
	if err != nil {
		pkg.PanicException(constant.DataNotFound, "")
	}

	userFieldDTOs := make([]dto.UserField, len(userFields))
	for i, field := range userFields {
		userFieldDTOs[i] = constructor.ConstructUserFieldFromModel(field)
	}

	return userFieldDTOs
}

func (u *InventoryServiceImpl) PlantField(c *gin.Context) dto.UserField {
	user, ok := c.MustGet("user").(dao.User)
	fieldIDStr := c.Query("fieldID")
	fieldID, err := strconv.Atoi(fieldIDStr)
	if err != nil {
		pkg.PanicException(constant.DataNotFound, "Invalid field id")
	}
	plantStr := c.Query("plant")
	plant := constant.Plant(plantStr)
	if !constant.IsValidPlant(plant) {
		pkg.PanicException(constant.DataNotFound, "Plant not found")
	}
	if !ok {
		pkg.PanicException(constant.DataNotFound, "User not found")
	}

	userField, err := u.inventoryRepository.GetMyField(user.ID, fieldID)
	if err != nil {
		log.Errorln(err)
		pkg.PanicException(constant.DataNotFound, "Error to access field")
	}
	if userField != nil {
		pkg.PanicException(constant.InvalidRequest, "Already planted")
	}

	plantQuantity, err := u.inventoryRepository.GetItemQuantity(user.ID, plant)

	if err != nil {
		log.Errorln(err)
		pkg.PanicException(constant.InvalidRequest, "Access inventory error")
	}
	if !(plantQuantity > 0) {
		pkg.PanicException(constant.InvalidRequest, "Not enough item to plant")
	}

	err = u.inventoryRepository.AdjustItemQuantity(user.ID, plant, -1)
	if err != nil {
		log.Errorln(err)
		pkg.PanicException(constant.InvalidRequest, "Cant decrease")
	}

	userFieldUpdated, err := u.inventoryRepository.PlantField(user.ID, fieldID, plant)
	if err != nil {
		log.Errorln(err)
		pkg.PanicException(constant.DataNotFound, "Failed to plant field")
	}

	return constructor.ConstructUserFieldFromModel(userFieldUpdated)
}

func InventoryServiceInit(inventoryRepository repository.InventoryRepository) *InventoryServiceImpl {
	return &InventoryServiceImpl{
		inventoryRepository: inventoryRepository,
	}
}
