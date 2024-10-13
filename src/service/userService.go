package service

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/constructor"
	"crazyfarmbackend/src/domain/dao"
	"crazyfarmbackend/src/domain/dto"
	"crazyfarmbackend/src/pkg"
	"crazyfarmbackend/src/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

type UserService interface {
	AuthUser(*gin.Context) (dto.UserAuthResponse, error)
	GetMe(*gin.Context) dto.User
	GetUserUpgrade(c *gin.Context) dto.UserUpgrade
	GetMyReferrals(c *gin.Context) []dto.UserReferral
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func (u *UserServiceImpl) logAndReturnError(context string, err error) {
	log.Error(context, err)
	pkg.PanicException(constant.UnknownError, "")
}

func (u *UserServiceImpl) AuthUser(c *gin.Context) (dto.UserAuthResponse, error) {
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
		u.logAndReturnError("Creating JWT token failed: ", err)
	}

	return dto.UserAuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *UserServiceImpl) AuthUserTelegram(data string) (dto.User, dao.UserAuth) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramInitDataTtl, err := strconv.ParseFloat(os.Getenv("TELEGRAM_INIT_DATA_TTL"), 64)
	if err != nil {
		u.logAndReturnError("Parsing TELEGRAM_INIT_DATA_TTL failed: ", err)
	}

	telegramInitData, err := pkg.ParseTelegramData(data)
	if err != nil {
		u.logAndReturnError("Initdata parse: ", err)
	}

	if err := pkg.ValidateTelegramData(data, telegramToken, time.Duration(telegramInitDataTtl*float64(time.Second))); err != nil {
		u.logAndReturnError("Validating initdata: ", err)
	}

	telegramUserIDStr := strconv.FormatInt(telegramInitData.TelegramUser.ID, 10)
	userAuth, isFirst, err := u.userRepository.GetOrCreateAuth(telegramUserIDStr, constant.Telegram)
	if err != nil {
		u.logAndReturnError("GetByAuth from database error: ", err)
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
		u.logAndReturnError("Updating data: ", err)
	}

	if isFirst {
		u.handleReferral(telegramInitData.StartParam, user.ID)
	}
	return constructor.ConstructUserFromModel(user), userAuth
}

func (u *UserServiceImpl) handleReferral(startParam string, userID uuid.UUID) {
	decodedParam := pkg.DecodeStartParam(startParam)
	if decodedParam.Method == "ref" {
		referrerId, err := uuid.Parse(decodedParam.Data)
		if err == nil {
			_, _ = u.userRepository.SetReferrals(userID, referrerId)
		}
	}
}

func (u *UserServiceImpl) GetMe(c *gin.Context) dto.User {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "User not found")
	}
	return constructor.ConstructUserFromModel(user)
}

func (u *UserServiceImpl) GetUserUpgrade(c *gin.Context) dto.UserUpgrade {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "")
	}

	userUpgrade, err := u.userRepository.GetUserUpgrade(user.ID)
	if err != nil {
		pkg.PanicException(constant.DataNotFound, "")
	}
	return constructor.ConstructUserUpgradeFromModel(userUpgrade)
}

func (u *UserServiceImpl) GetMyReferrals(c *gin.Context) []dto.UserReferral {
	user, ok := c.MustGet("user").(dao.User)
	if !ok {
		pkg.PanicException(constant.DataNotFound, "")
	}

	userReferrals, err := u.userRepository.GetMyReferrals(user.ID)
	if err != nil {
		pkg.PanicException(constant.DataNotFound, "")
	}

	userReferralsDTOs := make([]dto.UserReferral, len(userReferrals))
	for i, referral := range userReferrals {
		userReferralsDTOs[i] = constructor.ConstructUserReferralFromModel(referral)
	}

	return userReferralsDTOs
}

func UserServiceInit(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}
