package dao

import (
	"crazyfarmbackend/src/constant"
	"github.com/google/uuid"
)

type InventoryItem struct {
	ID       uuid.UUID      `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	UserID   uuid.UUID      `gorm:"not null"`
	User     User           `gorm:"foreignKey:UserID;column:user_id;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Plant    constant.Plant `gorm:"not null"`
	Quantity int            `gorm:"default:0"`
	BaseModel
}
