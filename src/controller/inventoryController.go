package controller

import (
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InventoryController interface {
	GetInventoryItems(c *gin.Context)
	GetMyFields(c *gin.Context)
	PlantField(c *gin.Context)
}

type InventoryControllerImpl struct {
	inventoryService service.InventoryService
}

func (u *InventoryControllerImpl) GetInventoryItems(c *gin.Context) {
	defer pkg.PanicHandler(c)
	inventoryItems, err := u.inventoryService.GetAllItems(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, inventoryItems)
	return
}

func (u *InventoryControllerImpl) GetMyFields(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userResponse := u.inventoryService.GetMyFields(c)
	c.JSON(http.StatusOK, userResponse)
	return
}

func (u *InventoryControllerImpl) PlantField(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userResponse := u.inventoryService.PlantField(c)
	c.JSON(http.StatusOK, userResponse)
	return
}

func InventoryControllerInit(inventoryService service.InventoryService) *InventoryControllerImpl {
	return &InventoryControllerImpl{
		inventoryService: inventoryService,
	}
}
