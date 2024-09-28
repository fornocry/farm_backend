package dto

import (
	"crazyfarmbackend/src/constant"
	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID
	Name          string
	Icon          *string
	Reward        constant.Plant
	RewardAmount  int
	NeedDoneTimes int
	Type          constant.Task
	Data          map[string]interface{}
	Status        constant.TaskCompleteStatus
}
