package jwt

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	settingRepo := repositories.NewISettingRepo()
	jwtSign, _ := settingRepo.Get(settingRepo.WithByKey("JWTSigningKey"))
	return &JWT{
		[]byte(jwtSign.Value),
	}
}

type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	ID   uint
	Name string
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	claims := CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: constant.JWTBufferTime,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(constant.JWTBufferTime))),
			Issuer:    constant.JWTIssuer,
		},
	}
	return claims
}

func (j *JWT) CreateToken(request CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &request)
	return token.SignedString(j.SigningKey)
}
