package lib

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/onlineGo/conf"
	"time"
)

var tokenSign string

func init() {
	tokenSign = conf.GetConfig("jwtTokenSign")
}

type JwtClaims struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	UserName string `json:"user_name"`
	UserAccount string `json:"user_account"`
	UserRole string `json:"user_account"`
}


func NewJwt(param map[string]string) JwtClaims {
	jt := JwtClaims{
		UserId: param["user_id"],
		UserName: param["user_name"],
		UserAccount: param["user_account"],
		UserRole: param["user_account"],
	}
	jt.setIssuer()
	return jt
}

func (j JwtClaims) setIssuer () {
	j.Issuer = tokenSign
}

func (j JwtClaims) SetExpireAt (seconds uint) {
	afterSecond :=   time.Second * time.Duration(seconds)
	j.ExpiresAt = int64(time.Now().Add(afterSecond).Second())
}

func (j JwtClaims) GetToken () (string,error)  {
	Token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), j)
	return Token.SigningString()
}

func (j JwtClaims) SignValid () error  {
	//if j.Valid() != nil {
	//	return  j.Valid()
	//}
	if j.Issuer != tokenSign {
		return  errors.New("the sign is not valid")
	}
	return nil

}

func CheckToken (jwtToken string) (JwtClaims , error) {
	jwtClaim := JwtClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, &jwtClaim, func(token *jwt.Token) (interface{}, error) {
		if token.Valid {
			return token, nil
		}
		return token, errors.New("the token is not valid")
	})
	if jwtClaim.SignValid() != nil {
		return jwtClaim, jwtClaim.SignValid()
	}
	return jwtClaim,err
}


