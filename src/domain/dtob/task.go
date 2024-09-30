package dtob

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
)

func ConstructTaskByModel(item dao.Task, status constant.TaskCompleteStatus) dto.Task {
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
