package models

import (
	"os"

	u "github.com/Meexe/videocall/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (user *User) Validate() (map[string]interface{}, bool) {

	if len(user.Nickname) < 4 {
		return u.Message(false, "Nickname should contain 4 letters or more"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	tmp := &User{}

	err := GetDB().Table("users").Where("nickname = ?", user.Nickname).First(tmp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if tmp.Nickname != "" {
		return u.Message(false, "Nickname already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) Create() map[string]interface{} {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = ""

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}

func (user *User) Login() map[string]interface{} {

	tmp := &User{}
	err := GetDB().Table("users").Where("nickname = ?", user.Nickname).First(tmp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Nickname not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(tmp.Password), []byte(user.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	user = tmp
	user.Password = ""

	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

func GetUser(u uint) *User {

	user := &User{}
	GetDB().Table("users").Where("id = ?", u).First(user)
	if user.Nickname == "" {
		return nil
	}

	user.Password = ""
	return user
}
