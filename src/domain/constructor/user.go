package constructor

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
	"crazyfarmbackend/src/pkg"
)

func ConstructUserFromModel(user dao.User) dto.User {
	return dto.User{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ReferralLink: pkg.ConstructReferralLink(user.ID),
		Icon:         user.Icon,
		LanguageCode: user.LanguageCode,
	}
}

func ConstructUserUpgradeFromModel(userUpgrade dao.UserUpgrade) dto.UserUpgrade {
	return dto.UserUpgrade{
		FarmLvl:   userUpgrade.FarmLvl,
		MaxFields: constant.GetLvlMaxFields(userUpgrade.FarmLvl),
	}
}

func ConstructUserFieldFromModel(userField dao.UserField) dto.UserField {
	return dto.UserField{
		FieldID:   userField.FieldID,
		Plant:     userField.Plant,
		PlantTime: userField.CreatedAt.Unix(),
	}
}

func ConstructUserReferralFromModel(user dao.User) dto.UserReferral {
	return dto.UserReferral{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Icon:      user.Icon,
	}
}
