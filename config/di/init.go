package di

import (
	"crazyfarmbackend/config"
	"crazyfarmbackend/src/api/middlewares"
	"crazyfarmbackend/src/controller"
	"crazyfarmbackend/src/repository"
	"crazyfarmbackend/src/service"
)

type Initialization struct {
	UserRepository repository.UserRepository
	UserService    service.UserService
	UserController controller.UserController

	InventoryRepository repository.InventoryRepository
	InventoryService    service.InventoryService
	InventoryController controller.InventoryController

	TaskRepository repository.TaskRepository
	TaskService    service.TaskService
	TaskController controller.TaskController

	MiddlewareService middlewares.MiddlewareService
	Nats              config.NatsBroker
}

func NewInitialization(
	userRepository repository.UserRepository,
	userService service.UserService,
	userController controller.UserController,

	inventoryRepository repository.InventoryRepository,
	inventoryService service.InventoryService,
	inventoryController controller.InventoryController,

	taskRepository repository.TaskRepository,
	taskService service.TaskService,
	taskController controller.TaskController,

	middlewareService middlewares.MiddlewareService,
	nats config.NatsBroker) *Initialization {
	return &Initialization{
		UserRepository:      userRepository,
		UserService:         userService,
		UserController:      userController,
		InventoryRepository: inventoryRepository,
		InventoryService:    inventoryService,
		InventoryController: inventoryController,
		TaskRepository:      taskRepository,
		TaskService:         taskService,
		TaskController:      taskController,
		MiddlewareService:   middlewareService,
		Nats:                nats,
	}
}
