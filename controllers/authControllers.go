package controllers

import (
	"encoding/json"
	"github.com/Meexe/videocall/models"
	u "github.com/Meexe/videocall/utils"
	"net/http"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create() //Create account

	cookie := http.Cookie{
		Name:		"jwt",
		Value:		user.Token,
		Path:		"/api",
		HttpOnly: 	false,
	}
	http.SetCookie(w, &cookie)
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(user.Nickname, user.Password)
	u.Respond(w, resp)
}