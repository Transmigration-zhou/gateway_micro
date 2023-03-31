package public

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func JwtDecode(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("token is not jwt.RegisteredClaims")
	}
}

func JwtEncode(claims jwt.RegisteredClaims) (string, error) {
	mySigningKey := []byte(JwtSignKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
