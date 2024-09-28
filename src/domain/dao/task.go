package dao

import (
	"crazyfarmbackend/src/constant"
	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID              `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Name          string                 `gorm:"type:text;default:null"`
	Icon          *string                `gorm:"type:text;default:null"`
	Reward        constant.Plant         `gorm:"not null"`
	RewardAmount  int                    `gorm:"type:int;default:0"`
	NeedDoneTimes int                    `gorm:"type:int;default:0"`
	Type          constant.Task          `gorm:"type:text;default:null"`
	Data          map[string]interface{} `gorm:"serializer:json"`
	BaseModel
}

type TaskComplete struct {
	ID     uuid.UUID                   `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID                   `gorm:"not null"`
	User   User                        `gorm:"foreignKey:UserID;column:user_id;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TaskID uuid.UUID                   `gorm:"not null"`
	Task   Task                        `gorm:"not null;default:1"`
	Status constant.TaskCompleteStatus `gorm:"not null;"`
	BaseModel
}
