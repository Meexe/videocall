package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Meexe/videocall/models"
	"github.com/gorilla/websocket"
)

type Users struct {
	sync.Mutex
	Users map[string]bool
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (onlineUsers *Users) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}
	defer ws.Close()

	ctx := r.Context().Value("user")
	userID := ctx.(uint)
	user := models.GetUser(userID)
	onlineUsers.Lock()
	onlineUsers.Users[user.Nickname] = true // append(onlineUsers.users, user.Nickname)
	onlineUsers.Unlock()

	for {
		users := make([]string, len(onlineUsers.Users)-1)
		index := 0
		for key := range onlineUsers.Users {
			if key == user.Nickname {
				continue
			}
			users[index] = key
			index++
		}
		msg, err := json.Marshal(users)
		if err != nil {
			log.Printf("json error: %s", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write:", err)
			log.Printf("User %s disconnected\n", user.Nickname)
			delete(onlineUsers.Users, user.Nickname)
			break
		}
		time.Sleep(10 * time.Second)
	}
}
