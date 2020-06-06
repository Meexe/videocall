package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	u "github.com/Meexe/videocall/utils"
	"golang.org/x/crypto/bcrypt"
	"os"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type User struct {
	gorm.Model
	Nickname	string `json:"nickname"`
	Password	string `json:"password"`
	Token		string `json:"token";sql:"-"`
}

//Validate incoming user details...
func (user *User) Validate() (map[string]interface{}, bool) {

	if len(user.Nickname) < 4 {
		return u.Message(false, "Nickname should contain 4 letters or more"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Nickname must be unique
	temp := &User{}

	//check for errors and duplicate nicknames
	err := GetDB().Table("users").Where("nickname = ?", user.Nickname).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Nickname != "" {
		return u.Message(false, "Nickname already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (user *User) Create() (map[string]interface{}) {

	if resp, ok := user.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	GetDB().Create(user)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	//Create new JWT token for the newly registered user
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //delete password

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}

func Login(nickname, password string) (map[string]interface{}) {

	user := &User{}
	err := GetDB().Table("users").Where("nickname = ?", nickname).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Nickname not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	user.Password = ""

	//Create JWT token
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

func GetUser(u uint) *User {

	acc := &User{}
	GetDB().Table("users").Where("id = ?", u).First(acc)
	if acc.Nickname == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}