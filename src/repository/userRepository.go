package repository

import (
	"crazyfarmbackend/src/domain/dao"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user *dao.User) (dao.User, error)
	Get(userId uuid.UUID) (dao.User, error)
	GetByAuth(data, method string) (dao.User, error)
	GetByAuthObj(data, method string) (dao.UserAuth, error)
	GetByAuthId(ID uuid.UUID) (dao.UserAuth, error)
	Create(data, method string) (dao.User, error)
	GetOrCreate(data, method string) (dao.User, error)
	GetOrCreateAuth(data, method string) (dao.UserAuth, bool, error)
	UpdateUserFields(userId uuid.UUID, updates map[string]interface{}) (dao.User, error)
	GetUserUpgrade(userId uuid.UUID) (dao.UserUpgrade, error)
	GetMyReferrals(userId uuid.UUID) ([]dao.User, error)
	SetReferrals(userId, referrerId uuid.UUID) (dao.UserReferral, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u *UserRepositoryImpl) logAndReturnError(message string, err error) error {
	log.Error(message, err)
	return err
}

func (u *UserRepositoryImpl) Save(user *dao.User) (dao.User, error) {
	if err := u.db.Save(user).Error; err != nil {
		return dao.User{}, u.logAndReturnError("Error saving user: ", err)
	}
	return *user, nil
}

func (u *UserRepositoryImpl) Get(userId uuid.UUID) (dao.User, error) {
	var user dao.User
	if err := u.db.First(&user, userId).Error; err != nil {
		return dao.User{}, u.logAndReturnError("Error getting user: ", err)
	}
	return user, nil
}

func (u *UserRepositoryImpl) getByAuthCommon(authMethod dao.UserAuth) (dao.UserAuth, error) {
	if err := u.db.Where(&authMethod).Preload("User").First(&authMethod).Error; err != nil {
		return dao.UserAuth{}, u.logAndReturnError("Error getting user by auth: ", err)
	}
	return authMethod, nil
}

func (u *UserRepositoryImpl) GetByAuth(data, method string) (dao.User, error) {
	authMethod := dao.UserAuth{AuthData: data, AuthMethod: method}
	auth, err := u.getByAuthCommon(authMethod)
	if err != nil {
		return dao.User{}, err
	}
	return auth.User, nil
}

func (u *UserRepositoryImpl) GetByAuthObj(data, method string) (dao.UserAuth, error) {
	authMethod := dao.UserAuth{AuthData: data, AuthMethod: method}
	return u.getByAuthCommon(authMethod)
}

func (u *UserRepositoryImpl) GetByAuthId(ID uuid.UUID) (dao.UserAuth, error) {
	authMethod := dao.UserAuth{ID: ID}
	return u.getByAuthCommon(authMethod)
}

func (u *UserRepositoryImpl) Create(data, method string) (dao.User, error) {
	user := dao.User{}
	if _, err := u.Save(&user); err != nil {
		return dao.User{}, err
	}
	authMethod := dao.UserAuth{AuthData: data, AuthMethod: method, User: user}
	if err := u.db.Save(&authMethod).Error; err != nil {
		return dao.User{}, u.logAndReturnError("Error creating user: ", err)
	}
	return user, nil
}

func (u *UserRepositoryImpl) GetOrCreate(data, method string) (dao.User, error) {
	user, err := u.GetByAuth(data, method)
	if err == nil {
		return user, nil
	}
	log.Infof("Creating user with auth method %s and data %s", method, data)
	return u.Create(data, method)
}

func (u *UserRepositoryImpl) GetOrCreateAuth(data, method string) (dao.UserAuth, bool, error) {
	userAuth, err := u.GetByAuthObj(data, method)
	if err == nil {
		return userAuth, false, nil
	}
	log.Infof("Creating user with auth method %s and data %s", method, data)
	user, err := u.Create(data, method)
	if err != nil {
		return dao.UserAuth{}, false, u.logAndReturnError("Error creating user: ", err)
	}
	return dao.UserAuth{User: user}, true, nil
}

func (u *UserRepositoryImpl) UpdateUserFields(userId uuid.UUID, updates map[string]interface{}) (dao.User, error) {
	user := dao.User{ID: userId}
	if err := u.db.Model(&user).Updates(updates).Error; err != nil {
		return dao.User{}, u.logAndReturnError("Error updating user fields: ", err)
	}
	return user, nil
}

func (u *UserRepositoryImpl) GetUserUpgrade(userId uuid.UUID) (dao.UserUpgrade, error) {
	var userUpgrade dao.UserUpgrade
	if err := u.db.Where("user_id = ?", userId).First(&userUpgrade).Error; err == nil {
		return userUpgrade, nil
	}
	userUpgrade = dao.UserUpgrade{UserID: userId, FarmLvl: 1}
	if err := u.db.Create(&userUpgrade).Error; err != nil {
		return dao.UserUpgrade{}, u.logAndReturnError("Error creating user upgrade: ", err)
	}
	return userUpgrade, nil
}

func (u *UserRepositoryImpl) GetMyReferrals(userId uuid.UUID) ([]dao.User, error) {
	var userReferrals []dao.UserReferral
	if err := u.db.Where("referrer_id = ?", userId).Preload("Referral").Find(&userReferrals).Error; err != nil {
		return nil, u.logAndReturnError("Error getting user referrals: ", err)
	}
	var referrals []dao.User
	for _, referral := range userReferrals {
		referrals = append(referrals, referral.Referral)
	}
	return referrals, nil
}

func (u *UserRepositoryImpl) SetReferrals(userId, referrerId uuid.UUID) (dao.UserReferral, error) {
	// Initialize the UserReferral struct with provided userId and referrerId
	userReferral := dao.UserReferral{
		ReferralId: userId,
		ReferrerID: referrerId,
	}

	// Save the userReferral to the database and handle any errors
	if err := u.db.Save(&userReferral).Error; err != nil {
		// Log the error and return it
		return dao.UserReferral{}, u.logAndReturnError("Error creating user referral: ", err)
	}

	// Return the created userReferral
	return userReferral, nil
}
func UserRepositoryInit(db *gorm.DB) *UserRepositoryImpl {
	_ = db.AutoMigrate(&dao.User{}, &dao.UserAuth{}, &dao.UserUpgrade{}, &dao.UserField{}, &dao.UserReferral{})
	return &UserRepositoryImpl{db: db}
}
