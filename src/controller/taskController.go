package controller

import (
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskController interface {
	GetAllTasks(c *gin.Context)
	Check(c *gin.Context)
	Claim(c *gin.Context)
}

type TaskControllerImpl struct {
	taskService service.TaskService
}

func (u TaskControllerImpl) Check(c *gin.Context) {
	defer pkg.PanicHandler(c)
	tasks, err := u.taskService.Check(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
	return
}
func (u TaskControllerImpl) Claim(c *gin.Context) {
	defer pkg.PanicHandler(c)
	tasks, err := u.taskService.Claim(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
	return
}
func (u TaskControllerImpl) GetAllTasks(c *gin.Context) {
	defer pkg.PanicHandler(c)
	tasks, err := u.taskService.GetAllTasks(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, tasks)
	return
}

func TaskControllerInit(taskService service.TaskService) *TaskControllerImpl {
	return &TaskControllerImpl{
		taskService: taskService,
	}
}
