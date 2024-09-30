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

// Helper function to extract user from context
func (s *TaskServiceImpl) getUserFromContext(c *gin.Context) (dao.User, error) {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		return dao.User{}, fmt.Errorf("failed to get user from context")
	}
	return user, nil
}

// Helper function to construct a DTO from a model and status
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

// Helper function to convert status to string
func statusToString(status dao.TaskComplete, err error) constant.TaskCompleteStatus {
	if err != nil {
		return constant.TASK_COMPLETE_NULL
	}
	return status.Status
}

func (s *TaskServiceImpl) Check(c *gin.Context) (dto.Task, error) {
	user, err := s.getUserFromContext(c)
	if err != nil {
		return dto.Task{}, err
	}
	taskId := c.Query("taskId")
	if taskId == "" {
		pkg.PanicException(constant.WrongBody, "")
	}

	taskIdUUid, err := uuid.Parse(taskId)
	if err != nil {
		return dto.Task{}, fmt.Errorf("invalid task ID format: %v", err)
	}
	task, err := s.taskRepository.Get(taskIdUUid)
	if err != nil {
		return dto.Task{}, fmt.Errorf("failed to get task: %v", err)
	}
	status, statusErr := s.taskRepository.GetStatus(user.ID, task.ID)
	if statusErr == nil {
		return constructTaskByModel(task, statusToString(status, nil)), nil
	}
	checked, err := s.taskRepository.CheckTask(task, user)
	if err != nil || !checked {
		return constructTaskByModel(task, statusToString(status, statusErr)), nil
	}

	status, statusErr = s.taskRepository.MarkDone(user.ID, taskIdUUid)
	return constructTaskByModel(task, statusToString(status, statusErr)), nil
}

func (s *TaskServiceImpl) Claim(c *gin.Context) (dto.Task, error) {
	user, err := s.getUserFromContext(c)
	if err != nil {
		return dto.Task{}, err
	}

	taskId := c.Query("taskId")
	if taskId == "" {
		pkg.PanicException(constant.WrongBody, "")
	}

	taskIdUUid, err := uuid.Parse(taskId)
	if err != nil {
		return dto.Task{}, fmt.Errorf("invalid task ID format: %v", err)
	}

	task, err := s.taskRepository.Get(taskIdUUid)
	if err != nil {
		return dto.Task{}, fmt.Errorf("failed to get task: %v", err)
	}

	status, statusErr := s.taskRepository.GetStatus(user.ID, task.ID)
	if status.Status == constant.TASK_COMPLETE_FINISHED {
		return constructTaskByModel(task, statusToString(status, statusErr)), nil
	}

	checked, err := s.taskRepository.CheckTask(task, user)
	if err != nil || !checked {
		return constructTaskByModel(task, statusToString(status, statusErr)), nil
	}

	status, statusErr = s.taskRepository.MarkClaimed(user.ID, taskIdUUid)
	if statusErr != nil {
		return dto.Task{}, fmt.Errorf("failed to mark task as claimed: %v", statusErr)
	}

	err = s.inventoryRepository.IncreaseItemQuantity(user.ID, task.Reward, task.RewardAmount)
	if err != nil {
		return dto.Task{}, fmt.Errorf("failed to give reward for task: %v", err)
	}

	return constructTaskByModel(task, statusToString(status, statusErr)), nil
}

func (s *TaskServiceImpl) GetAllTasks(c *gin.Context) ([]dto.Task, error) {
	user, err := s.getUserFromContext(c)
	if err != nil {
		return nil, err
	}

	items, err := s.taskRepository.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var dtoItems []dto.Task
	for _, item := range items {
		status, err := s.taskRepository.GetStatus(user.ID, item.ID)
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
