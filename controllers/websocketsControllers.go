package controllers

import (
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/Meexe/videocall/models"
	u "github.com/Meexe/videocall/utils"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type        string `json:"type"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Payload     string `json:"payload"`
}

type Connections struct {
	sync.Mutex
	Connections map[string]*websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (connections *Connections) GetOnlineUsers(w http.ResponseWriter, r *http.Request) {

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
	connections.Lock()
	connections.Connections[user.Nickname] = ws
	connections.Unlock()

	var oldstate []string
	wg := sync.WaitGroup{}
	wg.Add(2)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		for {
			msg := Message{}
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				log.Printf("User %s disconnected\n", user.Nickname)
				connections.Lock()
				delete(connections.Connections, user.Nickname)
				connections.Unlock()
				return
			}
			log.Printf("got message from %s: %v\n", user.Nickname, msg)
			if msg.Type == "call" {
				conn := connections.Connections[msg.Destination] // ToDo validation
				conn.WriteJSON(msg)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			index := 0
			flag := false

			connections.Lock()
			users := make([]string, len(connections.Connections))
			for key := range connections.Connections {
				if key == user.Nickname {
					flag = true
					continue
				}
				users[index] = key
				index++
			}
			connections.Unlock()

			if !flag {
				return
			}

			if !reflect.DeepEqual(oldstate, users) {
				oldstate = users
				msg := u.WsMessage("users", users[:len(users)-1])
				err := ws.WriteJSON(msg)
				if err != nil {
					log.Println("write:", err)
					log.Printf("User %s disconnected\n", user.Nickname)
					connections.Lock()
					delete(connections.Connections, user.Nickname)
					connections.Unlock()
					return
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}
