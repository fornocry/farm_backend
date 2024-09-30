package service

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
	"crazyfarmbackend/src/domain/dtob"
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"strconv"
	"time"
)

type TaskService interface {
	GetAllTasks(c *gin.Context) ([]dto.Task, error)
	Check(c *gin.Context) (dto.Task, error)
	Claim(c *gin.Context) (dto.Task, error)
}

type TaskServiceImpl struct {
	taskRepository      repository.TaskRepository
	inventoryRepository repository.InventoryRepository
	userRepository      repository.UserRepository
	nc                  *nats.Conn
}

// Helper function to extract user from context
func (s *TaskServiceImpl) getUserFromContext(c *gin.Context) (dao.User, error) {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		return dao.User{}, fmt.Errorf("failed to get user from context")
	}
	return user, nil
}

// Helper function to convert status to string
func statusToString(status dao.TaskComplete, err error) constant.TaskCompleteStatus {
	if err != nil {
		return constant.TASK_COMPLETE_NULL
	}
	return status.Status
}

func (s *TaskServiceImpl) checkTask(task dao.Task, user dao.User) (bool, error) {
	switch task.Type {
	case constant.SUBSCRIBE:
		channelId, _ := task.Data["id"].(string)
		return s.checkTaskSubscribe(strconv.Itoa(int(user.TgId)), channelId)
	case constant.FRIENDS:
		return s.checkTaskFriend(user, task.NeedDoneTimes)
	case constant.INVENTORY:
		plant, _ := task.Data["item"].(string)
		return s.checkTaskInventory(user, task.NeedDoneTimes, constant.Plant(plant))
	default:
		return false, nil
	}
}

func (s *TaskServiceImpl) checkTaskSubscribe(userId string, channelId string) (bool, error) {
	requestData := userId + "," + channelId
	msg, err := s.nc.Request("check_subscribe", []byte(requestData), 10*time.Second)
	if err != nil {
		return false, err
	}
	return string(msg.Data) == "1", nil
}

func (s *TaskServiceImpl) checkTaskFriend(user dao.User, friendsRequired int) (bool, error) {
	userReferrals, err := s.userRepository.GetMyReferrals(user.ID)
	if err != nil {
		return false, err
	}
	if len(userReferrals) >= friendsRequired {
		return true, nil
	}
	return false, nil
}
func (s *TaskServiceImpl) checkTaskInventory(user dao.User, itemsRequired int, itemName constant.Plant) (bool, error) {
	quantity, err := s.inventoryRepository.GetItemQuantity(user.ID, itemName)
	if err != nil {
		return false, err
	}
	if quantity >= itemsRequired {
		return true, nil
	}
	return false, nil
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
		return dtob.ConstructTaskByModel(task, statusToString(status, nil)), nil
	}
	checked, err := s.checkTask(task, user)
	if err != nil || !checked {
		return dtob.ConstructTaskByModel(task, statusToString(status, statusErr)), nil
	}

	status, statusErr = s.taskRepository.MarkDone(user.ID, taskIdUUid)
	return dtob.ConstructTaskByModel(task, statusToString(status, statusErr)), nil
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
		return dtob.ConstructTaskByModel(task, statusToString(status, statusErr)), nil
	}

	checked, err := s.checkTask(task, user)
	if err != nil || !checked {
		return dtob.ConstructTaskByModel(task, statusToString(status, statusErr)), nil
	}

	status, statusErr = s.taskRepository.MarkClaimed(user.ID, taskIdUUid)
	if statusErr != nil {
		return dto.Task{}, fmt.Errorf("failed to mark task as claimed: %v", statusErr)
	}

	err = s.inventoryRepository.IncreaseItemQuantity(user.ID, task.Reward, task.RewardAmount)
	if err != nil {
		return dto.Task{}, fmt.Errorf("failed to give reward for task: %v", err)
	}

	return dtob.ConstructTaskByModel(task, statusToString(status, statusErr)), nil
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
		dtoItems = append(dtoItems, dtob.ConstructTaskByModel(item, statusToString(status, err)))
	}

	return dtoItems, nil
}

func TaskServiceInit(
	taskRepository repository.TaskRepository,
	inventoryRepository repository.InventoryRepository,
	userRepository repository.UserRepository,
	nc *nats.Conn) *TaskServiceImpl {
	return &TaskServiceImpl{
		taskRepository:      taskRepository,
		inventoryRepository: inventoryRepository,
		userRepository:      userRepository,
		nc:                  nc,
	}
}
