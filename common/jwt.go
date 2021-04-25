package common

import (
	"Demo/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("a_key_by_D_aemon")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

//发布token
func ReleaseToken(user model.User) (string, error) {
	//设置过期时间
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	//创建认证
	claims := &Claims{
		UserId:         user.ID,
		StandardClaims: jwt.StandardClaims{
			//过期时间
			ExpiresAt: expirationTime.Unix(),
			//发放时间
			IssuedAt: time.Now().Unix(),
			//发放者
			Issuer: "D_aemon",
			//主题
			Subject: "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{},err error) {
		return jwtKey, nil
	})

	return token, claims, err
}


//token 由三部分组成：1、协议头(header)，储存token使用的加密协议。
//2、有效载荷（Playload），有效载荷中存放了token的签发者（iss）、签发时间（iat）、过期时间（exp）等以及一些我们需要写进token中的信息。
//3、签名(Signature)，将头部和有效荷载，再加上key来哈希的一个值。
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
//	eyJVc2VySWQiOjQsImV4cCI6MTYxODk5MzU3MywiaWF0IjoxNjE4Mzg4NzczLCJpc3MiOiJ6cHBqIiwic3ViIjoidXNlciB0b2tlbiJ9.
//		SYCsoxkhdq1Ru-ZVb6AFic_0QQems4JuKLXxzs1Q6wM