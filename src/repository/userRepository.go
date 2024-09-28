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
	GetByAuth(data string,
		method string) (dao.User, error)
	GetByAuthObj(
		data string,
		method string) (dao.UserAuth, error)
	GetByAuthId(ID uuid.UUID) (dao.UserAuth, error)
	Create(data string,
		method string) (dao.User, error)
	GetOrCreate(data string,
		method string) (dao.User, error)
	GetOrCreateAuth(
		data string,
		method string) (dao.UserAuth, error)
	UpdateUserFields(userId uuid.UUID, updates map[string]interface{}) (dao.User, error)
	GetUserUpgrade(userId uuid.UUID) (dao.UserUpgrade, error)
	GetMyFields(userId uuid.UUID) ([]dao.UserField, error)
	GetMyReferrals(userId uuid.UUID) ([]dao.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u UserRepositoryImpl) Save(user *dao.User) (dao.User, error) {
	var err = u.db.Save(user).Error
	if err != nil {
		log.Error("Got an error when save user. Error: ", err)
		return dao.User{}, err
	}
	return *user, nil
}

func (u UserRepositoryImpl) Get(
	userId uuid.UUID) (dao.User, error) {

	user := dao.User{
		ID: userId,
	}
	var err = u.db.Where(&user).First(&user).Error
	if err != nil {
		log.Error("Got an error when get user. Error: ", err)
		return dao.User{}, err
	}
	return user, nil
}

func (u UserRepositoryImpl) GetByAuth(
	data string,
	method string) (dao.User, error) {

	authMethod := dao.UserAuth{
		AuthData:   data,
		AuthMethod: method,
	}
	var err = u.db.Where(&authMethod).Preload("User").First(&authMethod).Error
	if err != nil {
		log.Error("Got an error when get user. Error: ", err)
		return dao.User{}, err
	}
	user := authMethod.User
	return user, nil
}
func (u UserRepositoryImpl) GetByAuthObj(
	data string,
	method string) (dao.UserAuth, error) {

	authMethod := dao.UserAuth{
		AuthData:   data,
		AuthMethod: method,
	}
	var err = u.db.Where(&authMethod).Preload("User").First(&authMethod).Error
	if err != nil {
		log.Error("Got an error when get user. Error: ", err)
		return dao.UserAuth{}, err
	}
	return authMethod, nil
}

func (u UserRepositoryImpl) GetByAuthId(ID uuid.UUID) (dao.UserAuth, error) {

	authMethod := dao.UserAuth{
		ID: ID,
	}
	var err = u.db.Where(&authMethod).Preload("User").First(&authMethod).Error
	if err != nil {
		log.Error("Got an error when get user. Error: ", err)
		return dao.UserAuth{}, err
	}
	return authMethod, nil
}

func (u UserRepositoryImpl) Create(
	data string,
	method string) (dao.User, error) {
	user := dao.User{}
	user, err := u.Save(&user)
	if err != nil {
		log.Error("Got an error when creating user. Error: ", err)
		return dao.User{}, err
	}
	authMethod := dao.UserAuth{
		AuthData:   data,
		AuthMethod: method,
		User:       user,
	}
	err = u.db.Save(&authMethod).Error
	if err != nil {
		log.Error("Got an error when creating user. Error: ", err)
		return dao.User{}, err
	}
	return user, nil
}

func (u UserRepositoryImpl) GetOrCreate(
	data string,
	method string) (dao.User, error) {
	user, err := u.GetByAuth(data, method)
	if err != nil {
		log.Infoln("Creating user with auth method ", method, " and data ", data)
		user, err = u.Create(data, method)
		if err != nil {
			return dao.User{}, err
		}
	}
	return user, nil
}
func (u UserRepositoryImpl) GetOrCreateAuth(
	data string,
	method string) (dao.UserAuth, error) {
	userAuth, err := u.GetByAuthObj(data, method)
	if err != nil {
		log.Infoln("Creating user with auth method ", method, " and data ", data)
		_, err := u.Create(data, method)
		if err != nil {
			return dao.UserAuth{}, err
		}
	}
	return userAuth, nil
}

func (u UserRepositoryImpl) UpdateUserFields(userId uuid.UUID, updates map[string]interface{}) (dao.User, error) {
	user := dao.User{
		ID: userId,
	}
	if err := u.db.Model(&user).Updates(updates).Error; err != nil {
		return dao.User{}, err
	}
	return user, nil
}

func (u UserRepositoryImpl) GetUserUpgrade(userId uuid.UUID) (dao.UserUpgrade, error) {
	var userUpgrade dao.UserUpgrade
	if err := u.db.Where("user_id = ?", userId).First(&userUpgrade).Error; err == nil {
		return userUpgrade, nil
	}
	userUpgrade = dao.UserUpgrade{
		UserID:  userId,
		FarmLvl: 1,
	}
	if err := u.db.Create(&userUpgrade).Error; err != nil {
		return dao.UserUpgrade{}, err
	}
	return userUpgrade, nil
}

func (u UserRepositoryImpl) GetMyFields(userId uuid.UUID) ([]dao.UserField, error) {
	var userFields []dao.UserField
	if err := u.db.Where("user_id = ?", userId).Find(&userFields).Error; err != nil {
		return nil, err
	}
	return userFields, nil
}

func (u UserRepositoryImpl) GetMyReferrals(userId uuid.UUID) ([]dao.User, error) {
	var userReferrals []dao.UserReferral
	if err := u.db.Where("referrer_id = ?", userId).Preload("Referral").Find(&userReferrals).Error; err != nil {
		return nil, err
	}
	var userReferralsAsUser []dao.User
	for _, referral := range userReferrals {
		userReferralsAsUser = append(userReferralsAsUser, referral.Referral)
	}
	return userReferralsAsUser, nil
}

func UserRepositoryInit(db *gorm.DB) *UserRepositoryImpl {
	_ = db.AutoMigrate(&dao.User{}, &dao.UserAuth{}, &dao.UserUpgrade{}, &dao.UserField{}, &dao.UserReferral{})
	return &UserRepositoryImpl{
		db: db,
	}
}
