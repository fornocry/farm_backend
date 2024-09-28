package dao

import (
	"time"
)

type BaseModel struct {
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at" json:"-"`
}
