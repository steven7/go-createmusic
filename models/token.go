package models

import (

	"github.com/dgrijalva/jwt-go"
	//m "steve.com/go-textbox/pkg/models"

)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}