package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/MarselBissengaliyev/ggp-blog/config"
	"github.com/MarselBissengaliyev/ggp-blog/models"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       uint
	Role      string
	jwt.StandardClaims
}

func GenerateAllTokens(
	email string,
	firstName string,
	lastName string,
	uid uint,
	role string,
	tokenSecret string,
) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       uid,
		Role:      role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(tokenSecret))
	if err != nil {
		log.Panicf("error occured while signing token: %s", err.Error())
		return
	}

	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	).SignedString([]byte(tokenSecret))
	if err != nil {
		log.Panicf("error occured while signing refresh token: %s", err.Error())
		return
	}

	return token, refreshToken, err
}

func CreateTokens(
	signedToken string,
	signedRefreshToken string,
	userAgent string,
	userId uint,
	db *gorm.DB,
) (models.Token, error) {
	token := models.Token{
		AccessToken:  signedToken,
		RefreshToken: signedRefreshToken,
		UserId:       userId,
		UserAgent:    userAgent,
	}

	if err := db.Create(&token).Error; err != nil {
		return token, err
	}

	return token, nil
}

func DeleteTokens(db *gorm.DB, uid uint, ua string) error {
	var token models.Token

	if err := db.Delete(&token, fmt.Sprintf("user_id = %d AND user_agent = '%s'", uid, ua)).Error; err != nil {
		return err
	}

	return nil
}

func ValidateToken(signedToken string, config *config.Config) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Token_Secret), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}

	return claims, msg
}
