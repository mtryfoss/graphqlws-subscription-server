package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/functionalfoundry/graphqlws"
)

type ConnectedUser struct {
	jwt.StandardClaims
}

func (u ConnectedUser) Name() string {
	return u.Subject
}

func AuthenticateCallback(secretkey string) graphqlws.AuthenticateFunc {
	return func(tokenstring string) (interface{}, error) {
		user := ConnectedUser{}
		_, err := jwt.ParseWithClaims(tokenstring, &user, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretkey), nil
		})
		return user, err
	}
}
