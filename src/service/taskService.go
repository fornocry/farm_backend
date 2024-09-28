package service

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskService interface {
	GetAllTasks(c *gin.Context) ([]dto.Task, error)
	Check(c *gin.Context) (dto.Task, error)
	Claim(c *gin.Context) (dto.Task, error)
}

type TaskServiceImpl struct {
	taskRepository      repository.TaskRepository
	inventoryRepository repository.InventoryRepository
}

func constructTaskByModel(item dao.Task, status constant.TaskCompleteStatus) dto.Task {
	return dto.Task{
		ID:            item.ID,
		Name:          item.Name,
		Icon:          item.Icon,
		Reward:        item.Reward,
		RewardAmount:  item.RewardAmount,
		NeedDoneTimes: item.NeedDoneTimes,
		Type:          item.Type,
		Data:          item.Data,
		Status:        status,
	}
}

func statusToString(status dao.TaskComplete, err error) constant.TaskCompleteStatus {
	var statusStr constant.TaskCompleteStatus
	if err != nil {
		statusStr = constant.TASK_COMPLETE_NULL
	} else {
		statusStr = status.Status
	}
	return statusStr
}

func (u TaskServiceImpl) Check(c *gin.Context) (dto.Task, error) {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		return dto.Task{}, fmt.Errorf("failed to get user from context")
	}
	taskId := c.Query("taskId")
	if taskId == "" {
		pkg.PanicException(constant.WrongBody, "")
	}
	taskIdUUid, err := uuid.Parse(taskId)
	task, err := u.taskRepository.Get(taskIdUUid)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return dto.Task{}, fmt.Errorf("failed to get user from context")
	}

	status, statusErr := u.taskRepository.GetStatus(user.ID, task.ID)
	if statusErr == nil {
		return constructTaskByModel(task, statusToString(status, statusErr)), nil
	}
	checked, err := u.taskRepository.CheckTask(task, user)
	if !checked {
		if err != nil {
			pkg.PanicException(constant.WrongBody, "")
		}
		return constructTaskByModel(task, statusToString(status, statusErr)), nil
	}

	status, statusErr = u.taskRepository.MarkDone(user.ID, taskIdUUid)
	return constructTaskByModel(task, statusToString(status, statusErr)), nil
}

func (u TaskServiceImpl) Claim(c *gin.Context) (dto.Task, error) {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		return dto.Task{}, fmt.Errorf("failed to get user from context")
	}
	taskId := c.Query("taskId")
	if taskId == "" {
		pkg.PanicException(constant.WrongBody, "")
	}
	taskIdUUid, err := uuid.Parse(taskId)
	task, err := u.taskRepository.Get(taskIdUUid)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return dto.Task{}, fmt.Errorf("failed to get user from context")
	}
	status, statusErr := u.taskRepository.GetStatus(user.ID, task.ID)
	if status.Status == constant.TASK_COMPLETE_FINISHED {
		return constructTaskByModel(task, statusToString(status, statusErr)), nil
	}
	checked, err := u.taskRepository.CheckTask(task, user)
	if !checked {
		if err != nil {
			pkg.PanicException(constant.WrongBody, "")
		}
		return constructTaskByModel(task, statusToString(status, statusErr)), nil
	}
	status, statusErr = u.taskRepository.MarkClaimed(user.ID, taskIdUUid)
	err = u.inventoryRepository.IncreaseItemQuantity(user.ID, task.Reward, task.RewardAmount)
	if err != nil {
		return dto.Task{}, fmt.Errorf("failed to give money for task")
	}
	return constructTaskByModel(task, statusToString(status, statusErr)), nil
}

func (u TaskServiceImpl) GetAllTasks(c *gin.Context) ([]dto.Task, error) {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		return []dto.Task{}, fmt.Errorf("failed to get user from context")
	}
	items, err := u.taskRepository.GetAllTasks()
	if err != nil {
		return []dto.Task{}, err
	}
	var dtoItems []dto.Task
	for _, item := range items {
		status, err := u.taskRepository.GetStatus(user.ID, item.ID)
		dtoItems = append(dtoItems, constructTaskByModel(item, statusToString(status, err)))
	}

	return dtoItems, nil
}

func TaskServiceInit(
	taskRepository repository.TaskRepository,
	inventoryRepository repository.InventoryRepository) *TaskServiceImpl {
	return &TaskServiceImpl{
		taskRepository:      taskRepository,
		inventoryRepository: inventoryRepository,
	}
}
