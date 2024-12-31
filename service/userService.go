package service

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"test1/infrastructure"
	"test1/model"
	"test1/repository"
)

type userService struct {
	db                *gorm.DB
	profileRepository repository.ProfileRepository
}
type UserService interface {
	CheckUsernameAndPassword(username string, password string) (*model.Profile, error)
	CreateUser(username string, password string, phone string, email string) (*model.Profile, error)
	CheckUserExist(Phone string, Email string) bool
	CheckGoogleToken(token string) (*model.Profile, error)
}

func NewUserService() UserService {
	profileRepo := repository.NewProfileRepository()
	db := infrastructure.GetDB()
	return &userService{
		db:                db,
		profileRepository: profileRepo,
	}
}

func (u *userService) CheckUsernameAndPassword(username string, password string) (*model.Profile, error) {
	profile, err := u.profileRepository.GetProfileByUsername(u.db, username)
	if err != nil {
		return nil, err
	}
	// nếu dùng băm mật khẩu trong csdl thì dùng hàm này sau này sẽ tiện hơn
	decryptpassword, err := infrastructure.RsaDecrypt(profile.User.Password)
	if err != nil {
		return nil, err
	}
	if string(decryptpassword) != password {
		return nil, err
	}

	//if !comparePassword(profile.User.Password, password) {
	//	return nil, errors.New("Password is incorrect")
	//}
	return profile, nil
}
func (u *userService) CreateUser(username string, password string, phone string, email string) (*model.Profile, error) {
	user := model.User{
		Username: username,
		Password: password,
	}
	profile := model.Profile{
		Phone: phone,
		Email: email,
		User:  user,
	}
	if err := u.profileRepository.CreateProfile(u.db, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}
func (u *userService) CheckUserExist(Phone string, Email string) bool {
	profile, err := u.profileRepository.GetProfileByPhoneOrEmail(u.db, Phone, Email)
	if err != nil {
		return false
	}
	if profile != nil {
		return true
	}
	return false
}
func (u *userService) CheckGoogleToken(token string) (*model.Profile, error) {
	return nil, nil
}

// nếu dùng băm mật khẩu trong csdl thì dùng hàm này sau này sẽ tiện hơn
func comparePassword(hashedPwd string, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		return false
	}
	return true
}
