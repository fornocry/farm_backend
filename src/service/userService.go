package service

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/repository"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type UserService interface {
	AuthUser(*gin.Context) (dto.UserAuthResponse, error)
	GetMe(*gin.Context) dto.User
	GetUserUpgrade(c *gin.Context) dto.UserUpgrade
	GetMyFields(c *gin.Context) []dto.UserField
	GetMyReferrals(c *gin.Context) []dto.UserReferral
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func constructUserByModel(user dao.User) dto.User {
	return dto.User{ID: user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ReferralLink: pkg.ConstructReferralLink(user.ID),
		Icon:         user.Icon,
		LanguageCode: user.LanguageCode,
	}
}
func constructUserUpgradeByModel(userUpgrade dao.UserUpgrade) dto.UserUpgrade {
	return dto.UserUpgrade{
		FarmLvl:   userUpgrade.FarmLvl,
		MaxFields: constant.GetLvlMaxFields(userUpgrade.FarmLvl),
	}
}
func constructUserFieldByModel(userField dao.UserField) dto.UserField {
	return dto.UserField{
		FieldID: userField.FieldID,
		Plant:   userField.Plant,
	}
}

func constructUserReferralByModel(user dao.User) dto.UserReferral {
	return dto.UserReferral{ID: user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Icon:      user.Icon,
	}
}

func (u UserServiceImpl) AuthUser(c *gin.Context) (dto.UserAuthResponse, error) {
	method := c.Query("method")
	data := c.Query("data")
	if method == "" || data == "" {
		pkg.PanicException(constant.WrongBody, "")
	}
	var user dto.User
	var userAuth dao.UserAuth
	switch method {
	case "telegram":
		user, userAuth = u.AuthUserTelegram(data)
	default:
		pkg.PanicException(constant.WrongMethod, "")
	}
	token, err := pkg.CreateJwtToken(userAuth.ID)
	if err != nil {
		pkg.PanicException(constant.UnknownError, "")
	}
	return dto.UserAuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u UserServiceImpl) AuthUserTelegram(data string) (dto.User, dao.UserAuth) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramInitDataTtl, err := strconv.ParseFloat(os.Getenv("TELEGRAM_INIT_DATA_TTL"), 64)
	telegramInitData, err := pkg.ParseTelegramData(data)
	if err != nil {
		log.Error("Initdata parse: ", err)
		pkg.PanicException(constant.InvalidRequest, "")
	}
	if err := pkg.ValidateTelegramData(data, telegramToken, time.Duration(telegramInitDataTtl*float64(time.Second))); err != nil {
		log.Error("Validating initdata: ", err)
		pkg.PanicException(constant.Unauthorized, "")
	}
	telegramUserIDStr := strconv.FormatInt(telegramInitData.TelegramUser.ID, 10)
	userAuth, err := u.userRepository.GetOrCreateAuth(telegramUserIDStr, constant.Telegram)
	if err != nil {
		log.Error("GetByAuth from database error: ", err)
		pkg.PanicException(constant.DataNotFound, "")
	}
	user := userAuth.User
	updates := map[string]interface{}{
		"tg_id":         telegramInitData.TelegramUser.ID,
		"first_name":    pkg.GetNullableString(telegramInitData.TelegramUser.FirstName),
		"last_name":     pkg.GetNullableString(telegramInitData.TelegramUser.LastName),
		"username":      pkg.GetNullableString(telegramInitData.TelegramUser.Username),
		"icon":          pkg.GetNullableString(telegramInitData.TelegramUser.PhotoURL),
		"language_code": pkg.GetNullableString(telegramInitData.TelegramUser.LanguageCode),
	}

	user, err = u.userRepository.UpdateUserFields(user.ID, updates)
	if err != nil {
		log.Error("Updating data: ", err)
		pkg.PanicException(constant.UnknownError, "")
	}
	return constructUserByModel(user), userAuth
}

func (u UserServiceImpl) GetMe(c *gin.Context) dto.User {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "")
	}
	return constructUserByModel(user)
}

func (u UserServiceImpl) GetUserUpgrade(c *gin.Context) dto.UserUpgrade {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "")
	}
	userUpgrade, err := u.userRepository.GetUserUpgrade(user.ID)
	if err != nil {
		pkg.PanicException(constant.DataNotFound, "")
	}
	return constructUserUpgradeByModel(userUpgrade)
}

func (u UserServiceImpl) GetMyFields(c *gin.Context) []dto.UserField {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "")
	}
	userFields, err := u.userRepository.GetMyFields(user.ID)
	if err != nil {
		pkg.PanicException(constant.DataNotFound, "")
	}
	if len(userFields) == 0 {
		return []dto.UserField{}
	}
	var userFieldDTOs []dto.UserField
	for _, field := range userFields {
		userFieldDTOs = append(userFieldDTOs, constructUserFieldByModel(field))
	}

	return userFieldDTOs
}

func (u UserServiceImpl) GetMyReferrals(c *gin.Context) []dto.UserReferral {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "")
	}
	userReferrals, err := u.userRepository.GetMyReferrals(user.ID)
	if err != nil {
		pkg.PanicException(constant.DataNotFound, "")
	}
	if len(userReferrals) == 0 {
		return []dto.UserReferral{}
	}
	var userReferralsDTOs []dto.UserReferral
	for _, referral := range userReferrals {
		userReferralsDTOs = append(userReferralsDTOs, constructUserReferralByModel(referral))
	}

	return userReferralsDTOs
}

func UserServiceInit(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}
