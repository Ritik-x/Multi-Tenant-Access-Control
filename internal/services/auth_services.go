package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)



type AuthService struct {
	jwtSecret string
}

type Claims struct {
	UserID         string `json:"user_id"`
	Email          string `json:"email"`
	OrganizationID string `json:"organization_id"`
jwt.RegisteredClaims
	
}


func NewAuthService (jwtSecret string ) *AuthService{
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}


func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *AuthService) CheckPassword( password string,
	passwordHash string,
) bool{
	err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	)

	return err == nil
}



func ( s *AuthService) GenerateAcessTokens (
	userID string ,
		email string,
	organiztionID string , 

)(string, error) {
	now := time.Now()
claims := Claims{
		UserID:         userID,
		Email:          email,
		OrganizationID: organiztionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
		token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
return token.SignedString([]byte(s.jwtSecret))
}





func ( s *AuthService) GenerateRefreshToken () (string ,string , error){
	randomBytes := make([] byte , 32)
	_,err := rand.Read(randomBytes)
	if err != nil {
				return "", "", err
	}

	refreshToken := base64.RawURLEncoding.EncodeToString(randomBytes)

	refreshTokenHash := s.HashRefreshToken(refreshToken)

	return refreshToken, refreshTokenHash, nil
}

func (s *AuthService) HashRefreshToken(
	refreshToken string,
) string {

	hash := sha256.Sum256([]byte(refreshToken))

	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func ( s *AuthService) GenerateInvitationToken() (string , string , error){
	randomBytes:= make([]byte , 32)
	_,err := rand.Read(randomBytes)
	if err != nil {
		return "", "", err
	}
	invitationToken := base64.URLEncoding.EncodeToString(randomBytes)
		invitationTokenHash := s.HashRefreshToken(invitationToken)
		return invitationToken, invitationTokenHash, nil
}