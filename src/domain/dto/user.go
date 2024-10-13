package dto

import (
	"crazyfarmbackend/src/constant"
	"github.com/google/uuid"
)

type AuthRequest struct {
	Method string `json:"method"`
	Data   string `json:"data"`
}

type User struct {
	ID           uuid.UUID `json:"ID"`
	FirstName    *string   `json:"FirstName"`
	LastName     *string   `json:"LastName"`
	ReferralLink string    `json:"ReferralLink"`
	Icon         *string   `json:"Icon"`
	LanguageCode *string   `json:"LanguageCode"`
}

type UserUpgrade struct {
	FarmLvl   int `json:"FarmLvl"`
	MaxFields int `json:"MaxFields"`
}

type UserField struct {
	FieldID   int            `json:"FieldID"`
	Plant     constant.Plant `json:"Plant"`
	PlantTime int64          `json:"PlantTime"`
}

type UserReferral struct {
	ID        uuid.UUID `json:"ID"`
	FirstName *string   `json:"FirstName"`
	LastName  *string   `json:"LastName"`
	Username  *string   `json:"Username"`
	Icon      *string   `json:"Icon"`
}

type UserAuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
