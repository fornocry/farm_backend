package dao

import (
	"crazyfarmbackend/src/constant"
	"github.com/google/uuid"
)

type UserAuth struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"not null"`
	User       User      `gorm:"foreignKey:UserID;column:user_id;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AuthData   string    `gorm:"type:text;not null"`
	AuthMethod string    `gorm:"type:text;not null"`
	BaseModel
}

type User struct {
	ID           uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	TgId         int64     `gorm:"not null;default:0"`
	FirstName    *string   `gorm:"type:text;default:null"`
	LastName     *string   `gorm:"type:text;default:null"`
	Username     *string   `gorm:"type:text;default:null"`
	Icon         *string   `gorm:"type:text;default:null"`
	LanguageCode *string   `gorm:"type:text;default:null"`
	BaseModel
}

type UserUpgrade struct {
	ID      uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	UserID  uuid.UUID `gorm:"not null"`
	User    User      `gorm:"foreignKey:UserID;column:user_id;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FarmLvl int       `gorm:"not null;default:1"`
	BaseModel
}

type UserField struct {
	ID      uuid.UUID      `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	UserID  uuid.UUID      `gorm:"not null"`
	User    User           `gorm:"foreignKey:UserID;column:user_id;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FieldID int            `gorm:"not null;"`
	Plant   constant.Plant `gorm:"not null;"`
	BaseModel
}

type UserReferral struct {
	ID         uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	ReferrerID uuid.UUID `gorm:"not null"`
	Referrer   User      `gorm:"foreignKey:ReferrerID;column:referrer_id;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReferralId uuid.UUID `gorm:"not null;uniqueIndex"`
	Referral   User      `gorm:"foreignKey:ReferralId;column:referral_id;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
