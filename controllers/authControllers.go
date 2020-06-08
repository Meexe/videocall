package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Meexe/videocall/models"
	u "github.com/Meexe/videocall/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create()

	if user.Token != "" {
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    user.Token,
			Path:     "/api",
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)
	}
	u.Respond(w, resp)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Login()

	if user.Token != "" {
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    user.Token,
			Path:     "/api",
			HttpOnly: false,
		}
		http.SetCookie(w, &cookie)
	}
	u.Respond(w, resp)
}

func Echo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("echo"))
}
