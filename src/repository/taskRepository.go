package repository

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"errors"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type TaskRepository interface {
	Save(task *dao.TaskComplete) (dao.TaskComplete, error)
	Get(taskId uuid.UUID) (dao.Task, error)
	GetAllTasks() ([]dao.Task, error)
	GetStatus(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error)
	MarkDone(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error)
	MarkClaimed(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error)
	CheckTask(task dao.Task, user dao.User) (bool, error)
	CheckTaskSubscribe(userId string, channelId string) (bool, error)
}

type TaskRepositoryImpl struct {
	db *gorm.DB
	nc *nats.Conn
}

func (r *TaskRepositoryImpl) logError(message string, err error) {
	log.Error(message, err)
}

func (r *TaskRepositoryImpl) CheckTask(task dao.Task, user dao.User) (bool, error) {
	switch task.Type {
	case constant.SUBSCRIBE:
		channelId, _ := task.Data["id"].(string)
		return r.CheckTaskSubscribe(strconv.Itoa(int(user.TgId)), channelId)
	case constant.FRIENDS:
		return false, nil
	default:
		return false, nil
	}
}

func (r *TaskRepositoryImpl) CheckTaskSubscribe(userId string, channelId string) (bool, error) {
	requestData := userId + "," + channelId
	msg, err := r.nc.Request("check_subscribe", []byte(requestData), 10*time.Second)
	if err != nil {
		r.logError("Error checking subscription: ", err)
		return false, err
	}
	return string(msg.Data) == "1", nil
}

func (r *TaskRepositoryImpl) Save(task *dao.TaskComplete) (dao.TaskComplete, error) {
	if err := r.db.Save(task).Error; err != nil {
		r.logError("Error saving task: ", err)
		return dao.TaskComplete{}, err
	}
	return *task, nil
}

func (r *TaskRepositoryImpl) Get(taskId uuid.UUID) (dao.Task, error) {
	var task dao.Task
	if err := r.db.First(&task, taskId).Error; err != nil {
		r.logError("Error retrieving task: ", err)
		return dao.Task{}, err
	}
	return task, nil
}

func (r *TaskRepositoryImpl) GetAllTasks() ([]dao.Task, error) {
	var tasks []dao.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		r.logError("Error retrieving all tasks: ", err)
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepositoryImpl) GetStatus(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error) {
	var taskComplete dao.TaskComplete
	if err := r.db.Where("user_id = ? AND task_id = ?", userId, taskId).First(&taskComplete).Error; err != nil {
		r.logError("Error retrieving task status: ", err)
		return dao.TaskComplete{}, err
	}
	return taskComplete, nil
}

func (r *TaskRepositoryImpl) MarkDone(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error) {
	taskComplete := &dao.TaskComplete{
		UserID: userId,
		TaskID: taskId,
		Status: constant.TASK_COMPLETE_DONE,
	}
	return r.Save(taskComplete)
}

func (r *TaskRepositoryImpl) MarkClaimed(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error) {
	var taskComplete dao.TaskComplete
	err := r.db.Where("user_id = ? AND task_id = ?", userId, taskId).First(&taskComplete).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logError("Error retrieving claimed task: ", err)
		return dao.TaskComplete{}, err
	}

	// If task not found, create a new one
	if errors.Is(err, gorm.ErrRecordNotFound) {
		taskComplete = dao.TaskComplete{
			UserID: userId,
			TaskID: taskId,
			Status: constant.TASK_COMPLETE_FINISHED,
		}
		if err := r.db.Create(&taskComplete).Error; err != nil {
			r.logError("Error creating claimed task: ", err)
			return dao.TaskComplete{}, err
		}
	} else {
		// If task already exists, update its status
		taskComplete.Status = constant.TASK_COMPLETE_FINISHED
		if err := r.db.Save(&taskComplete).Error; err != nil {
			r.logError("Error updating claimed task: ", err)
			return dao.TaskComplete{}, err
		}
	}
	return taskComplete, nil
}

func TaskRepositoryInit(db *gorm.DB, nc *nats.Conn) *TaskRepositoryImpl {
	if err := db.AutoMigrate(&dao.Task{}, &dao.TaskComplete{}); err != nil {
		log.Error("Error during AutoMigrate: ", err)
	}
	return &TaskRepositoryImpl{
		db: db,
		nc: nc,
	}
}
