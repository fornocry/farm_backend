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
	CheckTask(
		task dao.Task,
		user dao.User) (bool, error)
	CheckTaskSubscribe(userId string, channelId string) (bool, error)
}

type TaskRepositoryImpl struct {
	db *gorm.DB
	nc *nats.Conn
}

func (u TaskRepositoryImpl) CheckTask(
	task dao.Task,
	user dao.User) (bool, error) {
	switch os := task.Type; os {
	case constant.SUBSCRIBE:
		channelId, _ := task.Data["id"].(string)
		return u.CheckTaskSubscribe(strconv.Itoa(int(user.TgId)), channelId)
	case constant.FRIENDS:
		return false, nil
	}
	return false, nil
}

func (u TaskRepositoryImpl) CheckTaskSubscribe(userId string, channelId string) (bool, error) {
	requestData := userId + "," + channelId
	msg, err := u.nc.Request("check_subscribe", []byte(requestData), 10*time.Second)
	if err != nil {
		log.Error("Got an error when checking subscription. Error: ", err)
		return false, err
	}
	response := string(msg.Data)
	subscribed := response == "1"
	return subscribed, nil
}

func (u TaskRepositoryImpl) Save(task *dao.TaskComplete) (dao.TaskComplete, error) {
	var err = u.db.Save(task).Error
	if err != nil {
		log.Error("Got an error when save user. Error: ", err)
		return dao.TaskComplete{}, err
	}
	return *task, nil
}

func (u TaskRepositoryImpl) Get(
	taskId uuid.UUID) (dao.Task, error) {

	task := dao.Task{
		ID: taskId,
	}
	var err = u.db.Where(&task).First(&task).Error
	if err != nil {
		log.Error("Got an error when get user. Error: ", err)
		return dao.Task{}, err
	}
	return task, nil
}

func (u TaskRepositoryImpl) GetAllTasks() ([]dao.Task, error) {
	var tasks []dao.Task
	if err := u.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (u TaskRepositoryImpl) GetStatus(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error) {
	taskComplete := dao.TaskComplete{
		UserID: userId,
		TaskID: taskId,
	}

	var err = u.db.Where(&taskComplete).First(&taskComplete).Error
	if err != nil {
		return dao.TaskComplete{}, err
	}
	return taskComplete, nil
}

func (u TaskRepositoryImpl) MarkDone(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error) {
	taskComplete := dao.TaskComplete{
		UserID: userId,
		TaskID: taskId,
		Status: constant.TASK_COMPLETE_DONE,
	}

	taskComplete, err := u.Save(&taskComplete)
	if err != nil {
		log.Error("Got an error when creating taskComplete as check. Error: ", err)
		return dao.TaskComplete{}, err
	}
	return taskComplete, nil
}

func (u TaskRepositoryImpl) MarkClaimed(userId uuid.UUID, taskId uuid.UUID) (dao.TaskComplete, error) {
	var taskComplete dao.TaskComplete
	if err := u.db.Where("user_id = ? AND task_id = ?", userId, taskId).First(&taskComplete).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			taskComplete = dao.TaskComplete{
				UserID: userId,
				TaskID: taskId,
				Status: constant.TASK_COMPLETE_FINISHED,
			}
			if err := u.db.Create(&taskComplete).Error; err != nil {
				return dao.TaskComplete{}, err
			}
		} else {
			return dao.TaskComplete{}, err
		}
	} else {
		taskComplete.Status = constant.TASK_COMPLETE_FINISHED
		if err := u.db.Save(&taskComplete).Error; err != nil {
			return dao.TaskComplete{}, err
		}
	}
	return taskComplete, nil

}

func TaskRepositoryInit(db *gorm.DB,
	nc *nats.Conn) *TaskRepositoryImpl {
	_ = db.AutoMigrate(&dao.Task{}, &dao.TaskComplete{})
	return &TaskRepositoryImpl{
		db: db,
		nc: nc,
	}
}
