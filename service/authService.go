package service

import (
	"context"
	"fmt"
	"github.com/twinj/uuid"
	"gorm.io/gorm"
	"strconv"
	"test1/infrastructure"
	"test1/model"
	"test1/repository"
	"test1/utils"
	"time"
)

type authService struct {
	db                *gorm.DB
	authRepository    repository.AuthRepository
	profileRepository repository.ProfileRepository
}

func (a *authService) CreateToken(profileId uint) (*model.TokenDetail, error) {
	var err error
	td := &model.TokenDetail{}
	// lay ho so nguoi dung tu db theo profileId
	profile, err := a.profileRepository.GetProfile(a.db, profileId)
	if err != nil {
		return nil, fmt.Errorf("Can't create token")
	}
	// tao token
	td.Email = profile.Email
	td.AtExpires = time.Now().Add(time.Hour * time.Duration(model.AccessTokenTime)).Unix()
	td.AccessUUID = utils.GetPattern(profileId) + uuid.NewV4().String()
	td.RtExpires = time.Now().Add(time.Hour * time.Duration(model.RefreshTokenTime)).Unix()
	td.RefreshUUID = utils.GetPattern(profileId) + uuid.NewV4().String()

	//tao access token
	atClaims := make(map[string]interface{})
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = profileId
	atClaims["profile_id"] = profile.ID
	atClaims["exp"] = td.AtExpires
	_, td.AccessToken, err = infrastructure.GetEncodeAuth().Encode(atClaims)
	if err != nil {
		return nil, err
	}
	rtClaims := make(map[string]interface{})
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = profileId
	rtClaims["profile_id"] = profile.ID
	rtClaims["exp"] = td.RtExpires
	_, td.RefreshToken, err = infrastructure.GetEncodeAuth().Encode(rtClaims)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (a *authService) CreateAuth(userID uint, td *model.TokenDetail) error {
	// convert Unix to UTC(to Time object)
	accessToken := time.Unix(td.AtExpires, 0)
	refreshToken := time.Unix(td.RtExpires, 0)
	now := time.Now()

	//
	if errAccess := infrastructure.
		GetRedisClient().
		Set(context.Background(), td.AccessUUID, strconv.Itoa(int(userID)), accessToken.Sub(now)).
		Err(); errAccess != nil {
		return errAccess
	}
	if errRefresh := infrastructure.GetRedisClient().
		Set(context.Background(), td.RefreshUUID, strconv.Itoa(int(userID)), refreshToken.Sub(now)).
		Err(); errRefresh != nil {
		return errRefresh
	}
	return nil
}

type AuthService interface {
	CreateToken(profileId uint) (*model.TokenDetail, error)
	CreateAuth(userID uint, td *model.TokenDetail) error
}

func NewAuthService() *authService {
	db := infrastructure.GetDB()
	authRepository := repository.NewAuthRepository()
	profileRepository := repository.NewProfileRepository()
	return &authService{
		db:                db,
		authRepository:    authRepository,
		profileRepository: profileRepository,
	}
}
