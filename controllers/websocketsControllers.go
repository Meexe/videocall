package controllers

import (
	"log"
	"net/http"
	"reflect"
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
	log.Printf("User %s connected\n", user.Nickname)
	onlineUsers.Lock()
	onlineUsers.Users[user.Nickname] = true
	onlineUsers.Unlock()

	var oldstate []string
	wg := sync.WaitGroup{}
	wg.Add(2)
	defer wg.Wait()

	go func() {
		var msg string
		defer wg.Done()
		for {
			err := ws.ReadJSON(msg)
			if err != nil {
				log.Println("read:", err)
				log.Printf("User %s disconnected\n", user.Nickname)
				onlineUsers.Lock()
				delete(onlineUsers.Users, user.Nickname)
				onlineUsers.Unlock()
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			index := 0
			flag := false

			onlineUsers.Lock()
			users := make([]string, len(onlineUsers.Users))
			for key := range onlineUsers.Users {
				if key == user.Nickname {
					flag = true
					continue
				}
				users[index] = key
				index++
			}
			onlineUsers.Unlock()

			if !flag {
				return
			}

			if !reflect.DeepEqual(oldstate, users) {
				oldstate = users
				err := ws.WriteJSON(users[:len(users)-1])
				if err != nil {
					log.Println("write:", err)
					log.Printf("User %s disconnected\n", user.Nickname)
					onlineUsers.Lock()
					delete(onlineUsers.Users, user.Nickname)
					onlineUsers.Unlock()
					return
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
