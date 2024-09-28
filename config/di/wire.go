// go:build wireinject
//go:build wireinject
// +build wireinject

package di

import (
	"crazyfarmbackend/config"
	"crazyfarmbackend/src/api/middlewares"
	"crazyfarmbackend/src/controller"
	"crazyfarmbackend/src/repository"
	"crazyfarmbackend/src/service"
	"github.com/google/wire"
)

var connectionsSet = wire.NewSet(
	config.ConnectToDB,
	config.ConnectToNatsBroker,
)

var middlewareServiceSet = wire.NewSet(middlewares.MiddlewareServiceInit,
	wire.Bind(new(middlewares.MiddlewareService), new(*middlewares.MiddlewareServiceImpl)))

var natsBrokerSet = wire.NewSet(
	config.NatsBrokerInit,
	wire.Bind(new(config.NatsBroker), new(*config.NatsBrokerImpl)),
)

var userSet = wire.NewSet(
	repository.UserRepositoryInit,
	wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
	service.UserServiceInit,
	wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
	controller.UserControllerInit,
	wire.Bind(new(controller.UserController), new(*controller.UserControllerImpl)),
)

var inventorySet = wire.NewSet(
	repository.InventoryRepositoryInit,
	wire.Bind(new(repository.InventoryRepository), new(*repository.InventoryRepositoryImpl)),
	service.InventoryServiceInit,
	wire.Bind(new(service.InventoryService), new(*service.InventoryServiceImpl)),
	controller.InventoryControllerInit,
	wire.Bind(new(controller.InventoryController), new(*controller.InventoryControllerImpl)),
)

var taskSet = wire.NewSet(
	repository.TaskRepositoryInit,
	wire.Bind(new(repository.TaskRepository), new(*repository.TaskRepositoryImpl)),
	service.TaskServiceInit,
	wire.Bind(new(service.TaskService), new(*service.TaskServiceImpl)),
	controller.TaskControllerInit,
	wire.Bind(new(controller.TaskController), new(*controller.TaskControllerImpl)),
)

func Init() *Initialization {
	wire.Build(NewInitialization,
		natsBrokerSet,
		connectionsSet,
		userSet,
		inventorySet,
		taskSet,
		middlewareServiceSet)
	return nil
}
