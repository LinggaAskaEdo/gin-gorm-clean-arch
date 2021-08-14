package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"
	dto "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/dto"
	entity "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/entity"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

// JWTAuthService service relating to authorization
type JWTAuthService struct {
	env        lib.Env
	logger     lib.Logger
	repository repository.RedisRepository
}

// NewJWTAuthService creates a new auth service
func NewJWTAuthService(env lib.Env, logger lib.Logger, repository repository.RedisRepository) JWTAuthService {
	return JWTAuthService{
		env:        env,
		logger:     logger,
		repository: repository,
	}
}

// Authorize authorizes the generated token
func (s JWTAuthService) Authorize(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.env.AccessSecret), nil
	})
	if token.Valid {
		return true, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New("token malformed")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, errors.New("token expired")
		}
	}
	return false, errors.New("couldn't handle token")
}

// CreateToken creates jwt auth token
func (s JWTAuthService) CreateToken(user entity.User) (*dto.TokenDetails, error) {
	tokenDetails := &dto.TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	tokenDetails.AccessUUID = uuid.NewV4().String()
	tokenDetails.RefreshUUID = uuid.NewV4().String()

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDetails.AccessUUID
	atClaims["id"] = user.ID
	atClaims["name"] = user.Name
	atClaims["email"] = user.Email
	atClaims["exp"] = tokenDetails.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tokenDetails.AccessToken, err = at.SignedString([]byte(s.env.AccessSecret))
	if err != nil {
		s.logger.Error("JWT validation AccessToken failed: ", err)
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDetails.RefreshUUID
	rtClaims["id"] = user.ID
	rtClaims["name"] = user.Name
	rtClaims["email"] = user.Email
	rtClaims["exp"] = tokenDetails.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	tokenDetails.RefreshToken, err = rt.SignedString([]byte(s.env.RefreshSecret))
	if err != nil {
		s.logger.Error("JWT validation RefreshToken failed: ", err)
		return nil, err
	}

	return tokenDetails, nil
}

// StoreToken stores jwt auth token into redis
func (s JWTAuthService) StoreToken(user entity.User, token dto.TokenDetails) error {
	at := time.Unix(token.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(token.RtExpires, 0)
	now := time.Now()

	errAccess := s.repository.Set(token.AccessUUID, strconv.Itoa(int(user.ID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := s.repository.Set(token.RefreshUUID, strconv.Itoa(int(user.ID)), rt.Sub(now)).Err()

	if errRefresh != nil {
		return errRefresh
	}

	return nil
}
