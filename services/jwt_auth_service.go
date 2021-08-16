package services

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/lib"
	dto "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/dto"
	entity "github.com/LinggaAskaEdo/gin-gorm-clean-arch/models/entity"
	"github.com/LinggaAskaEdo/gin-gorm-clean-arch/repository"

	"github.com/golang-jwt/jwt"
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

// SplitToken
func (s JWTAuthService) ExtractToken(authHeader string) (string, error) {
	s.logger.Debug("ExtractToken")
	t := strings.Split(authHeader, " ")

	if len(t) == 2 {
		return t[1], nil
	}

	return "", errors.New("token malformed")
}

// VerifyToken
func (s JWTAuthService) VerifyToken(tokenString string) (*jwt.Token, error) {
	s.logger.Debug("VerifyToken")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(s.env.AccessSecret), nil
	})
	if err != nil {
		return nil, errors.New("token parse error")
	}

	return token, nil
}

// Authorize authorizes the generated token
func (s JWTAuthService) AuthorizeToken(tokenString string) (bool, error) {
	s.logger.Debug("AuthorizeToken")

	token, err := s.VerifyToken(tokenString)
	if err != nil {
		return false, errors.New("token malformed")
	}

	if token.Valid {
		// TODO: add logic check uuid is exist in redis, if exist next, if not exist abort
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			accessUUID, ok := claims["access_uuid"].(string)
			if !ok {
				return false, errors.New("token malformed")
			}

			if s.FetchToken(accessUUID) {
				return true, nil
			}
		}

		return false, errors.New("token deleted")
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

// ExtractTokenMetadata function
func (s JWTAuthService) ExtractTokenMetadata(tokenString string) (*dto.TokenDetails, error) {
	s.logger.Debug("ExtractTokenMetadata")

	token, err := s.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		refreshUUID, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, err
		}

		return &dto.TokenDetails{
			AccessUUID:  accessUUID,
			RefreshUUID: refreshUUID,
		}, nil
	}

	return nil, err
}

// CreateToken creates jwt auth token
func (s JWTAuthService) CreateToken(user entity.User) (*dto.TokenDetails, error) {
	s.logger.Debug("CreateToken")

	tokenDetails := &dto.TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	tokenDetails.RtExpires = time.Now().Add(time.Hour * 3).Unix()
	tokenDetails.AccessUUID = uuid.NewV4().String()
	tokenDetails.RefreshUUID = uuid.NewV4().String()

	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenDetails.AccessUUID
	atClaims["refresh_uuid"] = tokenDetails.RefreshUUID
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
	s.logger.Debug("StoreToken")

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

// FetchToken check uuid is exist in redis
func (s JWTAuthService) FetchToken(uuid string) bool {
	_, err := s.repository.Get(uuid).Result()

	return err == nil
}
